package tracer

import (
	"fmt"
	"log"
	"time"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rbmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
)

// Tracer is a provenance service.
type Tracer struct {
	deliveries <-chan *util.Entity
	rbmq       *rbmq.Session
	db         *db.Client
	conf       *config.Config
	cache      cache
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
	txn := tracer.db.NewTransaction()
	for {
		select {
		case derivative := <-derivatives:
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
	assigned, err := tracer.db.RunMutation(&txn.Mutation)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%+v", assigned)
		log.Printf("successfully commited %v items\n", txn.Size)
	}
}

func (tracer *Tracer) prepareActivity(activity *util.Activity, query *db.Query) bool {
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
