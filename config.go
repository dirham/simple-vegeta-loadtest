package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	EnableTokenGen bool `json:"enable_token_gen"`
	TokenGenAttr   struct {
		URL     string `json:"url"`
		Headers struct {
			Authorization struct {
				Type     string `json:"type"`
				Username string `json:"username"`
				Password string `json:"password"`
			} `json:"Authorization"`
		} `json:"headers"`
		Body struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"body"`
	} `json:"token_gen_attr"`
	Tergets []struct {
		URL string `json:"url"`
	} `json:"tergets"`
}

func LoadConfiguration(file_path string) Config {
	var config Config
	configFile, err := os.Open(file_path)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)

	jsonParser.Decode(&config)

	return config
}
