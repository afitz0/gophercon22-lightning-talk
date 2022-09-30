package common

import (
	"fmt"
	"time"
)

type Order struct {
	Id string
}
type OrderRequest struct{}

func Sleep(duration_seconds int, label string) {
	sec := "second"
	if duration_seconds != 1 {
		sec = sec + "s"
	}
	fmt.Printf("%s will take about \033[91m%v\033[0m %s.\n", label, duration_seconds, sec)
	time.Sleep(time.Second * time.Duration(duration_seconds))
}
