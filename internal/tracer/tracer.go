package tracer

import (
	"encoding/json"
	"fmt"
	"log"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/dgraph"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/mongodb"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rabbitmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
)

type Tracer struct {
	deliveries <-chan rabbitmq.Delivery
	consumer   *rabbitmq.Consumer
	mongoDB    *mongodb.Client
	dgraph     *dgraph.Client
	config     *config.Config
}

type Entity struct {
	UID            string          `json:"uid,omitempty"`
	ID             string          `json:"id,omitempty"`
	URI            string          `json:"uri,omitempty"`
	Name           string          `json:"name,omitempty"`
	CreationDate   string          `json:"creationDate,omitempty"`
	Optional       json.RawMessage `json:"optional,omitempty"`
	WasGeneratedBy struct {
		UID               string          `json:"uid,omitempty"`
		ID                string          `json:"id,omitempty"`
		StartDate         string          `json:"startDate,omitempty"`
		EndDate           string          `json:"endDate,omitempty"`
		Optional          json.RawMessage `json:"optional,omitempty"`
		WasAssociatedWith struct {
			UID             string          `json:"uid,omitempty"`
			ID              string          `json:"id,omitempty"`
			Name            string          `json:"name,omitempty"`
			Optional        json.RawMessage `json:"optional,omitempty"`
			ActedOnBehalfOf struct {
				UID      string          `json:"uid,omitempty"`
				ID       string          `json:"id,omitempty"`
				Name     string          `json:"name,omitempty"`
				Optional json.RawMessage `json:"optional,omitempty"`
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

type uids struct {
	used       []string
	agent      string
	supervisor string
}

func New(config *config.Config) *Tracer {
	msgChan := make(chan rabbitmq.Delivery)
	tracer := Tracer{
		deliveries: msgChan,
		consumer:   rabbitmq.NewConsumer(config.RabbitURL, msgChan, config.ConsumerTag),
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
	return tracer.consumer.Shutdown()
}

func (tracer *Tracer) handleDelivery(delivery rabbitmq.Delivery) {
	var entity Entity
	err := json.Unmarshal(delivery, &entity)
	if err != nil {
		log.Println(err)
		return
	}

	//fmt.Println(dgraph.Derivate(entity))
}

func (tracer *Tracer) createGraphEntry(e *Entity) error {
	var err error
	createAgent, createSupevisor := false, false
	entity := dgraph.NewEntity("derivate")
	entity.WasGeneratedBy.UID = "activity"

	for i, used := range e.WasGeneratedBy.Used {
		usedUID, err := tracer.mongoDB.EntityUID(used.ID)
		if err != nil {
			return err
		}
		e.WasGeneratedBy.Used[i].UID = usedUID
		entity.WasGeneratedBy.Used = append(entity.WasGeneratedBy.Used, dgraph.NewEntity(usedUID))
	}

	e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID, err = tracer.mongoDB.AgentUID(e.WasGeneratedBy.WasAssociatedWith.ID)
	if err != nil {
		return err
	}

	if e.WasGeneratedBy.WasAssociatedWith.UID != "" {
		e.WasGeneratedBy.WasAssociatedWith.UID = "agent"
		createAgent = true
	}

	e.WasGeneratedBy.WasAssociatedWith.UID, err = tracer.mongoDB.AgentUID(e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID)
	if err != nil {
		return err
	}

	if e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID != "" {
		e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID = "supervisor"
		createSupevisor = true
	}

	e.mapStruct(entity, createAgent, createSupevisor)

	assigned, err := tracer.dgraph.AddDerivate(entity)
	fmt.Println(assigned)

	e.UID = assigned["derivate"]
	e.WasGeneratedBy.UID = assigned["activity"]
	if createAgent {
		e.WasGeneratedBy.WasAssociatedWith.UID = assigned["agent"]
	}
	if createSupevisor {
		e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID = assigned["supervisor"]
	}

	return err
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
	entity.WasGeneratedBy.WasAssociatedWith.UID = e.WasGeneratedBy.WasAssociatedWith.UID
	entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID = e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.UID

	if createAgent {
		entity.WasGeneratedBy.WasAssociatedWith.ID = e.WasGeneratedBy.WasAssociatedWith.ID
		entity.WasGeneratedBy.WasAssociatedWith.Name = e.WasGeneratedBy.WasAssociatedWith.Name
	}

	if createSupevisor {
		entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID = e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.ID
		entity.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.Name = e.WasGeneratedBy.WasAssociatedWith.ActedOnBehalfOf.Name
	}
}

func generateEntry()    {}
func generateMutation() {}
