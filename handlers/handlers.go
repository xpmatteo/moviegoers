package handlers

import (
	"github.com/xpmatteo/gomovies/view"
	"html/template"
	"net/http"
)

func Index(templ *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		view.Render(w, r, templ)
	})
}
