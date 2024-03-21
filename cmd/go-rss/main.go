package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/venikx/go-rss/internal/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: routes.Multiplexer(),
	}

	fmt.Printf("Running server on http://localhost:%s", port)

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
