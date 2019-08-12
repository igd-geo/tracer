package tracer

import (
	"bytes"
	"encoding/json"
	"log"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/dgraph"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/mongodb"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rabbitmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/provutil"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
)

type Tracer struct {
	deliveries <-chan rabbitmq.Delivery
	rbSession  *rabbitmq.Session
	mongoDB    *mongodb.Client
	dgraph     *dgraph.Client
	config     *config.Config
}

type Delivery struct {
	Entity *provutil.Entity `json:"entity,omitempty"`
}

func New(config *config.Config) *Tracer {
	msgChan := make(chan rabbitmq.Delivery)
	tracer := Tracer{
		deliveries: msgChan,
		rbSession:  rabbitmq.New(config.RabbitURL, msgChan, config.ConsumerTag, "notifications", "topic"),
		mongoDB: mongodb.NewClient(
			config.MongoURL,
			config.MongoDatabase,
			config.MongoCollectionEntity,
			config.MongoCollectionAgent,
			config.MongoCollectionActivity,
		),
		dgraph: dgraph.NewClient(config.DgraphURL),
	}
	return &tracer
}

func (tracer *Tracer) Listen() {
	go func() {
		for delivery := range tracer.deliveries {
			go tracer.handleDelivery(delivery)
		}
	}()
}

func (tracer *Tracer) Cleanup() error {
	return tracer.rbSession.Shutdown()
}

func (tracer *Tracer) handleDelivery(rbDelivery rabbitmq.Delivery) {
	delivery := Delivery{Entity: provutil.NewEntity()}

	err := json.Unmarshal(rbDelivery, &delivery)
	if err != nil {
		log.Println(err)
		return
	}

	createAgent, createSupervisor, err := tracer.createGraphEntry(delivery.Entity)
	if err != nil {
		log.Println(err)
		return
	}

	err = tracer.createMongoEntries(delivery.Entity, createAgent, createSupervisor)
	if err != nil {
		log.Println(err)
		return
	}
}

func (tracer *Tracer) createGraphEntry(entity *provutil.Entity) (bool, bool, error) {
	var err error
	createAgent, createSupevisor := false, false

	if err := tracer.fetchUsedEntities(entity); err != nil {
		return createAgent, createSupevisor, err
	}

	if createAgent, createSupevisor, err = tracer.fetchAgents(entity); err != nil {
		return createAgent, createSupevisor, err
	}

	assigned, err := tracer.dgraph.AddDerivate(entity)
	if err != nil {
		return createAgent, createSupevisor, err
	}

	log.Printf("inserted dgraph entries as %v\n", assigned)

	entity.UID = assigned["entity"]
	entity.WasGeneratedBy.UID = assigned["activity"]
	if createAgent {
		entity.WasGeneratedBy.WasAssociatedWith.UID = assigned["agent"]
		log.Printf("agent <%s> not in database, creating entry\n", entity.WasGeneratedBy.WasAssociatedWith.ID)
	}
	if createSupevisor {
		entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID = assigned["supervisor"]
		log.Printf("agent <%s> not in database, creating entry\n", entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID)
	}

	return createAgent, createSupevisor, err
}

func (tracer *Tracer) fetchUsedEntities(entity *provutil.Entity) error {
	numEntities := len(entity.WasGeneratedBy.Used)
	if numEntities > 0 {
		for i, used := range entity.WasGeneratedBy.Used {
			usedUID := tracer.mongoDB.EntityUID(used.ID)
			if usedUID != "" {
				entity.WasGeneratedBy.Used[i].Attributes = provutil.NewAttributes()
				entity.WasGeneratedBy.Used[i].UID = usedUID
			}
		}
		log.Println("fetched uids of used entities")
	} else {
		log.Println("no used entities to fetch, skipping")
	}
	return nil
}

func (tracer *Tracer) fetchAgents(entity *provutil.Entity) (bool, bool, error) {
	var err error
	createAgent, createSupevisor := true, true

	agentUID := tracer.mongoDB.AgentUID(entity.WasGeneratedBy.WasAssociatedWith.ID)

	if agentUID != "" {
		entity.WasGeneratedBy.WasAssociatedWith.Attributes = provutil.NewAttributes()
		entity.WasGeneratedBy.WasAssociatedWith.UID = agentUID
		createAgent = false
	}

	supervisorUID := tracer.mongoDB.AgentUID(entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID)

	if supervisorUID != "" {
		entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.Attributes = provutil.NewAttributes()
		entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID = supervisorUID
		createSupevisor = false
	}
	return createAgent, createSupevisor, err
}

func (tracer *Tracer) createMongoEntries(entity *provutil.Entity, createAgent bool, createSupevisor bool) error {
	err := tracer.createMongoEntity(entity)
	if err != nil {
		return err
	}

	err = tracer.createMongoActivity(entity)
	if err != nil {
		return err
	}

	if createAgent {
		tracer.createMongoAgent(entity, false)
		if err != nil {
			return err
		}
	}

	if createSupevisor {
		tracer.createMongoAgent(entity, true)
		if err != nil {
			return err
		}
	}

	log.Println("inserted mongodb entries")
	return nil
}

func (tracer *Tracer) createMongoEntity(entity *provutil.Entity) error {
	var buffer bytes.Buffer
	err := json.Compact(&buffer, entity.Data)

	if err != nil {
		return err
	}

	entity.Data = buffer.Bytes()
	return tracer.mongoDB.InsertEntity(entity)
}

func (tracer *Tracer) createMongoAgent(entity *provutil.Entity, isSupervisor bool) error {
	var agent *provutil.Agent
	var buffer bytes.Buffer

	if isSupervisor {
		agent = entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf
	} else {
		agent = entity.WasGeneratedBy.WasAssociatedWith
	}
	err := json.Compact(&buffer, agent.Data)

	if err != nil {
		return err
	}

	agent.Data = buffer.Bytes()

	return tracer.mongoDB.InsertAgent(agent)
}

func (tracer *Tracer) createMongoActivity(entity *provutil.Entity) error {
	var buffer bytes.Buffer
	err := json.Compact(&buffer, entity.WasGeneratedBy.Data)

	if err != nil {
		return err
	}

	entity.WasGeneratedBy.Data = buffer.Bytes()
	return tracer.mongoDB.InsertActivity(entity.WasGeneratedBy)
}
