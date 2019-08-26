package tracer

import (
	"fmt"
	"log"
	"time"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/broker"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
)

// Tracer is a provenance service.
type Tracer struct {
	deliveries <-chan *util.Entity
	rbSession  *broker.Session
	db         *db.Client
	config     *config.Config
	cache      cache
}

type mutation struct {
	entities []*util.Entity
}

type cache struct {
	items map[string]string
}

// Setup sets up the tracer service and initializes all required components
func Setup(config *config.Config, db *db.Client, broker *broker.Session, deliveries chan *util.Entity) *Tracer {
	tracer := Tracer{
		deliveries: deliveries,
		rbSession:  broker,
		db:         db,
		cache: cache{
			items: make(map[string]string),
		},
	}
	return &tracer
}

// Cleanup initlializes RabbitMQ Shutdown.
func (tracer *Tracer) Cleanup() error {
	return tracer.rbSession.Shutdown()
}

// Listen starts the tracer service.
func (tracer *Tracer) Listen() {
	derivatives := make(chan *util.Entity)
	commit := make(chan struct{})
	batchTimeout := time.Duration(tracer.config.BatchTimeout) * time.Millisecond

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
	txn := tracer.db.NewTransaction()
	for {
		select {
		case derivative := <-derivatives:
			tracer.prepare(derivative, txn)
			if txn.Size == tracer.config.BatchSizeLimit {
				log.Println("batch full, commiting...")
				tracer.commitTransaction(txn)
				txn = tracer.db.NewTransaction()
				tracer.cache.items = make(map[string]string)
			}
		case <-commit:
			if len(txn.Mutation) == 0 {
				continue
			}
			log.Println("no new delivery within time window, commiting...")
			tracer.commitTransaction(txn)
			txn = tracer.db.NewTransaction()
			tracer.cache.items = make(map[string]string)
		}
	}
}

func (tracer *Tracer) commitTransaction(txn *db.Transaction) {
	_, err := tracer.db.RunMutation(&txn.Mutation)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("successfully commited %v items\n", txn.Size)
	}
}

func (tracer *Tracer) prepareActivity(activity *util.Activity, query *db.Query) bool {
	missed := false

	if uid, ok := tracer.cache.get(activity.ID); ok {
		activity.UID = uid
		return missed
	} else if activity.IsBatch {
		query.SetVariable(db.VariableActivityID, activity.ID)
		missed = true
	}

	return missed
}

func (tracer *Tracer) prepareAgent(agent *util.Agent, query *db.Query, isSupervisor bool) bool {
	missed := false

	if uid, ok := tracer.cache.get(agent.ID); ok {
		agent.UID = uid
		return missed
	}
	if isSupervisor {
		query.SetVariable(db.VariableSupervisorID, agent.ID)
	} else {
		query.SetVariable(db.VariableAgentID, agent.ID)
	}
	missed = true

	return missed
}
func (tracer *Tracer) prepare(derivative *util.Entity, txn *db.Transaction) error {
	activity := derivative.WasGeneratedBy
	agent := derivative.WasGeneratedBy.WasAssociatedWith
	supervisor := derivative.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf

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

	log.Printf("appending entity %s to mutation", derivative.ID)
	txn.Mutation = append(txn.Mutation, derivative)
	txn.Size++

	return nil
}

func (tracer *Tracer) fetchAndCacheUids(derivative *util.Entity, query *db.Query) error {
	activity := derivative.WasGeneratedBy
	agent := derivative.WasGeneratedBy.WasAssociatedWith
	supervisor := derivative.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf

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

	if len(result.Agent) == 1 {
		agentUID = result.Agent[0].UID
	}

	if len(result.Supervisor) == 1 {
		supervisorUID = result.Supervisor[0].UID
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
