package handlers

import (
	"github.com/xpmatteo/gomovies/model"
	"github.com/xpmatteo/gomovies/view"
	"html/template"
	"net/http"
	"strconv"
)

type MovieRepository interface {
	Query(options model.QueryOptions) []model.Movie
}

func Index(templ *template.Template, repo MovieRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = 1
		}
		genre, err := strconv.Atoi(r.URL.Query().Get("genre"))
		if err != nil {
			genre = 0
		}
		options := model.QueryOptions{
			Page:  page,
			Genre: genre,
		}
		data := view.Model(repo.Query(options), options)
		if err := templ.Execute(w, data); err != nil {
			panic(err)
		}
	})
}
