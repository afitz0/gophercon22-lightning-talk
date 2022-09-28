package api_client

import (
	"lightning/app"
	"lightning/app/constants"
)

type ApiClient struct{}
type RetryableError bool

type FulfillResult struct {
	Order  app.Order
	Status string
	Error  error
}

func New() (ApiClient, error) {
	return ApiClient{}, nil
}

func (a *ApiClient) Close() {}

func (a *ApiClient) CreateOrder(req app.Order) error {
	return nil
}

func (a *ApiClient) FulfillOrder(req app.Order) <-chan FulfillResult {
	c := make(chan FulfillResult)
	go func() {
		c <- FulfillResult{
			Order:  req,
			Status: constants.ORDER_FULFILLED,
			Error:  nil,
		}
	}()
	return c

	//return app.OrderStatus{}, nil
}

func (r RetryableError) Error() string {
	return "An error occured during the API call that is retryable"
}
