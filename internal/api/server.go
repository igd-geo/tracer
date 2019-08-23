package api

import (
	"encoding/json"
	"log"
	"net/http"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"github.com/graphql-go/graphql"
)

type Server struct {
	config *config.Config
	db     *db.Client
}

func NewServer(config *config.Config, db *db.Client) *Server {
	return &Server{
		config: config,
		db:     db,
	}
}

func (s *Server) Run() {
	schema := initGraphQL(s.db)

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	log.Fatal(http.ListenAndServe(":1234", nil))
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		log.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
