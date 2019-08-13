package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/provutil"
	"github.com/graphql-go/graphql"
)

type Server struct {
	config *config.Config
	provDB provutil.ProvDB
	infoDB provutil.InfoDB
}

func NewServer(config *config.Config, infoDB provutil.InfoDB, provDB provutil.ProvDB) *Server {
	return &Server{
		config: config,
		infoDB: infoDB,
		provDB: provDB,
	}
}

func (s *Server) Run() {

	schema := initGraphQL(s.infoDB, s.provDB)

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	log.Fatal(http.ListenAndServe(":1234", nil))
}

func (s *Server) Cleanup() error {
	return nil
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
