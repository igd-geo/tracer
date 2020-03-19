package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"geocode.igd.fraunhofer.de/hummer/tracer/client/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/broker"
	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/database"
	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/provenance"
)

// Tracer is a provenance service.
type Tracer struct {
	broker       *broker.Session
	database     *database.Client
	config       *config.Config
	txn          *database.Transaction
	mux          sync.Mutex
	batchTimeout time.Duration
	derivatives  chan *provenance.Entity
	cache        cache
}

type cache struct {
	items map[string]string
}

// New Returns a new tracer instance
func New(config *config.Config, database *database.Client, brokerSession *broker.Session) *Tracer {
	tracer := Tracer{
		config:       config,
		database:     database,
		broker:       brokerSession,
		derivatives:  make(chan *provenance.Entity),
		batchTimeout: time.Duration(config.Batch.Timeout) * time.Millisecond,
		cache: cache{
			items: make(map[string]string),
		},
	}
	return &tracer
}

// Close Initializes a graceful shutdown of all used components
func (t *Tracer) Close() error {
	return t.broker.Close()
}

// Listen Starts the tracer service
func (t *Tracer) Listen() error {
	t.txn = t.database.NewTransaction()
	err := t.broker.NewProducer()
	if err != nil {
		return err
	}

	err = t.broker.NewConsumer(t.msgHandler)
	if err != nil {
		return err
	}

	go t.batchHandler()

	return nil
}

func (t *Tracer) msgHandler(derivative *provenance.Entity) {
	prepared, err := t.prepare(derivative)
	if err != nil {
		log.Printf("failed to prepare provenance derivative: %e", err)
		return
	}
	t.derivatives <- prepared
}

func (t *Tracer) batchHandler() {
	for {
		select {
		case derivative := <-t.derivatives:
			t.txn.Add(derivative)
			if t.txn.Size >= t.config.Batch.Size {
				log.Printf("mutation batch full, commiting transaction")
				t.commitTransaction()
				t.txn = t.database.NewTransaction()
			}
		case <-time.After(t.batchTimeout):
			t.commitTransaction()
			derivative := <-t.derivatives
			t.txn = t.database.NewTransaction()
			t.txn.Add(derivative)
		}
	}
}

func (t *Tracer) commitTransaction() error {
	_, err := t.database.RunMutation(t.txn.Mutation)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tracer) prepareActivity(activity *provenance.Activity, query *database.Query) bool {
	for _, entity := range activity.Used {
		t.prepareUsedEntities(entity)
	}

	if uid, ok := t.cache.get(activity.ID); ok {
		activity.UID = uid
		query.SetVariable(database.VariableActivityID, "")
		return false
	} else if activity.IsBatch {
		query.SetVariable(database.VariableActivityID, activity.ID)
		return true
	}
	return false
}

func (t *Tracer) prepareUsedEntities(entity *provenance.Entity) {
	entityUID := fmt.Sprintf("_:%s", entity.ID)

	if uid, ok := t.cache.get(entity.ID); ok {
		entityUID = uid
	} else {
		query := database.NewQuery(database.QueryEntityUIDByID)
		query.SetVariable(database.VariableEntityID, entity.ID)
		result, _ := t.database.RunQueryWithVars(query)

		if len(result.Entity) == 1 {
			entityUID = result.Entity[0].UID
		}

		t.cache.set(entity.ID, entityUID)
	}

	entity.UID = entityUID
}

func (t *Tracer) prepareAgent(agent *provenance.Agent, query *database.Query, isSupervisor bool) bool {
	if uid, ok := t.cache.get(agent.ID); ok {
		agent.UID = uid
		return false
	}
	if isSupervisor {
		query.SetVariable(database.VariableSupervisorID, agent.ID)
	} else {
		query.SetVariable(database.VariableAgentID, agent.ID)
	}
	return true
}

func (t *Tracer) prepare(derivative *provenance.Entity) (*provenance.Entity, error) {
	activity := derivative.WasGeneratedBy[0]
	agent := derivative.WasGeneratedBy[0].WasAssociatedWith[0]
	supervisor := derivative.WasGeneratedBy[0].WasAssociatedWith[0].ActedOnBehalfOf[0]

	query := database.NewQuery(database.QueryAllUIDsByID)

	activityMissed := t.prepareActivity(activity, query)
	agentMissed := t.prepareAgent(agent, query, false)
	supervisorMissed := t.prepareAgent(supervisor, query, true)

	if activityMissed || agentMissed || supervisorMissed {
		err := t.fetchAndCacheUids(derivative, query)
		if err != nil {
			return nil, err
		}
	}

	derivative.WasDerivedFrom = activity.Used

	return derivative, nil
}

func (t *Tracer) fetchAndCacheUids(derivative *provenance.Entity, query *database.Query) error {
	activity := derivative.WasGeneratedBy[0]
	agent := derivative.WasGeneratedBy[0].WasAssociatedWith[0]
	supervisor := derivative.WasGeneratedBy[0].WasAssociatedWith[0].ActedOnBehalfOf[0]

	activityUID := fmt.Sprintf("_:%s", activity.ID)
	agentUID := fmt.Sprintf("_:%s", agent.ID)
	supervisorUID := fmt.Sprintf("_:%s", supervisor.ID)

	result, err := t.database.RunQueryWithVars(query)
	if err != nil {
		return err
	}

	if len(result.Activity) == 1 {
		activityUID = result.Activity[0].UID
	} else {
		t.fetchActivityInfoFromSysService(activity)
	}

	if len(result.Agent) == 1 {
		agentUID = result.Agent[0].UID
	} else {
		t.fetchAgentInfoFromSysService(agent, false)
	}

	if len(result.Supervisor) == 1 {
		supervisorUID = result.Supervisor[0].UID
	} else {
		t.fetchAgentInfoFromSysService(agent, true)
	}

	activity.UID = activityUID
	agent.UID = agentUID
	supervisor.UID = supervisorUID

	t.mux.Lock()
	t.cache.set(activity.ID, activityUID)
	t.cache.set(agent.ID, agentUID)
	t.cache.set(supervisor.ID, supervisorUID)
	t.mux.Unlock()

	return nil
}

func (c *cache) get(key string) (string, bool) {
	value, ok := c.items[key]
	return value, ok
}

func (c *cache) set(key string, value string) {
	c.items[key] = value
}

func (t *Tracer) fetchAgentInfoFromSysService(agent *provenance.Agent, isSupervisor bool) error {
	type mockServiceStruct struct {
		ID          string
		Name        string
		Description string
		Type        string
		CreatedBy   string
	}

	type mockUserStruct struct {
		ID   string
		Name string
		Type string
	}
	if isSupervisor {
		supervisor := agent.ActedOnBehalfOf[0]
		url := fmt.Sprintf("%s/%s", t.config.Arbiter.UserRegistry.URL, supervisor.UID)
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		mockUser := mockUserStruct{}
		json.NewDecoder(res.Body).Decode(&mockUser)

		supervisor.Name = mockUser.Name
		supervisor.Type = mockUser.Type
	} else {
		url := fmt.Sprintf("%s/%s", t.config.Arbiter.UserRegistry.URL, agent.UID)
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		mockService := mockServiceStruct{}
		json.NewDecoder(res.Body).Decode(&mockService)

		agent.Name = mockService.Name
		agent.Type = mockService.Type
		agent.Description = mockService.Type
		agent.ActedOnBehalfOf[0].ID = mockService.CreatedBy
	}
	return nil
}

func (t *Tracer) fetchActivityInfoFromSysService(activity *provenance.Activity) error {
	type mockProcessStruct struct {
		ID         string
		InstanceOf string
		StartDate  string
		EndDate    string
		Input      string
		Output     string
		Status     string
	}

	url := fmt.Sprintf("%s/%s", t.config.Arbiter.ScenarioRegistry.URL, activity.UID)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	mockProcess := mockProcessStruct{}
	return json.NewDecoder(res.Body).Decode(&mockProcess)
}
