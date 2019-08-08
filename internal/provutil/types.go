package provutil

import (
	"encoding/json"
)

type Entity struct {
	*Attributes
	*Edges `bson:"-"`
	Data   json.RawMessage `json:"data,omitempty" bson:"data,omitempty"`
}

type Activity struct {
	*Attributes
	*Edges `bson:"-"`
	Data   json.RawMessage `json:"data,omitempty" bson:"data,omitempty"`
}

type Agent struct {
	*Attributes
	*Edges `bson:"-"`
	Data   json.RawMessage `json:"data,omitempty" bson:"data,omitempty"`
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

func NewEntity() *Entity {
	return &Entity{
		Attributes: &Attributes{
			UID: "_:entity",
		},
		Edges: &Edges{
			WasGeneratedBy: &Activity{
				Attributes: &Attributes{
					UID: "_:activity",
				},
				Edges: &Edges{
					WasAssociatedWith: &Agent{
						Attributes: &Attributes{
							UID: "_:agent",
						},
						Edges: &Edges{
							ActedOnBehalfOf: &Agent{
								Attributes: &Attributes{
									UID: "_:supervisor",
								},
							},
						},
					},
					Used: []*Entity{&Entity{}},
				},
			},
		},
	}
}

func NewAttributes() *Attributes {
	return &Attributes{}
}
