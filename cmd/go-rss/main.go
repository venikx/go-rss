package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Hello Nix!")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	fmt.Printf("Running server on http://localhost:%s", port)
}
