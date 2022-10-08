package api_client

import (
	"math/rand"

	"lightning/app/common"
	"lightning/app/plain/constants"
)

type ApiClient struct{}
type RetryableError bool

type FulfillResult struct {
	Order  common.Order
	Status string
	Error  error
}

func New() (ApiClient, error) {
	return ApiClient{}, nil
}

func (a *ApiClient) Close() {}

func (a *ApiClient) InitOrder(req common.Order) error {
	common.Sleep(rand.Intn(5), "InitOrder")
	return nil
}

func (a *ApiClient) FulfillOrder(req common.Order) (string, error) {
	common.Sleep(8, "FulfillOrder")
	return constants.ORDER_FULFILLED, nil
}

func (a *ApiClient) FulfillOrderChan(req common.Order) <-chan FulfillResult {
	c := make(chan FulfillResult)
	go func() {
		common.Sleep(rand.Intn(10), "FulfillOrder")
		c <- FulfillResult{
			Order:  req,
			Status: constants.ORDER_FULFILLED,
			Error:  nil,
		}
	}()
	return c
}

func (r RetryableError) Error() string {
	return "An error occured during the API call that is retryable"
}
