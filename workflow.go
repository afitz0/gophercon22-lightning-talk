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

	err := workflow.ExecuteActivity(ctx, InitOrder, o).Get(ctx, nil)
	if err != nil {
		log.Error("InitOrder failed", "Err", err)
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
	//fmt.Println("No more. :(")
	return nil
}

func InitOrder(ctx context.Context, o common.Order) error {
	//fmt.Println("Hello GopherCon!")
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
