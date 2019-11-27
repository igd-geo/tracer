package main

import (
	"log"
	"os"
	"os/signal"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rbmq"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	conf := config.New()

	msgBroker := rbmq.NewBroker(conf.Broker)

	httpDummySession := msgBroker.NewProducerOnlySession("notifications", "topic")

	db := db.NewClient(conf.DB)

	server := api.NewServer(conf, db, httpDummySession)
	go server.Run()

	<-signalChan
}
