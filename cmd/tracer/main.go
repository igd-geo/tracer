package main

import (
	"log"
	"os"
	"os/signal"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/rbmq"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/util"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	conf := config.New()

	deliveries := make(chan *util.Entity)

	msgBroker := rbmq.NewBroker(conf.Broker)

	provSession := msgBroker.NewSession(deliveries, "notifications", "topic")
	db := db.NewClient(conf.DB)

	tracer := tracer.Setup(conf, db, provSession, deliveries)
	tracer.Listen()

	<-signalChan
	err := tracer.Cleanup()
	if err != nil {
		log.Fatal(err)
	}
}
