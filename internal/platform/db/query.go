package db

const (
	QueryEntityUIDByID = `query Entity($entity: string) {
							entity(func: eq(id, $entity)) {
								uid
							}
						}`

	QueryActivityUIDByID = `query Activity($activity: string) {
							activity(func: eq(id, $activity)) {
								uid
							}
						}`

	QueryAgentUIDByID = `query Agent($agent: string) {
							agent(func: eq(id, $agent)) {
								uid
							}
						}`

	QuerySupervisorUIDByID = `query Supervisor($supervisor: string) {
								supervisor(func: eq(id, $supervisor)) {
									uid
								}
							}`

	QueryAllUIDsByID = `query All($entity: string, $activity: string, agent: string, supervisor: string) {
							entity(func: eq(id, $entity)) {
								uid
							}
							activity(func: eq(id, $activity)) {
								uid
							}
							agent(func: eq(id, $agent)) {
								uid
							}
							supervisor(func: eq(id, $supervisor)) {
								uid
							}
						}`

	QueryProvenanceGraph = `query Graph($root: string) {
								graph(func: eq(id, $root)) {
									expand(_all_) {
										expand(_all_) {
											expand(_all_) {
												expand(_all_)
											}
										}
									}
								}
							}`

	VariableEntityID     = "$entity"
	VariableActivityID   = "$activity"
	VariableAgentID      = "$agent"
	VariableSupervisorID = "$supervisor"
	VariableGraphRootID  = "$root"
)

type Query struct {
	queryString string
	variables   map[string]string
}

func NewQuery(queryString string) *Query {
	return &Query{
		queryString: queryString,
		variables:   make(map[string]string),
	}
}

func (q *Query) SetVariable(key string, value string) {
	q.variables[key] = value
}
