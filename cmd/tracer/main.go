package main

import (
	"log"
	"os"
	"os/signal"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/broker"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	config := config.New()

	deliveries := make(chan *util.Entity)
	broker := broker.New(config.Broker, deliveries, "notifications", "topic")
	db := db.NewClient(config.DB)

	tracer := tracer.Setup(config, db, broker, deliveries)
	tracer.Listen()

	<-signalChan
	err := tracer.Cleanup()
	if err != nil {
		log.Fatal(err)
	}
}
