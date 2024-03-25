package handlers

import (
	"github.com/xpmatteo/gomovies/model"
	"github.com/xpmatteo/gomovies/view"
	"html/template"
	"net/http"
)

type MovieRepository interface {
	Query(options model.QueryOptions) []model.Movie
}

func Index(templ *template.Template, repo MovieRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		view.Render(w, r, templ, repo.Query(model.QueryOptions{}))
	})
}
