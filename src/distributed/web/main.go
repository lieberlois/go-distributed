package main

import (
	"go-distributed/src/distributed/web/controller"
	"log"
	"net/http"
)

func main() {
	controller.Initialize()
	port := ":3000"
	log.Printf("Now listening on port %s", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err.Error())
	}
}
