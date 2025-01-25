package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello-world", s.HelloWorldHandler)

	// api
	mux.HandleFunc("GET /api/health", s.healthHandler)

	return mux
	//e := echo.New()
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	//fileServer := http.FileServer(http.FS(web.Files))
	//e.GET("/js/*", echo.WrapHandler(fileServer))

	//e.GET("/", echo.WrapHandler(templ.Handler(web.Base())))
	//e.GET("/users", s.HandleUsers)
	//e.POST("/users/new", s.HandleNewUser)

	//e.GET("/feeds", s.HandleFeeds)
	//e.POST("/feeds/new", s.HandleNewFeed)

	//e.HTTPErrorHandler = func(err error, c echo.Context) {
	//	// Take required information from error and context and send it to a service like New Relic
	//	//log.Println(c.Path(), c.QueryParams(), err.Error())

	//	// Call the default handler to return the HTTP response
	//	e.DefaultHTTPErrorHandler(err, c)
	//}

	//return e
}
