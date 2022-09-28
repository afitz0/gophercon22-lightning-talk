package archival

import (
	"lightning/app"
)

type ArchivalStorage struct{}

func NewClient() (ArchivalStorage, error) {
	return ArchivalStorage{}, nil
}

func (a *ArchivalStorage) Close() {}

func (a *ArchivalStorage) Persist(o app.Order, s app.OrderStatus) error {
	return nil
}
