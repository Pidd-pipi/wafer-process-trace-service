package main

import (
	"example.com/wafer-process-trace-service/config"
	"example.com/wafer-process-trace-service/httpapi"
	"example.com/wafer-process-trace-service/store"
	"example.com/wafer-process-trace-service/web"
	"log"
)

func main() {
	address := ":" + config.Port()
	log.Printf("wafer-process-trace-service listening on %s", address)
	if err := serveAddress(address, httpapi.NewHandler(store.New(), web.FS)); err != nil {
		log.Fatal(err)
	}
}
