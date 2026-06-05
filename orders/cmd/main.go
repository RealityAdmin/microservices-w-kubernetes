package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"k-microserv-kuber.com/orders/controller"
	"k-microserv-kuber.com/orders/database/memory"
	serviceGateway "k-microserv-kuber.com/orders/gateway/http"
	httphandler "k-microserv-kuber.com/orders/handler/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting products service on port %d", port)

	db := memory.New()
	ug := serviceGateway.NewUserGateway("http://localhost:8081")
	pg := serviceGateway.NewProductGateway("http://localhost:8082")

	ctrl := controller.New(db, ug, pg)
	h := httphandler.New(ctrl)

	mux := http.NewServeMux()

	mux.HandleFunc("/users", h.HandleUserEndpoint)
	mux.HandleFunc("/products", h.HandleProductEndpoint)

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)

	// go func() {
	// 	http.Handle("/users", http.HandlerFunc(h.HandleUserEndpoint))
	// 	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// go func() {
	// 	http.Handle("/products", http.HandlerFunc(h.HandleProductEndpoint))
	// 	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// err := http.ListenAndServe(":9999", nil)
	// if err != nil {
	// 	panic(err)
	// }
}
