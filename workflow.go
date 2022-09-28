package app

import (
	"context"
	"math/rand"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Order struct {
	Id string
}
type OrderRequest struct{}

func Workflow(ctx workflow.Context, o Order) error {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2,
		MaximumInterval:        10 * time.Second,
		MaximumAttempts:        100,
		NonRetryableErrorTypes: []string{"PaymentFailed"},
	}
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		RetryPolicy:         retryPolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	log := workflow.GetLogger(ctx)

	err := workflow.ExecuteActivity(ctx, CreateOrder, o).Get(ctx, &o)
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

	return nil
}

func CreateOrder(ctx context.Context, o Order) (Order, error) {
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	return o, nil
}

func FulfillOrder(ctx context.Context, o Order) (string, error) {
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	return "", nil
}

func ArchiveOrder(ctx context.Context, o Order, s string) error {
	time.Sleep(time.Second * time.Duration(rand.Intn(10)))
	return nil
}
