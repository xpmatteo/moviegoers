package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/gomovies/model"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockRepository struct {
	willReturnMovies []model.Movie
}

var passedOptions model.QueryOptions

func (repo MockRepository) Query(opts model.QueryOptions) []model.Movie {
	passedOptions = opts
	return repo.willReturnMovies
}

type ZC struct{}

func (z ZC) Today() time.Time {
	return time.Time{}
}

var zeroCalendar = ZC{}

func Test_index(t *testing.T) {
	tests := []struct {
		name                 string
		template             string
		url                  string
		expectedBody         string
		expectedQueryOptions model.QueryOptions
	}{
		{
			name:                 "movies",
			template:             "{{ .movies }}",
			url:                  "/",
			expectedBody:         "[]",
			expectedQueryOptions: model.QueryOptions{Page: 1},
		},
		{
			name:                 "nextPage",
			template:             "{{.nextPage}}",
			url:                  "/?page=7",
			expectedBody:         "8",
			expectedQueryOptions: model.QueryOptions{Page: 7},
		},
		{
			name:                 "nextPage default",
			template:             "{{.nextPage}}",
			url:                  "/",
			expectedBody:         "2",
			expectedQueryOptions: model.QueryOptions{Page: 1},
		},
		{
			name:                 "genre",
			template:             "{{.selectedGenre}}",
			url:                  "/?genre=123",
			expectedBody:         "123",
			expectedQueryOptions: model.QueryOptions{Page: 1, Genre: 123},
		},
		{
			name:                 "genre default",
			template:             "{{.selectedGenre}}",
			url:                  "/",
			expectedBody:         "0",
			expectedQueryOptions: model.QueryOptions{Page: 1, Genre: 0},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			templ := template.Must(template.New("index").Parse(test.template))
			repo := MockRepository{}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, test.url, nil)

			Index(templ, repo, zeroCalendar).ServeHTTP(w, r)

			assert.Equal(t, test.expectedBody, w.Body.String())
			assert.Equal(t, test.expectedQueryOptions, passedOptions)
		})
	}
}
