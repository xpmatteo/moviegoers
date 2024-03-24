package handlers

import (
	"github.com/xpmatteo/gomovies/model"
	"github.com/xpmatteo/gomovies/view"
	"html/template"
	"net/http"
)

func Index(templ *template.Template) http.Handler {
	movies := []*model.Movie{
		{
			Title:    "ABC",
			Overview: "AHA OHO",
		},
		{
			Title:    "ABC",
			Overview: "AHA OHO",
		},
		{
			Title:    "ABC",
			Overview: "AHA OHO",
		},
		{
			Title:    "ABC",
			Overview: "AHA OHO",
		},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		view.Render(w, r, templ, movies)
	})
}
