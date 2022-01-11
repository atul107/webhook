package main

import (
	"log"
	"net/http"
)

//server endpoints and callbacks
func startServer(host, port string) {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/proxy", proxy)
	serverPort := host + port
	log.Fatal(http.ListenAndServe(serverPort, nil))
}
