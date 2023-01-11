package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func handleToken(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}

// func (s *models.Server) Greeting() http.HandlerFunc {

// 	s.Router.HandleFunc("hello", handleToken)
// }

func Hello(router *mux.Router) {

	router.HandleFunc("/hello", handleToken)
}
