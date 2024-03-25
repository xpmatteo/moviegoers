package view

import (
	"github.com/xpmatteo/gomovies/model"
	"html/template"
	"net/http"
	"strconv"
)

func Render(w http.ResponseWriter, r *http.Request, templ *template.Template, movies []model.Movie) {
	data := Model(movies, r)
	err := templ.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func Model(movies []model.Movie, r *http.Request) map[string]any {
	thisPage, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		thisPage = 1
	}
	data := map[string]any{
		"movies":   movies,
		"nextPage": thisPage + 1,
		"genres":   model.Genres,
	}
	return data
}
