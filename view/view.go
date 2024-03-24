package view

import (
	"html/template"
	"net/http"
)

func Render(w http.ResponseWriter, r *http.Request, templ *template.Template) {
	data := map[string]any{
		"movies": nil,
		"path":   "path",
	}
	err := templ.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
