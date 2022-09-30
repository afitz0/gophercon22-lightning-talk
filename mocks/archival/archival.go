package archival

import (
	"math/rand"

	"lightning/app/common"
)

type ArchivalStorage struct{}
type RetryableError bool

func NewClient() (ArchivalStorage, error) {
	return ArchivalStorage{}, nil
}

func (a *ArchivalStorage) Close() {}

func (a *ArchivalStorage) Persist(o common.Order, s string) error {
	common.Sleep(rand.Intn(10), "ArchiveOrder")
	return nil
}

func (r RetryableError) Error() string {
	return "An error occured during the API call that is retryable"
}
