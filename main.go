package main

import (
	"github.com/joho/godotenv"
	"github.com/xpmatteo/gomovies/handlers"
	"github.com/xpmatteo/gomovies/mtdb"
	"html/template"
	"log"
	"net/http"
	"os"
)

var templ = template.Must(template.ParseFiles("web/index.tmpl"))

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %s", err)
	}

	mtdb := mtdb.NewMemoizer(mtdb.NewMtdb())
	http.Handle("GET /{$}", handlers.Index(templ, mtdb, handlers.DefaultCalendar))
	http.Handle("GET /", http.FileServer(http.Dir("./public/")))

	var port = os.Getenv("PORT")
	log.Print("Serving HTTP from port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
