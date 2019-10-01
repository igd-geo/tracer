package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
	"github.com/graphql-go/graphql"
)

// Server is the API server
type Server struct {
	config *config.Config
	db     *db.Client
}

// NewServer returns a new Server with a given config and database connection
func NewServer(config *config.Config, db *db.Client) *Server {
	return &Server{
		config: config,
		db:     db,
	}
}

// Run starts the server
func (s *Server) Run() {
	schema := initGraphQL(s.db)

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/metadata", func(w http.ResponseWriter, r *http.Request) {
		metadata := fetchMetadata("https://inspire-geoportal.ec.europa.eu/resources/INSPIRE-c1e5f7f2-3b35-11e9-a83c-52540023a883_20190902-141544/services/1/PullResults/33741-33760/1.iso19139.xml")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/xml")
		w.Write(metadata)
	})

	log.Fatal(http.ListenAndServe(s.config.Port, nil))
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

func fetchMetadata(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	util.ParseMetadataToJSON(body)
	return body
}
