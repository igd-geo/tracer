package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"geocode.igd.fraunhofer.de/hummer/tracer/api/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/database"

	"github.com/graphql-go/graphql"
)

// Server is the API server
type Server struct {
	config   *config.Config
	database *database.Client
}

type dbConnectionKey string

// NewServer returns a new Server with a given conf and database connection
func NewServer(config *config.Config, database *database.Client) *Server {
	return &Server{
		config:   config,
		database: database,
	}
}

// Run starts the server
func (s *Server) Run() {
	schema := initGraphQL(s.database)

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		w.Header().Set("Content-Type", "application/json")
		result := s.executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	log.Fatal(http.ListenAndServe(s.config.API.Port, nil))
}

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, "+
		"Accept-Encoding, X-CSRF-Token, Authorization")
}

func (s *Server) executeQuery(query string, schema graphql.Schema) *graphql.Result {
	dbKey := dbConnectionKey("db")

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       context.WithValue(context.Background(), dbKey, s.database),
	})
	if len(result.Errors) > 0 {
		log.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
