package main

import (
	"log"

	"github.com/hum/mc-rcon/rcon"
)

var commands = [...]string{
	"help",
	"help 2",
	"time set day",
}

func main() {
	client, err := rcon.CreateClient("ip_address", "port")
	if err != nil {
		log.Fatalf("Could not connect to the server: %v\n", err)
	}
	defer client.Close()

	err = client.Authenticate("password")
	if err != nil {
		log.Fatalf("Failed to authenticate: %v\n", err)
	}

	for _, command := range commands {
		response, err := client.SendCommand(command)
		if err != nil {
			log.Fatalf("Could not process the command: %v\n", err)
			continue
		}
		log.Println(response)
	}
}
