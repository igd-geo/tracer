package provutil

import (
	"encoding/json"
)

type Entity struct {
	UID            string           `bson:"uid" json:"uid"`
	ID             string           `bson:"id" json:"id"`
	URI            string           `bson:"uri" json:"uri"`
	Type           string           `bson:"type" json:"type"`
	Name           string           `bson:"name,omitempty" json:"name,omitempty"`
	CreationDate   string           `bson:"creationDate,omitempty" json:"creationDate,omitempty"`
	Data           json.RawMessage  `bson:"data,omitempty" json:"data,omitempty"`
	Graph          *json.RawMessage `bson:"-" json:"graph"`
	WasDerivedFrom []*Entity        `bson:"-" json:"wasDerivedFrom,omitempty"`
	WasGeneratedBy *Activity        `bson:"-" json:"wasGeneratedBy,omitempty"`
}

type Activity struct {
	UID               string           `bson:"uid" json:"uid"`
	ID                string           `bson:"id" json:"id"`
	Type              string           `bson:"type" json:"type"`
	Name              string           `bson:"name,omitempty" json:"name,omitempty"`
	StartDate         string           `bson:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate           string           `bson:"endDate,omitempty" json:"endDate,omitempty"`
	Data              json.RawMessage  `bson:"data,omitempty" json:"data,omitempty"`
	Graph             *json.RawMessage `bson:"-" json:"graph"`
	WasAssociatedWith *Agent           `bson:"-" json:"wasAssociatedWith,omitempty"`
	Used              []*Entity        `bson:"-" json:"used,omitempty"`
}

type Agent struct {
	UID             string           `bson:"uid" json:"uid"`
	ID              string           `bson:"id" json:"id"`
	Type            string           `bson:"type" json:"type"`
	Name            string           `bson:"name,omitempty" json:"name,omitempty"`
	Data            json.RawMessage  `bson:"data,omitempty" json:"data,omitempty"`
	Graph           *json.RawMessage `bson:"-" json:"graph"`
	ActedOnBehalfOf *Agent           `bson:"-" json:"actedOnBehalfOf,omitempty"`
}

func NewEntity() *Entity {
	return &Entity{
		UID: "_:entity",
		WasGeneratedBy: &Activity{
			UID: "_:activity",
			WasAssociatedWith: &Agent{
				UID: "_:agent",
				ActedOnBehalfOf: &Agent{
					UID: "_:supervisor",
				},
			},
			Used: []*Entity{&Entity{}},
		},
	}
}

/*
func NewAttributes() *Attributes {
	return &Attributes{}
}
*/
