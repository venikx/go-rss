package server

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/venikx/go-rss/typings"
)

type UsersPage struct {
	Title string
	Users []typings.User
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	users, err := s.db.ReadUsers(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	t, err := template.ParseFiles("views/base.html", "views/partials/users.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = t.Execute(w, UsersPage{Title: "All Users", Users: users})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	name := r.FormValue("name")

	_, err := s.db.CreateUser(ctx, name)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.WriteHeader(http.StatusCreated)
}

type FeedsPage struct {
	Title string
	Feeds []typings.Feed
}

func (s *Server) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	feeds, err := s.db.ReadFeeds(ctx)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	t, err := template.ParseFiles("views/base.html", "views/partials/feeds.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = t.Execute(w, FeedsPage{Title: "All Feeds", Feeds: feeds})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// TODO(Kevin): Get from authentication somehow
var userId = 1

func (s *Server) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	name := r.FormValue("name")
	url := r.FormValue("url")
	hasFollow := r.FormValue("follow")

	feed, err := s.db.CreateFeed(ctx, name, url, userId)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if hasFollow != "" {
		_, err = s.db.FollowFeed(ctx, feed.Id, userId)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}

	w.WriteHeader(http.StatusCreated)
}
