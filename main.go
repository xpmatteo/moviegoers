package main

import (
	"github.com/xpmatteo/gomovies/handlers"
	"html/template"
	"log"
	"net/http"
)

func main() {
	templ := template.Must(template.ParseFiles("view/index.gotmpl"))
	http.Handle("GET /{$}", handlers.Index(templ))
	http.Handle("GET /", http.FileServer(http.Dir("./public/")))
	log.Panic(http.ListenAndServe(":8080", nil))
}
