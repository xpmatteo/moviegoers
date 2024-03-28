package main

import (
	"github.com/joho/godotenv"
	"github.com/xpmatteo/gomovies/adapters"
	"github.com/xpmatteo/gomovies/handlers"
	"html/template"
	"log"
	"net/http"
	"os"
)

var templ = template.Must(template.ParseFiles("view/index.tmpl"))

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %s", err)
	}

	http.Handle("GET /{$}", handlers.Index(templ, &adapters.Mtdb{Agent: http.DefaultClient}, handlers.DefaultCalendar))
	http.Handle("GET /", http.FileServer(http.Dir("./public/")))

	var port = os.Getenv("PORT")
	log.Print("Serving HTTP from port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
