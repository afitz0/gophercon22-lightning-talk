package main

import (
	"lightning/app"
	"lightning/app/plain/api_client"
	"lightning/app/plain/archival"
	"lightning/app/plain/constants"
	"lightning/app/plain/dist_store"
	"lightning/app/plain/log"

	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/slices"
)

const (
	BACKOFF_COEFFICIENT = 2
	MAX_TIMEOUT         = 300 * time.Second
	MAX_ATTEMPTS        = 50
	INITIAL_TIMEOUT     = 1 * time.Second
)

var (
	retryableFulfillmentStatuses = []string{
		constants.E_ORDER_LOST,
		constants.E_ORDER_DROPPED,
		constants.E_TIMEOUT,
	}
)

func OrderItem(o app.Order) error {
	c, err := api_client.New()
	if err != nil {
		log.Fatal("Failed to create new API client", err)
		return err
	}
	defer c.Close()

	err = dist_store.InitClient()
	if err != nil {
		log.Fatal("Failed to create new dist store client", err)
		return err
	}
	defer dist_store.Close()

	// First check where we currently are, so that we can skip ahead if necessary
	s, ok := dist_store.Get(o.Id)
	log.Debug("Got status from dist_store:", s)

	if !ok || s == constants.ORDER_RECIEVED {
		log.Info("New order, attempting to place.")
		dist_store.Set(o.Id, constants.ORDER_RECIEVED)
		retryDelay := INITIAL_TIMEOUT
		success := false
		for try := 0; try < MAX_ATTEMPTS && !success; try++ {
			err = c.CreateOrder(o)
			if err != nil && !err.(api_client.RetryableError) {
				log.Fatal("Unretryable error from trying to create order. Crashing.", err)
			} else if err == nil {
				log.Info("Successfully placed order.")
				success = true
			} else {
				log.Error("Placing order failed. Retrying after", retryDelay, "s. Error:", err)
				retryDelay = expoBackoff(retryDelay)
			}
		}

		// here because we exhausted retry budget? gotta die
		if !success {
			log.Fatal("Exhausted retries trying to create order. Crashing. Last error:", err)
		}

		dist_store.Set(o.Id, constants.ORDER_PLACED)
	}

	// Update status in case we got here from a retry
	s, ok = dist_store.Get(o.Id)
	log.Debug("Got status from dist_store:", s)

	// The valid states that we can initiate fulfillment from
	if s == constants.ORDER_RECIEVED || s == constants.ORDER_PLACED {
		log.Info("Attempting to fulfill order.")
		retryDelay := INITIAL_TIMEOUT

		// push the id and current status to the distributed store, for resumability
		dist_store.Set(o.Id, constants.ORDER_INPROGRESS)

		success := false
		for try := 0; try < MAX_ATTEMPTS && !success; try++ {
			fulfillment := c.FulfillOrder(o)

			select {
			case res := <-fulfillment:
				if res.Error != nil && res.Status != constants.E_ORDER_DUPLICATE {
					log.Error("Error from fulfillment", res.Error)
					if slices.Contains(retryableFulfillmentStatuses, res.Status) {
						// retryable? yep, so let's do it again after a delay
						log.Info("Retrying fulfillment after", retryDelay, "seconds")
						retryDelay = expoBackoff(retryDelay)
						break
					} else {
						log.Fatal("Fulfillment hit unrecoverable error. Crashing.", res.Error)
					}
				} else {
					log.Info("Successfully fulfilled order")
					dist_store.Set(o.Id, res.Status)
					success = true
				}
				break
			case <-time.After(MAX_TIMEOUT):
				if slices.Contains(retryableFulfillmentStatuses, constants.E_TIMEOUT) {
					log.Error("Fulfillment timed out, but I guess it's retryable, so....", errors.New("timeout"))

					// leave the original channel open to avoid a panic if it just happens to be taking a while.
					fulfillment = c.FulfillOrder(o)
				} else {
					log.Fatal("Fulfillment timed out, unretryable. Crashing.", errors.New("timeout"))
				}
				break
			}
		}

		// here because we exhausted retry budget? gotta die
		if !success {
			log.Fatal("Exhausted retries trying to fulfill order. Crashing. Last error:", err)
		}
	}

	// Update status in case we got here from a retry
	s, ok = dist_store.Get(o.Id)
	log.Debug("Got status from dist_store:", s)

	// allowable state in order to archive
	if s == constants.ORDER_FULFILLED && s != constants.ORDER_ARCHIVED {
		log.Info("Attempting to archive order")
		db, err := archival.NewClient()
		if err != nil {
			log.Fatal("Failed to create archival client", err)
		}
		defer db.Close()

		success := false
		retryDelay := INITIAL_TIMEOUT
		for try := 0; try < MAX_ATTEMPTS && !success; try++ {
			err = db.Persist(o, s)
			if err == nil {
				log.Info("Successfully archived order")
				dist_store.Set(o.Id, constants.ORDER_ARCHIVED)
				success = true
			} else if !err.(archival.RetryableError) {
				log.Fatal("Failed trying to archive order. Unretryable. Crashing. Last error:", err)
			} else {
				log.Error("Archiving order failed. Retrying after", retryDelay, "s. Error:", err)
				retryDelay = expoBackoff(retryDelay)
			}
		}

		// here because we exhausted retry budget? gotta die
		if !success {
			log.Fatal("Exhausted retries trying to archive order. Crashing. Last error:", err)
		}
	}

	// final error check. If successful, status should be archived.
	s, ok = dist_store.Get(o.Id)
	if s != constants.ORDER_ARCHIVED {
		return errors.New(fmt.Sprintf("Order reached unknown state. Expected %s, but got %s",
			constants.ORDER_ARCHIVED,
			s))
	}

	// finished; yay!
	return nil
}

// waits for the given duration and then returns what the "next" delay should be
func expoBackoff(d time.Duration) time.Duration {
	time.Sleep(d)
	n := d * BACKOFF_COEFFICIENT
	if n > MAX_TIMEOUT {
		n = MAX_TIMEOUT
	}
	return n
}

func main() {
	o := app.Order{
		Id: "123-abc",
	}

	err := OrderItem(o)
	if err != nil {
		fmt.Println("Error trying to order something", err)
	}
}
