package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

//structure for response body
type postBody struct {
	Url     string            `json:"url"`
	Payload json.RawMessage   `json:"payload"`
	Headers map[string]string `json:"headers"`
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
		var resp *http.Response
		err = retry(config.RetryAttemp, (config.RetryInterval)*time.Second, func() (err error) {
			resp, err = client.Do(req)
			return
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(string(err.Error())))
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		// Success and Failure Messages
		if resp.StatusCode == 200 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("200 - Success Response"))
			log.Println("Success Response from webhook")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Bad Response from webhook"))
			log.Println("Bad Response from webhook")
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
