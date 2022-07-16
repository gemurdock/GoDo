package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type PingResponse struct {
	Count int32
}

func main() {
	fmt.Printf("Starting server\n\n")

	responseCount := PingResponse{
		Count: 0,
	}

	compiledTemplates, err := template.ParseGlob("./templates/*")
	check(err)

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		responseCount.Count += 1

		compiledTemplates.ExecuteTemplate(w, "ping.html", responseCount)
	})

	http.ListenAndServe(":80", router)
}
