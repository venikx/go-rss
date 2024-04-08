package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/venikx/go-rss/cmd/web"
)

func (s *Server) HandleUsers(c echo.Context) error {
	component := web.UsersPage()
	return component.Render(context.Background(), c.Response().Writer)
}

func (s *Server) HandleNewUser(c echo.Context) error {
	name := c.FormValue("name")
	component := web.User(name)
	return component.Render(context.Background(), c.Response().Writer)
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
