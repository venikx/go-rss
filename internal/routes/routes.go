package routes

import (
	"fmt"
	"html"
	"net/http"
)

func Multiplexer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q using Nix!", html.EscapeString(r.URL.Path))
}
