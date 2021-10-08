package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type Token struct {
	AccessToken string `json:"access_token,omitempty"`
	UserId      string `json:"user_id,omitempty"`
	Expires_in  int    `json:"expires_in,omitempty"`
	Message     string `json:"message,omitempty"`
}

func main() {
	// load config.json
	cfg := LoadConfiguration("config.json")
	var token Token

	if cfg.EnableTokenGen {
		// prepare for get accesstoken
		body, err := json.Marshal(cfg.TokenGenAttr.Body)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error Found: %s", err.Error()))
		}
		reqBody := bytes.NewReader(body)
		// peform http request to get token
		head := make(map[string]string)

		sEnc := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.TokenGenAttr.Headers.Authorization.Username, cfg.TokenGenAttr.Headers.Authorization.Password)))
		head["Authorization"] = fmt.Sprintf("%s %s", cfg.TokenGenAttr.Headers.Authorization.Type, sEnc)
		var heads []map[string]string
		heads = append(heads, head)
		tokenRes := PostHttp(cfg.TokenGenAttr.URL, reqBody, 5, heads)

		defer tokenRes.Body.Close()

		err = json.NewDecoder(tokenRes.Body).Decode(&token)
		if err != nil {
			log.Fatal("Got Error decode response ", err.Error())
		}
		log.Printf("%v", token)
	}

	// peform others url with vegeta test
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 4 * time.Second
	for _, v := range cfg.Tergets {

		targeter := vegeta.NewStaticTargeter(vegeta.Target{
			Method: "GET",
			URL:    v.URL,
			Header: http.Header{"Authorization": []string{fmt.Sprintf("%s %s", "Token", v.Token)}},
		})

		for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
			metrics.Add(res)
		}
		metrics.Close()

		fmt.Printf("========== Target %s ============\n", v.URL)
		fmt.Printf("Total latencie: %s \n", metrics.Latencies.Total)
		fmt.Printf("Mean latencie: %s \n", metrics.Latencies.Mean)
		fmt.Printf("maximum observed request latency: %s \n", metrics.Latencies.Max)
		fmt.Printf("Duration of the request attack: %s \n", metrics.Duration)
		fmt.Printf("Throughput is the rate of successful requests per second: %f \n", metrics.Throughput)
		fmt.Printf("Success is the percentage of non-error responses: %f \n", metrics.Success)
		fmt.Printf("rate of sent requests per second : %f\n", metrics.Rate)

	}

}
