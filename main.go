package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/xpmatteo/gomovies/handlers"
	"github.com/xpmatteo/gomovies/model"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

const port = "8080"

var data struct {
	Results []model.Movie `json:"Results"`
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	readMovies()
	fmt.Printf("Read %d movies\n", len(data.Results))
	
	templ := template.Must(template.ParseFiles("view/index.tmpl"))
	http.Handle("GET /{$}", handlers.Index(templ, data.Results))
	http.Handle("GET /", http.FileServer(http.Dir("./public/")))
	log.Print("Serving HTTP from port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func readMovies() {
	file, err := os.Open("movies.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Fatal("Error:", err)
	}
}
