package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	/**
	*	Reading Configruation
	 */
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	config := readConfig(file)
	fmt.Println("Server running on", config.BindIp, config.BindPort)

	/**
	*	Persistance storage in file
	 */
	f, err := os.OpenFile("./logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f) //log is recorded into file

	/**
	*	Server
	 */
	startServer(config.BindIp, config.BindPort)
}
