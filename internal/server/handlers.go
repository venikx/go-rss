package server

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/venikx/go-rss/cmd/web"
)

func (s *Server) HandleUsers(c echo.Context) error {
	users, err := s.db.ReadUsers()
	if err != nil {
		return err
	}
	log.Printf("%v", users)

	component := web.UsersPage(users)
	return component.Render(context.Background(), c.Response().Writer)
}

func (s *Server) HandleNewUser(c echo.Context) error {
	name := c.FormValue("name")
	user, _ := s.db.CreateUser(name)
	component := web.User(user)
	return component.Render(context.Background(), c.Response().Writer)
}

func (s *Server) HandleFeeds(c echo.Context) error {
	feeds, err := s.db.ReadFeeds()
	if err != nil {
		return err
	}
	log.Printf("%v", feeds)

	component := web.FeedsPage(feeds)
	return component.Render(context.Background(), c.Response().Writer)
}

func (s *Server) HandleNewFeed(c echo.Context) error {
	name := c.FormValue("name")
	url := c.FormValue("url")
	userIdStr := c.FormValue("userId")

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return err
	}

	feed, err := s.db.CreateFeed(name, url, userId)
	if err != nil {
		return err
	}

	component := web.Feed(feed)
	return component.Render(context.Background(), c.Response().Writer)
}

func (s *Server) healthHandler(c echo.Context) error {
	msg, _ := s.db.Health()
	return c.JSON(http.StatusOK, msg)
}
