package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"k-microserv-kuber.com/users/controller"
	"k-microserv-kuber.com/users/database/mysql"

	httphandler "k-microserv-kuber.com/users/handler/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting user service on port %d", port)

	db, err := mysql.New()
	if err != nil {
		panic(err)
	}
	svc := controller.New(db)
	h := httphandler.New(svc)
	http.Handle("/user", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
