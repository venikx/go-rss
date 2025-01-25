package server

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Helmet struct {
	Title  string
	Author string
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/base.html", "views/hello.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	hello, err := s.db.HelloWorld()
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

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
