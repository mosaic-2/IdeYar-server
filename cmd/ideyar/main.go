package main

import (
	"log"
)

func main() {
	log.Println("Starting server...")

	if err := serve(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
