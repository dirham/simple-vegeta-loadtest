package main

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

func PostHttp(url string, body *bytes.Reader, timeout int64, headers []map[string]string) *http.Response {
	client := &http.Client{
		Timeout: time.Duration(time.Second * time.Duration(timeout)),
	}
	// peform http request to get token
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatal("Got Error create new request", err.Error())
	}

	for _, v := range headers {
		for k, val := range v {
			req.Header.Set(k, val)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatal("Got Error", err.Error())
	}

	return response
}
