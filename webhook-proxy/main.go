package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

//structure for response body
type postBody struct {
	Url     string            `json:"url"`
	Payload json.RawMessage   `json:"payload"`
	Headers map[string]string `json:"headers"`
}

//Structure for config
type Configuration struct {
	Users  []string
	Groups []string
}

func main() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(configuration.Users)

	// persistance storage in file
	f, err := os.OpenFile("./logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	//log is recorded into file
	log.SetOutput(f)
	//starts server
	handleRequests()
}

//server endpoints and callbacks
func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/proxy", proxy)
	log.Fatal(http.ListenAndServe(":3004", nil))
}

//checks for retriable errors
func checkRetry(status int) bool {
	validStatus := [3]int{502, 503, 504}
	for _, v := range validStatus {
		if status == v {
			log.Println("Retry Tried")
			return true
		}
	}
	return false
}
