package server

import (
	//"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/venikx/go-rss/cmd/web"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/js/*", echo.WrapHandler(fileServer))

	e.GET("/", echo.WrapHandler(templ.Handler(web.Base())))
	e.GET("/users", s.HandleUsers)
	e.POST("/users/new", s.HandleNewUser)

	e.GET("/feeds", s.HandleFeeds)
	e.POST("/feeds/new", s.HandleNewFeed)

	e.GET("/healthz", s.healthHandler)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// Take required information from error and context and send it to a service like New Relic
		//log.Println(c.Path(), c.QueryParams(), err.Error())

		// Call the default handler to return the HTTP response
		e.DefaultHTTPErrorHandler(err, c)
	}

	return e
}
