package main

import (
	"encoding/json"
	"fmt"
	"os"
)

//Structure for config
type Config struct {
	BindIp   string `json:"bind_ip"`
	BindPort string `json:"bind_port"`
}

func readConfig(f *os.File) Config {
	decoder := json.NewDecoder(f)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}
