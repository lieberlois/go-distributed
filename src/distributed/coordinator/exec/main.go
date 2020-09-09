package main

import (
	"fmt"
	"go-distributed/src/distributed/coordinator"
)

func main() {
	ql := coordinator.NewQueueListener()
	go ql.ListenForNewSource()

	var s string
	_, _ = fmt.Scanln(&s)
}


