package main

import "github.com/gorilla/mux"

import (
	"fmt"
	"log"
	"net/http"
)

func startRouter(host string) {

	r := mux.NewRouter()
	r.HandleFunc("/humpbacknode/v1/_ping", ping).Methods("GET")
	http.Handle("/", r)
	fmt.Println("humpback-node api starting...")
	if err := http.ListenAndServe(host, nil); err != nil {
		log.Printf("humpback-node api error: %s\n", err.Error())
	}
}

func ping(w http.ResponseWriter, req *http.Request) {

	w.Write([]byte{'P', 'A', 'N', 'G'})
	w.WriteHeader(200)
}
