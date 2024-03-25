package main

import (
	"github.com/joho/godotenv"
	"github.com/xpmatteo/gomovies/adapters"
	"github.com/xpmatteo/gomovies/handlers"
	"html/template"
	"log"
	"net/http"
)

const port = "8080"

var templ = template.Must(template.ParseFiles("view/index.tmpl"))

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	http.Handle("GET /{$}", handlers.Index(templ, &adapters.Mtdb{}))
	http.Handle("GET /", http.FileServer(http.Dir("./public/")))
	log.Print("Serving HTTP from port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
