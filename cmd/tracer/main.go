package main

import (
	"log"
	"os"
	"os/signal"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/dgraph"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/mongodb"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	config := config.New()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	installFlags(config)

	infoDB := mongodb.NewClient(
		config.MongoURL,
		config.MongoDatabase,
		config.MongoCollectionEntity,
		config.MongoCollectionAgent,
		config.MongoCollectionActivity,
	)
	provDB := dgraph.NewClient(config.DgraphURL)

	tracer := tracer.New(config, infoDB, provDB)
	tracer.Listen()

	<-signalChan
	err := tracer.Cleanup()
	if err != nil {
		log.Fatal(err)
	}
}
