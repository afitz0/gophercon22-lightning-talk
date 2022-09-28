package main

import (
	"lightning/app"
	"lightning/app/constants"
	"lightning/app/plain/api_client"
	"lightning/app/plain/archival"
	"lightning/app/plain/dist_store"
	"lightning/app/plain/log"

	"errors"
	"fmt"
	"time"
)

func OrderItem(o app.Order) error {
	const MAX_ATTEMPTS = 50
	const BACKOFF_COEFFICIENT = 2
	const MAX_TIMEOUT = 300 * time.Second
	const INITIAL_TIMEOUT = 1 * time.Second

	c, err := api_client.New()
	if err != nil {
		log.Error("Failed to create new API client", err)
		return err
	}
	defer c.Close()

	s, ok := dist_store.Get(o.Id)

	if !ok || s == constants.ORDER_RECIEVED {
		dist_store.Set(o.Id, constants.ORDER_RECIEVED)
		retryDelay := INITIAL_TIMEOUT
		for try := 0; try < MAX_ATTEMPTS; try++ {
			err = c.CreateOrder(o)
			if err != nil {
				log.Error("Failed to create new order", err)
			}

			if !err.(api_client.RetryableError) {
				log.Fatal("Got unretryable error from API call. Crashing.", err)
			}

			// exponential backoff
			time.Sleep(retryDelay)
			retryDelay *= BACKOFF_COEFFICIENT
			if retryDelay > MAX_TIMEOUT {
				retryDelay = MAX_TIMEOUT
			}
		}
		dist_store.Set(o.Id, constants.ORDER_PLACED)
	}

	fulfillment := c.FulfillOrder(o)

	// push the id and current status to the distributed store, for resumability
	dist_store.Set(o.Id, constants.ORDER_INPROGRESS)

	select {
	case res := <-fulfillment:
		switch res.Error {
		case nil:
			dist_store.Set(o.Id, res.Status)
			break
		default:
			log.Error("Error from fulfillment", res.Error)
		}
		fmt.Print("Received result")
	case <-time.After(MAX_TIMEOUT):
		log.Fatal("Fulfillment timed out. Crashing.", errors.New("error"))
	}

	db, err := archival.NewClient()
	if err != nil {
		log.Error("Failed to create archival client", err)
		return err
	}
	defer db.Close()

	err = db.Persist(order, status)
	if err != nil {
		log.Error("Failed to persist order details", err)
		return err
	}

	return nil
}
