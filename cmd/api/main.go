package main

import (
	"log"
	"os"
	"os/signal"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/dgraph"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/mongodb"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	config := config.New()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	installFlags(config)

	infoDB := mongodb.NewClient(config.InfoDB)
	provDB := dgraph.NewClient(config.ProvDB)

	server := api.NewServer(config, infoDB, provDB)
	go server.Run()

	<-signalChan
	err := server.Cleanup()
	if err != nil {
		log.Fatal(err)
	}
}
