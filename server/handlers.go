package server

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/venikx/go-rss/database"
	"github.com/venikx/go-rss/typings"
)

type UsersPage struct {
	Title string
	Users []typings.User
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	users, err := database.ReadUsers(ctx)
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

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	name := r.FormValue("name")

	_, err := database.CreateUser(ctx, name)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.WriteHeader(http.StatusCreated)
}

type FeedsPage struct {
	Title string
	Feeds []typings.Feed
}

func handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	feeds, err := database.ReadFeeds(ctx)
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

func handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	name := r.FormValue("name")
	url := r.FormValue("url")
	hasFollow := r.FormValue("follow")

	feed, err := database.CreateFeed(ctx, name, url, userId)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if hasFollow != "" {
		_, err = database.FollowFeed(ctx, feed.Id, userId)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}

	w.WriteHeader(http.StatusCreated)
}
