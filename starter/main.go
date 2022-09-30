package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	"lightning/app"
	"lightning/app/common"
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

	workflowOptions := client.StartWorkflowOptions{
		ID:        "gophercon_lightning_talk_workflowID",
		TaskQueue: "lightning",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, app.Workflow, common.Order{})
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
