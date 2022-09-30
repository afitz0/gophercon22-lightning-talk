package app

import (
	"context"
	"fmt"
	golog "log"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"lightning/app/common"
	"lightning/app/mocks/api_client"
	"lightning/app/mocks/archival"
)

func Workflow(ctx workflow.Context, o common.Order) error {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2,
		MaximumInterval:        10 * time.Second,
		MaximumAttempts:        100,
		NonRetryableErrorTypes: []string{"PaymentFailed"},
	}
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy:         retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	log := workflow.GetLogger(ctx)

	err := workflow.ExecuteActivity(ctx, CreateOrder, o).Get(ctx, nil)
	if err != nil {
		log.Error("CreateOrder failed", "Err", err)
		return err
	}

	var status string
	err = workflow.ExecuteActivity(ctx, FulfillOrder, o).Get(ctx, &status)
	if err != nil {
		log.Error("FulfillOrder failed", "Err", err)
		return err
	}

	err = workflow.ExecuteActivity(ctx, ArchiveOrder, o, status).Get(ctx, nil)
	if err != nil {
		log.Error("ArchiveOrder failed", "Err", err)
		return err
	}

	fmt.Println("Done!")
	return nil
}

func CreateOrder(ctx context.Context, o common.Order) error {
	c, err := api_client.New()
	if err != nil {
		golog.Println("Could not create api client", err)
		return err
	}
	defer c.Close()

	err = c.InitOrder(o)
	return err
}

func FulfillOrder(ctx context.Context, o common.Order) (string, error) {
	c, err := api_client.New()
	if err != nil {
		golog.Println("Could not create api client", err)
		return "", err
	}
	defer c.Close()

	status, err := c.FulfillOrder(o)
	return status, err
}

func ArchiveOrder(ctx context.Context, o common.Order, s string) error {
	db, err := archival.NewClient()
	if err != nil {
		golog.Println("Could not create api client", err)
		return err
	}
	defer db.Close()

	err = db.Persist(o, s)
	return err
}
