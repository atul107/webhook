package main

import (
	"fmt"
	"net/http"
)

//callback method
func homePage(w http.ResponseWriter, r *http.Request) {
	//adding values to w response stream
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
