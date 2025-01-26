package main

import (
	"fmt"
	"time"

	"github.com/venikx/go-rss/database"
	"github.com/venikx/go-rss/scripts"
	"github.com/venikx/go-rss/server"
)

func main() {
	_ = database.New()
	server := server.NewServer()

	go scripts.StartScraping(10, time.Minute)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

}
