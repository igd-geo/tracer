package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rbmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
	"github.com/graphql-go/graphql"
)

// Server is the API server
type Server struct {
	conf *config.Config
	db   *db.Client
	rbmq *rbmq.Session
}

// NewServer returns a new Server with a given conf and database connection
func NewServer(conf *config.Config, db *db.Client, rbSession *rbmq.Session) *Server {
	return &Server{
		conf: conf,
		db:   db,
		rbmq: rbSession,
	}
}

// Run starts the server
func (s *Server) Run() {
	schema := initGraphQL(s.db)

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		w.Header().Set("Content-Type", "application/json")
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	http.HandleFunc("/metadata", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		metadata := fetchMetadata("https://inspire-geoportal.ec.europa.eu/resources/" +
			"INSPIRE-c1e5f7f2-3b35-11e9-a83c-52540023a883_20190902-141544/services/1/PullResults/" +
			"33741-33760/1.iso19139.xml")
		w.Header().Set("Content-Type", "application/xml")
		w.Write(metadata)
	})

	http.HandleFunc("/httpdummy", func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		msg, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.rbmq.Publish(string(msg), "tracer")
		log.Println(string(msg))
	})

	log.Fatal(http.ListenAndServe(s.conf.Port, nil))
}

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, " +
		"Accept-Encoding, X-CSRF-Token, Authorization")
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
