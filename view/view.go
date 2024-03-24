package view

import (
	"github.com/xpmatteo/gomovies/model"
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, r *http.Request, templ *template.Template, movies []*model.Movie) {
	data := map[string]any{
		"movies": movies,
		"path":   "/",
	}
	err := templ.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
