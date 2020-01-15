package tracer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rbmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
)

const mockServerURL = "http://localhost:8787"

// Tracer is a provenance service.
type Tracer struct {
	deliveries    <-chan *util.Entity
	rbmq          *rbmq.Session
	db            *db.Client
	conf          *config.Config
	cache         cache
	batchTimes    []time.Duration
	httpTimes     []time.Duration
	mutationTimes []time.Duration
}

type cache struct {
	items map[string]string
}

// Setup sets up the tracer service and initializes all required components
func Setup(conf *config.Config, db *db.Client, rbSession *rbmq.Session, deliveries chan *util.Entity) *Tracer {
	tracer := Tracer{
		conf:       conf,
		db:         db,
		rbmq:       rbSession,
		deliveries: deliveries,
		cache: cache{
			items: make(map[string]string),
		},
	}
	return &tracer
}

// Cleanup initializes RabbitMQ Shutdown.
func (tracer *Tracer) Cleanup() error {
	return tracer.rbmq.Shutdown()
}

// Listen starts the tracer service.
func (tracer *Tracer) Listen() {
	derivatives := make(chan *util.Entity)
	commit := make(chan struct{})
	batchTimeout := time.Duration(tracer.conf.BatchTimeout) * time.Millisecond

	go func() {
		for {
			select {
			case derivative := <-tracer.deliveries:
				derivatives <- derivative
			case <-time.After(batchTimeout):
				commit <- struct{}{}
				derivative := <-tracer.deliveries
				derivatives <- derivative
			}
		}
	}()

	go tracer.handleDerivatives(derivatives, commit)
}

func (tracer *Tracer) handleDerivatives(derivatives <-chan *util.Entity, commit <-chan struct{}) {
	first := true
	txn := tracer.db.NewTransaction()
	for {
		select {
		case derivative := <-derivatives:
			if first {
				txn.StartTime = time.Now()
			}
			err := tracer.prepare(derivative, txn)
			if err != nil {
				log.Println(err)
			}
			if txn.Size == tracer.conf.BatchSizeLimit {
				log.Println("batch full, committing...")
				tracer.commitTransaction(txn)
				txn = tracer.db.NewTransaction()
				tracer.cache.items = make(map[string]string)
			}
		case <-commit:
			if len(txn.Mutation) == 0 {
				continue
			}
			log.Println("no new delivery within time window, committing...")
			tracer.commitTransaction(txn)
			txn = tracer.db.NewTransaction()
			tracer.cache.items = make(map[string]string)
		}
	}
}

func (tracer *Tracer) commitTransaction(txn *db.Transaction) {
	mutationStart := time.Now()
	_, err := tracer.db.RunMutation(&txn.Mutation)
	if err != nil {
		log.Println(err)
	} else {
		//log.Printf("%+v", assigned)
		log.Printf("successfully commited %v items\n", txn.Size)
		tracer.mutationTimes = append(tracer.mutationTimes, time.Since(mutationStart))
		tracer.batchTimes = append(tracer.batchTimes, time.Since(txn.StartTime))
	}

	var totalBatch time.Duration = 0
	var totalHTTP time.Duration = 0
	var totalMutation time.Duration = 0
	batches := 0
	httpReqs := 0
	mutations := 0

	for count, dur := range tracer.batchTimes {
		totalBatch = totalBatch + dur
		batches = count + 1
	}

	for count, dur := range tracer.mutationTimes {
		totalMutation = totalMutation + dur
		mutations = count + 1
	}

	for count, dur := range tracer.httpTimes {
		totalHTTP = totalHTTP + dur
		httpReqs = count + 1
	}

	avgBatch := (float64(totalBatch.Nanoseconds()) / float64(batches)) / float64(time.Millisecond)
	avgHTTP := (float64(totalHTTP.Nanoseconds()) / float64(httpReqs)) / float64(time.Millisecond)
	avgMut := (float64(totalMutation.Nanoseconds()) / float64(mutations)) / float64(time.Millisecond)

	fmt.Printf("\nTotal Time: %v\n", totalBatch)
	//fmt.Printf("Batch Durations [%d batches]: %v\n", batches, tracer.batchTimes)
	fmt.Printf("Average Batch Duration [%d batches]: %f\n", batches, avgBatch)
	fmt.Printf("Total Mutation Request Time: %v\n", totalMutation)
	fmt.Printf("Average Mutation Request Duration: %v\n", avgMut)
	fmt.Printf("Total HTTP Request Time: %v\n", totalHTTP)
	fmt.Printf("Average HTTP Request Duration [%d requests]: %v\n\n\n", httpReqs, avgHTTP)
}

func (tracer *Tracer) prepareActivity(activity *util.Activity, query *db.Query) bool {
	for _, entity := range activity.Used {
		tracer.prepareUsedEntities(entity)
	}

	if uid, ok := tracer.cache.get(activity.ID); ok {
		activity.UID = uid
		query.SetVariable(db.VariableActivityID, "")
		return false
	} else if activity.IsBatch {
		query.SetVariable(db.VariableActivityID, activity.ID)
		return true
	}
	return false
}

func (tracer *Tracer) prepareUsedEntities(entity *util.Entity) {
	entityUID := fmt.Sprintf("_:%s", entity.ID)

	if uid, ok := tracer.cache.get(entity.ID); ok {
		entityUID = uid
	} else {
		query := db.NewQuery(db.QueryEntityUIDByID)
		query.SetVariable(db.VariableEntityID, entity.ID)
		result, _ := tracer.db.RunQueryWithVars(query)

		if len(result.Entity) == 1 {
			entityUID = result.Entity[0].UID
		}

		tracer.cache.set(entity.ID, entityUID)
	}

	entity.UID = entityUID
}

func (tracer *Tracer) prepareAgent(agent *util.Agent, query *db.Query, isSupervisor bool) bool {
	if uid, ok := tracer.cache.get(agent.ID); ok {
		agent.UID = uid
		return false
	}
	if isSupervisor {
		query.SetVariable(db.VariableSupervisorID, agent.ID)
	} else {
		query.SetVariable(db.VariableAgentID, agent.ID)
	}
	return true
}

func (tracer *Tracer) prepare(derivative *util.Entity, txn *db.Transaction) error {
	activity := derivative.WasGeneratedBy[0]
	agent := derivative.WasGeneratedBy[0].WasAssociatedWith[0]
	supervisor := derivative.WasGeneratedBy[0].WasAssociatedWith[0].ActedOnBehalfOf[0]

	query := db.NewQuery(db.QueryAllUIDsByID)

	activityMissed := tracer.prepareActivity(activity, query)
	agentMissed := tracer.prepareAgent(agent, query, false)
	supervisorMissed := tracer.prepareAgent(supervisor, query, true)

	if activityMissed || agentMissed || supervisorMissed {
		err := tracer.fetchAndCacheUids(derivative, query)
		if err != nil {
			return err
		}
	}

	derivative.WasDerivedFrom = activity.Used
	txn.Mutation = append(txn.Mutation, derivative)
	txn.Size++

	return nil
}

func (tracer *Tracer) fetchAndCacheUids(derivative *util.Entity, query *db.Query) error {
	activity := derivative.WasGeneratedBy[0]
	agent := derivative.WasGeneratedBy[0].WasAssociatedWith[0]
	supervisor := derivative.WasGeneratedBy[0].WasAssociatedWith[0].ActedOnBehalfOf[0]

	activityUID := fmt.Sprintf("_:%s", activity.ID)
	agentUID := fmt.Sprintf("_:%s", agent.ID)
	supervisorUID := fmt.Sprintf("_:%s", supervisor.ID)

	result, err := tracer.db.RunQueryWithVars(query)
	if err != nil {
		return err
	}

	if len(result.Activity) == 1 {
		activityUID = result.Activity[0].UID
	}
	// else {
	//fetchActivityInfoFromSysService(activity)
	//}

	if len(result.Agent) == 1 {
		agentUID = result.Agent[0].UID
	} else {
		fetchStart := time.Now()
		fetchAgentInfoFromSysService(agent, false)
		tracer.httpTimes = append(tracer.httpTimes, time.Since(fetchStart))
	}

	if len(result.Supervisor) == 1 {
		supervisorUID = result.Supervisor[0].UID
	} else {
		fetchStart := time.Now()
		fetchAgentInfoFromSysService(agent, true)
		tracer.httpTimes = append(tracer.httpTimes, time.Since(fetchStart))
	}

	activity.UID = activityUID
	agent.UID = agentUID
	supervisor.UID = supervisorUID

	tracer.cache.set(activity.ID, activityUID)
	tracer.cache.set(agent.ID, agentUID)
	tracer.cache.set(supervisor.ID, supervisorUID)

	return nil
}

func (c *cache) get(key string) (string, bool) {
	value, ok := c.items[key]
	return value, ok
}

func (c *cache) set(key string, value string) {
	c.items[key] = value
}

func fetchAgentInfoFromSysService(agent *util.Agent, isSupervisor bool) {
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
		url := fmt.Sprintf("%s/user/%s", mockServerURL, supervisor.ID)
		res, _ := http.Get(url)
		mockUser := mockUserStruct{}
		json.NewDecoder(res.Body).Decode(&mockUser)

		supervisor.Name = mockUser.Name
		supervisor.Type = mockUser.Type
	} else {
		url := fmt.Sprintf("%s/service/%s", mockServerURL, agent.ID)
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()

		mockService := mockServiceStruct{}
		json.NewDecoder(res.Body).Decode(&mockService)

		agent.Name = mockService.Name
		agent.Type = mockService.Type
		agent.Description = mockService.Type
		agent.ActedOnBehalfOf[0].ID = mockService.CreatedBy
	}
}

func fetchActivityInfoFromSysService(activity *util.Activity) {
	type mockProcessStruct struct {
		ID         string
		InstanceOf string
		StartDate  string
		EndDate    string
		Input      string
		Output     string
		Status     string
	}

	res, _ := http.Get(fmt.Sprintf("%s/user/%s", activity.ID, mockServerURL))
	mockProcess := mockProcessStruct{}
	json.NewDecoder(res.Body).Decode(&mockProcess)
}
