package tracer

import (
	"bytes"
	"encoding/json"
	"log"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/dgraph"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/mongodb"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rabbitmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/types"
)

type Tracer struct {
	deliveries <-chan rabbitmq.Delivery
	rbSession  *rabbitmq.Session
	mongoDB    *mongodb.Client
	dgraph     *dgraph.Client
	config     *config.Config
}

type Delivery struct {
	Entity Entity `json:"entity,omitempty"`
}

type DeliveryTest struct {
	Entity types.Entity `json:"entity,omitempty"`
}

type Entity struct {
	UID            string          `json:"uid,omitempty"`
	ID             string          `json:"id,omitempty"`
	URI            string          `json:"uri,omitempty"`
	Name           string          `json:"name,omitempty"`
	CreationDate   string          `json:"creationDate,omitempty"`
	Data           json.RawMessage `json:"data,omitempty"`
	WasGeneratedBy struct {
		UID               string          `json:"uid,omitempty"`
		ID                string          `json:"id,omitempty"`
		StartDate         string          `json:"startDate,omitempty"`
		EndDate           string          `json:"endDate,omitempty"`
		Data              json.RawMessage `json:"data,omitempty"`
		WasAssociatedWith struct {
			UID             string          `json:"uid,omitempty"`
			ID              string          `json:"id,omitempty"`
			Name            string          `json:"name,omitempty"`
			Data            json.RawMessage `json:"data,omitempty"`
			ActedOnBehalfOf struct {
				UID  string          `json:"uid,omitempty"`
				ID   string          `json:"id,omitempty"`
				Name string          `json:"name,omitempty"`
				Data json.RawMessage `json:"data,omitempty"`
			} `json:"actedOnBehalfOf,omitempty"`
		}
		Used []struct {
			UID string `json:"uid,omitempty"`
			ID  string `json:"id,omitempty"`
		} `json:"used,omitempty"`
	} `json:"wasGeneratedBy,omitempty"`
	WasDerivedFrom []struct {
		UID string `json:"uid,omitempty"`
		ID  string `json:"id,omitempty"`
	} `json:"wasDerivedFrom,omitempty"`
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
	delivery := Delivery{}
	deliveryTest := DeliveryTest{}
	err := json.Unmarshal(rbDelivery, &deliveryTest)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%+v", deliveryTest.Entity.Attributes)
	log.Printf("%+v", deliveryTest.Entity.ID)
	log.Printf("%+v", deliveryTest.Entity.Edges)
	log.Printf("%+v", string(deliveryTest.Entity.Data))
	log.Printf("%+v", deliveryTest.Entity.WasGeneratedBy.WasAssociatedWith.Attributes)
	log.Printf("%+v", deliveryTest.Entity.WasGeneratedBy.WasAssociatedWith.Edges)
	log.Printf("%+v", string(deliveryTest.Entity.WasGeneratedBy.WasAssociatedWith.Data))
	log.Printf("%+v", deliveryTest.Entity.WasGeneratedBy.Attributes)
	log.Printf("%+v", deliveryTest.Entity.WasGeneratedBy.Edges)
	log.Printf("%+v", string(deliveryTest.Entity.WasGeneratedBy.Data))

	err = json.Unmarshal(rbDelivery, &delivery)
	if err != nil {
		log.Println(err)
		return
	}

	createAgent, createSupervisor, err := tracer.createGraphEntry(&delivery.Entity)
	if err != nil {
		log.Println(err)
		return
	}

	err = tracer.createMongoEntries(&delivery.Entity, createAgent, createSupervisor)
	if err != nil {
		log.Println(err)
		return
	}

}

func (tracer *Tracer) createGraphEntry(e *Entity) (bool, bool, error) {
	var err error
	createAgent, createSupevisor := false, false
	entity := dgraph.NewEntity()

	if err := tracer.fetchUsedEntities(e, entity); err != nil {
		return createAgent, createSupevisor, err
	}

	if createAgent, createSupevisor, err = tracer.fetchAgents(e, entity); err != nil {
		return createAgent, createSupevisor, err
	}

	e.mapStruct(entity, createAgent, createSupevisor)

	assigned, err := tracer.dgraph.AddDerivate(entity)
	if err != nil {
		return createAgent, createSupevisor, err
	}

	log.Printf("inserted dgraph entries as %v\n", assigned)

	e.UID = assigned["entity"]
	e.WasGeneratedBy.UID = assigned["activity"]
	if createAgent {
		e.WasGeneratedBy.WasAssociatedWith.UID = assigned["agent"]
		log.Printf("agent not in database %s, creating entry\n", e.WasGeneratedBy.WasAssociatedWith.ID)
	}
	if createSupevisor {
		e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID = assigned["supervisor"]
		log.Printf("agent not in database %s, creating entry\n", e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID)
	}

	return createAgent, createSupevisor, err
}

func (tracer *Tracer) fetchUsedEntities(e *Entity, entity *dgraph.Entity) error {
	numEntities := len(e.WasGeneratedBy.Used)
	if numEntities > 0 {
		for i, used := range e.WasGeneratedBy.Used {
			usedUID := tracer.mongoDB.EntityUID(used.ID)
			e.WasGeneratedBy.Used[i].UID = usedUID
			entity.WasGeneratedBy.Used = append(entity.WasGeneratedBy.Used, dgraph.NewEntity())
		}
		log.Println("fetched uids of used entities")
	} else {
		log.Println("no used entities to fetch, skipping")
	}
	return nil
}

func (tracer *Tracer) fetchAgents(e *Entity, entity *dgraph.Entity) (bool, bool, error) {
	var err error
	createAgent, createSupevisor := true, true

	agentUID := tracer.mongoDB.AgentUID(e.WasGeneratedBy.WasAssociatedWith.ID)

	if agentUID != "" {
		e.WasGeneratedBy.WasAssociatedWith.UID = agentUID
		createAgent = false
	}

	supervisorUID := tracer.mongoDB.AgentUID(e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID)

	if supervisorUID != "" {
		e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID = supervisorUID
		createSupevisor = false
	}
	return createAgent, createSupevisor, err
}

func (tracer *Tracer) createMongoEntries(e *Entity, createAgent bool, createSupevisor bool) error {
	err := tracer.createMongoEntity(e)
	if err != nil {
		return err
	}

	err = tracer.createMongoActivity(e)
	if err != nil {
		return err
	}

	if createAgent {
		tracer.createMongoAgent(e, false)
		if err != nil {
			return err
		}
	}

	if createSupevisor {
		tracer.createMongoAgent(e, true)
		if err != nil {
			return err
		}
	}

	log.Println("inserted mongodb entries")
	return nil
}

func (tracer *Tracer) createMongoEntity(e *Entity) error {
	mongoEntity := mongodb.NewEntity()
	mongoEntity.UID = e.UID
	mongoEntity.ID = e.ID
	mongoEntity.URI = e.URI
	mongoEntity.Name = e.Name
	mongoEntity.CreationDate = e.CreationDate
	mongoEntity.Data = e.Data

	return tracer.mongoDB.InsertEntity(mongoEntity)
}

func (tracer *Tracer) createMongoAgent(e *Entity, isSupervisor bool) error {
	mongoAgent := mongodb.NewAgent()
	var buffer bytes.Buffer

	if isSupervisor {
		mongoAgent := mongodb.NewAgent()
		mongoAgent.UID = e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID
		mongoAgent.ID = e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID
		mongoAgent.Name = e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.Name

		err := json.Compact(&buffer, e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.Data)

		if err != nil {
			return err
		}

		mongoAgent.Data = buffer.Bytes()
	} else {
		mongoAgent.UID = e.WasGeneratedBy.WasAssociatedWith.UID
		mongoAgent.ID = e.WasGeneratedBy.WasAssociatedWith.ID
		mongoAgent.Name = e.WasGeneratedBy.WasAssociatedWith.Name

		err := json.Compact(&buffer, e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.Data)

		if err != nil {
			return err
		}

		mongoAgent.Data = buffer.Bytes()
	}

	return tracer.mongoDB.InsertAgent(mongoAgent)
}

func (tracer *Tracer) createMongoActivity(e *Entity) error {
	mongoActivity := mongodb.NewActivity()
	mongoActivity.UID = e.WasGeneratedBy.UID
	mongoActivity.ID = e.WasGeneratedBy.ID
	mongoActivity.StartDate = e.WasGeneratedBy.StartDate
	mongoActivity.EndDate = e.WasGeneratedBy.EndDate
	mongoActivity.Data = e.WasGeneratedBy.Data

	return tracer.mongoDB.InsertActivity(mongoActivity)
}

func (e *Entity) mapStruct(entity *dgraph.Entity, createAgent bool, createSupevisor bool) {
	entity.ID = e.ID
	entity.URI = e.URI
	entity.Name = e.Name
	entity.CreationDate = e.CreationDate
	entity.WasGeneratedBy.ID = e.WasGeneratedBy.ID
	entity.WasGeneratedBy.StartDate = e.WasGeneratedBy.StartDate
	entity.WasGeneratedBy.EndDate = e.WasGeneratedBy.EndDate
	entity.WasDerivedFrom = entity.WasGeneratedBy.Used

	if createAgent {
		entity.WasGeneratedBy.WasAssociatedWith.ID = e.WasGeneratedBy.WasAssociatedWith.ID
		entity.WasGeneratedBy.WasAssociatedWith.Name = e.WasGeneratedBy.WasAssociatedWith.Name
	} else {
		entity.WasGeneratedBy.WasAssociatedWith.UID = e.WasGeneratedBy.WasAssociatedWith.UID
	}

	if createSupevisor {
		entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID = e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID
		entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.Name = e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.Name
	} else {
		entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID = e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID
	}
}
