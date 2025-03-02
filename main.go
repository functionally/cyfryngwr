package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/functionally/cyfryngwr/cwtch"
	"github.com/functionally/cyfryngwr/dispatch"
)

func main() {

	config := make(map[string]interface{})
	dispatcher, err := dispatch.New(config)
	if err != nil {
		log.Fatalf("Failed to create dispatcher: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	cwtchbot := cwtch.Connect(".cyfryngwr/", "cyfryngwr", "Cyfryngwr, a cwtch agent")
	cwtch.Loop(dispatcher, cwtchbot, stop)

}
