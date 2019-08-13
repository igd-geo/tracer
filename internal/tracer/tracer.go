package tracer

import (
	"bytes"
	"encoding/json"
	"log"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rabbitmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/provutil"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
)

type Tracer struct {
	deliveries <-chan *provutil.Entity
	rbSession  *rabbitmq.Session
	infoDB     provutil.InfoDB
	provDB     provutil.ProvDB
	config     *config.Config
}

func New(config *config.Config, infoDB provutil.InfoDB, provDB provutil.ProvDB) *Tracer {
	msgChan := make(chan *provutil.Entity)
	tracer := Tracer{
		deliveries: msgChan,
		rbSession:  rabbitmq.New(config.RabbitURL, msgChan, config.ConsumerTag, "notifications", "topic"),
		infoDB:     infoDB,
		provDB:     provDB,
	}
	return &tracer
}

func (tracer *Tracer) Listen() {
	go func() {
		for derivate := range tracer.deliveries {
			go tracer.handleDelivery(derivate)
		}
	}()
}

func (tracer *Tracer) Cleanup() error {
	return tracer.rbSession.Shutdown()
}

func (tracer *Tracer) handleDelivery(derivate *provutil.Entity) {
	createAgent, createSupervisor, err := tracer.createProvEntry(derivate)
	if err != nil {
		log.Println(err)
		return
	}

	err = tracer.createInfoEntries(derivate, createAgent, createSupervisor)
	if err != nil {
		log.Println(err)
		return
	}
}

func (tracer *Tracer) createProvEntry(entity *provutil.Entity) (bool, bool, error) {
	var err error
	createAgent, createSupevisor := false, false

	if err := tracer.fetchUsedEntities(entity); err != nil {
		return createAgent, createSupevisor, err
	}

	if createAgent, createSupevisor, err = tracer.fetchAgents(entity); err != nil {
		return createAgent, createSupevisor, err
	}

	assigned, err := tracer.provDB.InsertDerivate(entity)
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
			usedUID := tracer.infoDB.EntityUID(used.ID)
			if usedUID != "" {
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

	agentUID := tracer.infoDB.AgentUID(entity.WasGeneratedBy.WasAssociatedWith.ID)

	if agentUID != "" {
		entity.WasGeneratedBy.WasAssociatedWith.UID = agentUID
		createAgent = false
	}

	supervisorUID := tracer.infoDB.AgentUID(entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID)

	if supervisorUID != "" {
		entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID = supervisorUID
		createSupevisor = false
	}
	return createAgent, createSupevisor, err
}

func (tracer *Tracer) createInfoEntries(entity *provutil.Entity, createAgent bool, createSupevisor bool) error {
	err := tracer.addEntity(entity)
	if err != nil {
		return err
	}

	err = tracer.addActivity(entity)
	if err != nil {
		return err
	}

	if createAgent {
		tracer.addAgent(entity, false)
		if err != nil {
			return err
		}
	}

	if createSupevisor {
		tracer.addAgent(entity, true)
		if err != nil {
			return err
		}
	}

	log.Println("inserted mongodb entries")
	return nil
}

func (tracer *Tracer) addEntity(entity *provutil.Entity) error {
	var buffer bytes.Buffer
	err := json.Compact(&buffer, entity.Data)

	if err != nil {
		return err
	}

	entity.Data = buffer.Bytes()
	return tracer.infoDB.InsertEntity(entity)
}

func (tracer *Tracer) addAgent(entity *provutil.Entity, isSupervisor bool) error {
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

	return tracer.infoDB.InsertAgent(agent)
}

func (tracer *Tracer) addActivity(entity *provutil.Entity) error {
	var buffer bytes.Buffer
	err := json.Compact(&buffer, entity.WasGeneratedBy.Data)

	if err != nil {
		return err
	}

	entity.WasGeneratedBy.Data = buffer.Bytes()
	return tracer.infoDB.InsertActivity(entity.WasGeneratedBy)
}
