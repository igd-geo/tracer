package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"geocode.igd.fraunhofer.de/hummer/tracer/client"
	"geocode.igd.fraunhofer.de/hummer/tracer/client/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/pkg/broker"
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

	brokerSession, err := broker.New(config.Broker)
	if err != nil {
		log.Fatalf("failed to establish connection to the message broker: %s", err)
	}

	database := database.New(config.Database)

	client := client.New(config, database, brokerSession)
	if err := client.Listen(); err != nil {
		log.Fatal(err)
	}

	<-signalChan
	err = client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
