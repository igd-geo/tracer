package util

import (
	"encoding/json"
)

// Entity is an implementation of PROV-DM entities
type Entity struct {
	UID            string          `json:"uid,omitempty"`
	ID             string          `json:"id,omitempty"`
	Type           string          `json:"type,omitempty"`
	URI            string          `json:"uri,omitempty"`
	Name           string          `json:"name,omitempty"`
	CreationDate   string          `json:"creationDate,omitempty"`
	Data           json.RawMessage `json:"data,omitempty"`
	Graph          json.RawMessage `json:"graph,omitempty"`
	WasDerivedFrom []*Entity       `json:"wasDerivedFrom,omitempty"`
	WasGeneratedBy *Activity       `json:"wasGeneratedBy,omitempty"`
}

// Activity is an implementation of PROV-DM activities
type Activity struct {
	UID               string          `json:"uid,omitempty"`
	ID                string          `json:"id,omitempty"`
	Type              string          `json:"type,omitempty"`
	IsBatch           bool            `json:"isBatch,omitempty"`
	Name              string          `json:"name,omitempty"`
	StartDate         string          `json:"startDate,omitempty"`
	EndDate           string          `json:"endDate,omitempty"`
	Data              json.RawMessage `json:"data,omitempty"`
	Graph             json.RawMessage `json:"graph,omitempty"`
	WasAssociatedWith *Agent          `json:"wasAssociatedWith,omitempty"`
	Used              []*Entity       `json:"used,omitempty"`
}

// Agent is an implementation of PROV-DM agents
type Agent struct {
	UID             string          `json:"uid,omitempty"`
	ID              string          `json:"id,omitempty"`
	Type            string          `json:"type,omitempty"`
	Name            string          `json:"name,omitempty"`
	Test            string          `json:"test,omitempty"`
	Data            json.RawMessage `json:"data,omitempty"`
	Graph           json.RawMessage `json:"graph,omitempty"`
	ActedOnBehalfOf *Agent          `json:"actedOnBehalfOf,omitempty"`
}

type Graph struct {
	json.RawMessage `json:"structure,omitempty"`
}

// NewEntity returns an empty Entity with initialized Agent and Activity fields
func NewEntity() *Entity {
	return &Entity{
		WasGeneratedBy: &Activity{
			WasAssociatedWith: &Agent{
				ActedOnBehalfOf: &Agent{},
			},
		},
	}
}
