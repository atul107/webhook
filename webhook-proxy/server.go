package main

import (
	"fmt"
	"log"
	"net/http"
)

//server endpoints and callbacks
func startServer() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/proxy", proxy)
	serverHost := fmt.Sprintf("%s:%s", config.BindIp, config.BindPort)
	// fmt.Print(serverHost)
	log.Fatal(http.ListenAndServe(serverHost, nil))
}
