package provutil

import "encoding/json"

type InfoDB interface {
	EntityUID(id string) string
	AgentUID(id string) string
	ActivitytUID(id string) string
	FetchEntity(id string) *Entity
	FetchAgent(id string) *Agent
	FetchActivity(id string) *Activity
	InsertEntity(entity *Entity) error
	InsertAgent(agent *Agent) error
	InsertActivity(activity *Activity) error
}

type ProvDB interface {
	InsertDerivate(derivate *Entity) (map[string]string, error)
	FetchProvenanceGraph(uid string) *json.RawMessage
}
