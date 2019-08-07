package types

import (
	"encoding/json"
)

type Entity struct {
	*Attributes
	*Edges
	Data json.RawMessage `bson:"data,omitempty"`
}

type Activity struct {
	*Attributes
	*Edges
	Data json.RawMessage `bson:"data,omitempty"`
}

type Agent struct {
	*Attributes
	*Edges
	Data json.RawMessage `bson:"data,omitempty"`
}

type Attributes struct {
	UID          string `bson:"uid,omitempty" json:"uid,omitempty"`
	ID           string `bson:"id,omitempty" json:"id,omitempty"`
	URI          string `bson:"uri,omitempty" json:"uri,omitempty"`
	Name         string `bson:"name,omitempty" json:"name,omitempty"`
	CreationDate string `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
	StartDate    string `bson:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate      string `bson:"endDate,omitempty" json:"endDate,omitempty"`
}

type Edges struct {
	WasDerivedFrom    []*Entity `json:"wasDerivedFrom,omitempty"`
	WasGeneratedBy    *Activity `json:"wasGeneratedBy,omitempty"`
	WasAssociatedWith *Agent    `json:"wasAssociatedWith,omitempty"`
	Used              []*Entity `json:"used,omitempty"`
	ActedOnBehalfOf   *Agent    `json:"actedOnBehalfOf,omitempty"`
}

type Data struct {
	json.RawMessage `bson:"data,omitempty"`
}

func NewEntity() *Entity {
	return &Entity{Attributes: &Attributes{}, Edges: &Edges{WasGeneratedBy: newActivity()}}
}

func newActivity() *Activity {
	return &Activity{Attributes: &Attributes{}, Edges: &Edges{WasAssociatedWith: newAgent()}}
}

func newAgent() *Agent {
	return &Agent{Attributes: &Attributes{}, Edges: &Edges{ActedOnBehalfOf: newAgent()}}
}
