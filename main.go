package main

import (
	"encoding/json"
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

type RequestedValue struct {
	Amount int64 `json:"amount"`
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
		Count each request (Middleware)
	*/
	router.Use(func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			responseCount.Count += 1
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	})

	/*
		Serve ping.html
	*/
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		compiledTemplates.ExecuteTemplate(w, "ping.html", responseCount)
	})

	/*
		Serve post.html
	*/
	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		amount := RequestedValue{
			Amount: -1,
		}
		err := decoder.Decode(&amount)
		check(err)

		if amount.Amount < 0 {
			http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
			return
		}

		compiledTemplates.ExecuteTemplate(w, "post.html", amount)
	})

	/*
		Start server
	*/
	http.ListenAndServe(":80", router)
}
