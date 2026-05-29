package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"k-microserv-kuber.com/products/controller"
	"k-microserv-kuber.com/products/database/memory"

	httphandler "k-microserv-kuber.com/products/handler/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()
	log.Printf("Starting products service on port %d", port)

	db := memory.New()
	svc := controller.New(db)
	h := httphandler.New(svc)
	http.Handle("/products", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
