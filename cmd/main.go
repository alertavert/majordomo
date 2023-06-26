package main

import (
	"log"

	"github.com/alertavert/gpt4-go/pkg/server"
)

func main() {
	if err := server.Setup(); err != nil {
		log.Fatalf("Error setting up server: %v", err)
		return
	}
	log.Fatal(server.Run(":5000"))
}
