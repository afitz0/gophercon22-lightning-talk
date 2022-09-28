package archival

import (
	"lightning/app"
)

type ArchivalStorage struct{}
type RetryableError bool

func NewClient() (ArchivalStorage, error) {
	return ArchivalStorage{}, nil
}

func (a *ArchivalStorage) Close() {}

func (a *ArchivalStorage) Persist(o app.Order, s string) error {
	return nil
}

func (r RetryableError) Error() string {
	return "An error occured during the API call that is retryable"
}
