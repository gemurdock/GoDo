package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"

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

	/*
		Print out entire request (Middleware)
	*/
	router.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestDump, err := httputil.DumpRequest(r, true)
			if err == nil {
				fmt.Print("\n\n")
				fmt.Printf(string(requestDump))
			}
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	})

	/*
		Serve ping.html
	*/
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		responseCount.Count += 1

		compiledTemplates.ExecuteTemplate(w, "ping.html", responseCount)
	})

	/*
		Start server
	*/
	http.ListenAndServe(":80", router)
}
