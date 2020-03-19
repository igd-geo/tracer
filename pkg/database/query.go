package database

const (
	// QueryEntityUIDByID GraphQL query to fetch an entity's uid
	QueryEntityUIDByID = `query Entity($entity: string) {
		entity(func: eq(id, $entity)) {
			uid
		}
	}`

	// QueryActivityUIDByID GraqhQL query to fetch an activity's uid
	QueryActivityUIDByID = `query Activity($activity: string) {
		activity(func: eq(id, $activity)) {
			uid
		}
	}`

	// QueryAgentUIDByID GraphQL query to fetch an agend's uid
	QueryAgentUIDByID = `query Agent($agent: string) {
		agent(func: eq(id, $agent)) {
			uid
		}
	}`

	// QuerySupervisorUIDByID GraphQL query to fetch an agend's uid
	QuerySupervisorUIDByID = `query Supervisor($supervisor: string) {
		supervisor(func: eq(id, $supervisor)) {
			uid
		}
	}`

	// QueryAllUIDsByID Combined GraphQL query to fetch the uids of several provenance components
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

	// QueryEntityFullByID GraphQL query to fetch an entity's attributes
	QueryEntityFullByID = `query Entity($entity: string) {
		entity(func: eq(id, $entity)) {
			uid
			id
			creationDate
		}
	}`

	// QueryActivityFullByID GraphQL query to fetch an activity's attributes
	QueryActivityFullByID = `
	query Activity($activity: string) {
		activity(func: eq(id, $activity)) {
			uid
			id
			startDate
			endDate
		}
	}`

	// QueryAgentFullByID GraphQL query to fetch an agent's attirbutes
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

	// QueryWasGeneratedBy GraphQL query to fetch the activity that generated an entity
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

	// QueryWasDerivedFrom GraphQL query to fetch an entity's ancestor
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

	// QueryWasAssociatedWith GraphQL query to fetch an activities acting agent
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

	// QueryActedOnBehalfOf GraphQL query to fetch an agents supervisor
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

	// QueryUsed GraphQL query to fetch an activity's used entities
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

	// QueryProvenanceGraph GraphQL query a complete provenance graph
	// structure starting with an entity as root
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

	// VariableEntityID Query variable for prepared query statements
	VariableEntityID = "$entity"
	// VariableActivityID Query variable for prepared query statements
	VariableActivityID = "$activity"
	// VariableAgentID Query variable for prepared query statements
	VariableAgentID = "$agent"
	// VariableSupervisorID Query variable for prepared query statements
	VariableSupervisorID = "$supervisor"
	// VariableGraphRootID Query variable for prepared query statements
	VariableGraphRootID = "$root"
)

// Query Wrapper for GraphQL prepared query statements
type Query struct {
	queryString string
	variables   map[string]string
}

// NewQuery Returns a new Query
func NewQuery(queryString string) *Query {
	return &Query{
		queryString: queryString,
		variables:   make(map[string]string),
	}
}

// SetVariable Adds query variables to be used with prepared statements
func (q *Query) SetVariable(key string, value string) {
	q.variables[key] = value
}
