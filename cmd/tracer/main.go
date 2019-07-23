package main

import (
	"log"
	"os"
	"os/signal"

	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer"
	"geocode.igd.fraunhofer.de/hummer/tracer/internal/tracer/config"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	config := config.New()
	installFlags(config)

	tracer := tracer.New(config)
	tracer.Listen()

	<-signalChan
	err := tracer.Cleanup()
	if err != nil {
		log.Fatal(err)
	}
}
