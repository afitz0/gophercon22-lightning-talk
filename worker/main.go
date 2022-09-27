package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"lightning/app"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "lightning", worker.Options{})

	w.RegisterWorkflow(app.Workflow)
	w.RegisterActivity(app.CreateOrder)
	w.RegisterActivity(app.FulfillOrder)
	w.RegisterActivity(app.ArchiveOrder)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
