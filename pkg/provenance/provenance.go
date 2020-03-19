package provenance

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
	WasGeneratedBy []*Activity     `json:"wasGeneratedBy,omitempty"`
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
	WasAssociatedWith []*Agent        `json:"wasAssociatedWith,omitempty"`
	Used              []*Entity       `json:"used,omitempty"`
}

// Agent is an implementation of PROV-DM agents
type Agent struct {
	UID             string          `json:"uid,omitempty"`
	ID              string          `json:"id,omitempty"`
	Type            string          `json:"type,omitempty"`
	Name            string          `json:"name,omitempty"`
	Description     string          `json:"description,omitempty"`
	Data            json.RawMessage `json:"data,omitempty"`
	Graph           json.RawMessage `json:"graph,omitempty"`
	ActedOnBehalfOf []*Agent        `json:"actedOnBehalfOf,omitempty"`
}

// Graph contains different representations of a provenance graph
type Graph struct {
	json.RawMessage `json:"json,omitempty"`
	Nodes           json.RawMessage `json:"nodes,omitempty"`
	Edges           json.RawMessage `json:"edges,omitempty"`
}

// Node is a map of attributes
type Node map[string]string

// Edge contains information to connect Nodes
type Edge struct {
	Source   string `json:"source,omitempty"`
	Target   string `json:"target,omitempty"`
	EdgeType string `json:"edgeType,omitempty"`
}

// NewEntity returns an empty Entity with initialized Agent and Activity fields
func NewEntity() *Entity {
	return &Entity{}
}
