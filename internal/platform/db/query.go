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

	QueryAllUIDsByID = `query All($entity: string, $activity: string, $agent: string, $supervisor: string) {
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

	QueryEntityFullByID = `query Entity($entity: string) {
		entity(func: eq(id, $entity)) {
			uid
			id
			creationDate
		}
	}`

	QueryActivityFullByID = `
	query Activity($activity: string) {
		activity(func: eq(id, $activity)) {
			uid
			id
			startDate
			endDate
		}
	}`

	QueryAgentFullByID = `
	query Agent($agent: string) {
		agent(func: eq(id, $agent)) {
			uid
			id
			name
			description
			type
		}
	}`

	QueryWasGeneratedBy = `
	query WasGeneratedBy($entity: string) {
		var(func: uid($entity)) {
			A as wasGeneratedBy {
			  uid
			}
		}

		activity (func: uid(A)) {
			uid
			id
			startDate
			endDate
	 	}
	}`

	QueryWasDerivedFrom = `
	query WasDerivedFrom($entity: string) {
		var(func: uid($entity)) {
			A as wasDerivedFrom {
			  uid
			}
		}

		entity (func: uid(A)) {
			uid
			id
			creationDate
	 	}
	}`

	QueryWasAssociatedWith = `
	query WasAssociatedWith($activity: string) {
		var(func: uid($activity)) {
			A as wasAssociatedWith {
			  uid
			}
		}

		agent (func: uid(A)) {
			uid
			id
			name
			description
			type
	 	}
	}`

	QueryActedOnBehalfOf = `
	query ActedOnBehalfOf($agent: string) {
		var(func: uid($agent)) {
			A as actedOnBehalfOf {
			  uid
			}
		}

		agent (func: uid(A)) @filter(has(name)){
			uid
			id
			name
			description
			type
	 	}
	}`

	QueryUsed = `
	query Used($activity: string) {
		var(func: uid($activity)) {
			A as used {
			  uid
			}
		}

		entity (func: uid(A)) {
			uid
			id
			creationDate
	 	}
	}`

	QueryProvenanceGraph = `
	query Graph($root: string) {
		graph(func: eq(id,$root)) {
			id
			wasDerivedFrom {
			  	id
			}
			wasGeneratedBy {
			  	id
			  	wasAssociatedWith {
					id
					actedOnBehalfOf {
				  		id
					}
			  	}
			  	used {
					id
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
