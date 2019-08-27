package main

import (
	"log"
	"os"
	"os/signal"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api/config"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/platform/db"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	config := config.New()

	db := db.NewClient(config.DB)

	server := api.NewServer(config, db)
	go server.Run()

	<-signalChan
}
