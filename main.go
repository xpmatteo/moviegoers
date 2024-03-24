package main

import (
	"encoding/json"
	"fmt"
	"github.com/xpmatteo/gomovies/handlers"
	"github.com/xpmatteo/gomovies/model"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

const port = "8080"

func main() {
	file, err := os.Open("movies.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var data struct {
		results []model.Movie
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Read %d movies\n", len(data.results))
	templ := template.Must(template.ParseFiles("view/index.gotmpl"))
	http.Handle("GET /{$}", handlers.Index(templ, data.results))
	http.Handle("GET /", http.FileServer(http.Dir("./public/")))
	log.Print("Serving HTTP from port ", port)
	log.Panic(http.ListenAndServe(":"+port, nil))
}
