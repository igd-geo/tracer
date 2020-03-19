package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"geocode.igd.fraunhofer.de/hummer/tracer/api"
	"geocode.igd.fraunhofer.de/hummer/tracer/api/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/database"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	configPath := flag.String("configPath", "config.yml", "location of config file")
	flag.Parse()

	config, err := config.Parse(*configPath)
	if err != nil {
		log.Fatalf("failed to parse config: %s", err)
	}

	log.Println("connecting to database")
	database := database.New(config.Database)

	server := api.NewServer(config, database)
	log.Println("starting server")
	go server.Run()

	<-signalChan
}
