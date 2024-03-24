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
	readMoviesFromApi()
	fmt.Printf("Read %d movies\n", len(data.Results))

	templ := template.Must(template.ParseFiles("view/index.tmpl"))
	http.Handle("GET /{$}", handlers.Index(templ, data.Results))
	http.Handle("GET /", http.FileServer(http.Dir("./public/")))
	log.Print("Serving HTTP from port ", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func readMoviesFromApi() {
	url := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc"
	url = "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&primary_release_date.lte=2024-03-24&sort_by=primary_release_date.desc"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+os.Getenv("TMDB_ACCESS_TOKEN"))
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Fatal("Error:", err)
	}
}
