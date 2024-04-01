package handlers

import (
	"github.com/xpmatteo/gomovies/domain"
	"github.com/xpmatteo/gomovies/web"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type MovieRepository interface {
	Query(options domain.QueryOptions) []domain.Movie
}

type CalendarFunc func() time.Time

func DefaultCalendar() time.Time {
	return time.Now()
}

func Index(templ *template.Template, repo MovieRepository, today CalendarFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = 1
		}
		genre, err := strconv.Atoi(r.URL.Query().Get("genre"))
		if err != nil {
			genre = 0
		}
		options := domain.QueryOptions{
			Page:           page,
			Genre:          genre,
			ReleaseDateMax: today(),
		}
		data := web.Model(repo.Query(options), options)
		if err := templ.Execute(w, data); err != nil {
			panic(err)
		}
	})
}
