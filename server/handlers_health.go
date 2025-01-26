package server

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/venikx/go-rss/database"
)

type Helmet struct {
	Title  string
	Author string
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	t, err := template.ParseFiles("views/base.html", "views/hello.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	hello, err := database.HelloWorld(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = t.Execute(w, Helmet{Title: hello, Author: "from PostgreSQL!"})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	jsonResp, err := json.Marshal(database.Health(ctx))

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
