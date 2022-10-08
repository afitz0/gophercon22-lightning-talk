package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"lightning/app"
	"lightning/app/zapadapter"
)

func main() {
	c, err := client.Dial(client.Options{
		Logger: zapadapter.NewZapAdapter(
			zapadapter.NewZapLogger()),
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "lightning", worker.Options{})

	w.RegisterWorkflow(app.Workflow)
	w.RegisterActivity(app.InitOrder)
	w.RegisterActivity(app.FulfillOrder)
	w.RegisterActivity(app.ArchiveOrder)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
