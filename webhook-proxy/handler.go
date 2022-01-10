package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

//callback method
func homePage(w http.ResponseWriter, r *http.Request) {
	//adding values to w response stream
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

//callback for /proxy endpoint
func proxy(w http.ResponseWriter, r *http.Request) {
	defer fmt.Println("Endpoint Hit: proxy ", r.Method)

	if r.URL.Path != "/proxy" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {

	case "GET":
		fmt.Fprintf(w, "Welcome to the Proxy!")

	case "POST":
		// Parse Body
		decoder := json.NewDecoder(r.Body)
		var t postBody
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		//url check
		_, err = url.ParseRequestURI(t.Url)
		if err != nil {
			panic(err)
		}

		// Webhook Call
		req, err := http.NewRequest("POST", t.Url, bytes.NewReader(t.Payload))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Not able to make request!"))
			log.Println(err)
		}
		req.Header.Set("host", t.Url)
		req.Header.Set("Content-Type", "application/json")
		for key, element := range t.Headers {
			req.Header.Set(key, element)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Bad Response"))
			log.Println(err)
		}

		//Performs one Retry for retriable errors
		for checkRetry(resp.StatusCode) {
			resp, err = client.Do(req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Bad Response"))
				log.Println(err)
				break
			}
		}

		// Success and Failure Messages
		if resp.StatusCode == 200 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 - Success Response"))
			log.Println("Success Response")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Bad Response"))
			log.Println("Bad Response")
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
