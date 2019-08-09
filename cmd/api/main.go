package main

import (
	"log"
	"os"
	"os/signal"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/api/config"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	config := config.New()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	installFlags(config)

	server := api.NewServer(config)
	go server.Run()

	<-signalChan
	err := server.Cleanup()
	if err != nil {
		log.Fatal(err)
	}
}
