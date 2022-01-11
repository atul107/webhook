package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

//structure for response body
type postBody struct {
	Url     string            `json:"url"`
	Payload json.RawMessage   `json:"payload"`
	Headers map[string]string `json:"headers"`
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

// Parse Request Json Body
func parseJsonBody(body io.ReadCloser) (postBody, error) {
	decoder := json.NewDecoder(body)
	var t postBody
	err := decoder.Decode(&t)
	return t, err
}

// Create Request for Webhook
func createRequest(t postBody) (*http.Request, error) {
	req, err := http.NewRequest("POST", t.Url, bytes.NewReader(t.Payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("host", t.Url)
	req.Header.Set("Content-Type", "application/json")
	for key, element := range t.Headers {
		req.Header.Set(key, element)
	}
	return req, nil
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
		t, err := parseJsonBody(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad body format!"))
			log.Println("Bad body format!")
			return
		}

		// Url check
		_, err = url.ParseRequestURI(t.Url)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Incorrect Url!"))
			log.Println("Incorrect Url!")
			return
		}

		// Webhook Request
		req, err := createRequest(t)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Not able to make request!"))
			log.Println("Not able to make request!")
			return
		}

		// Calling Webhook
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Bad Response"))
			log.Println("Bad Response")
			return
		}
		defer resp.Body.Close()

		//Performs one Retry for retriable errors
		for checkRetry(resp.StatusCode) {
			resp, err = client.Do(req)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Bad Response"))
				log.Println("Bad Response")
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
