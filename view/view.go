package view

import (
	"github.com/xpmatteo/gomovies/model"
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, r *http.Request, templ *template.Template, movies []model.Movie) {
	data := Model(movies, r)
	err := templ.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func Model(movies []model.Movie, r *http.Request) map[string]any {
	data := map[string]any{
		"movies": movies,
	}
	return data
}
