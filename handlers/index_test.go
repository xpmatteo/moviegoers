package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/gomovies/domain"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockRepository struct {
	willReturnMovies []domain.Movie
}

var passedOptions domain.QueryOptions

func (repo MockRepository) Query(opts domain.QueryOptions) []domain.Movie {
	passedOptions = opts
	return repo.willReturnMovies
}

func Test_index(t *testing.T) {
	tests := []struct {
		name                 string
		template             string
		url                  string
		maxReleaseDate       time.Time
		expectedBody         string
		expectedQueryOptions domain.QueryOptions
	}{
		{
			name:                 "movies",
			template:             "{{ .movies }}",
			url:                  "/",
			expectedBody:         "[]",
			expectedQueryOptions: domain.QueryOptions{Page: 1},
		},
		{
			name:                 "nextPage",
			template:             "{{.nextPage}}",
			url:                  "/?page=7",
			expectedBody:         "8",
			expectedQueryOptions: domain.QueryOptions{Page: 7},
		},
		{
			name:                 "nextPage default",
			template:             "{{.nextPage}}",
			url:                  "/",
			expectedBody:         "2",
			expectedQueryOptions: domain.QueryOptions{Page: 1},
		},
		{
			name:                 "genre",
			template:             "{{.selectedGenre}}",
			url:                  "/?genre=123",
			expectedBody:         "123",
			expectedQueryOptions: domain.QueryOptions{Page: 1, Genre: 123},
		},
		{
			name:                 "genre default",
			template:             "{{.selectedGenre}}",
			url:                  "/",
			expectedBody:         "0",
			expectedQueryOptions: domain.QueryOptions{Page: 1, Genre: 0},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			templ := template.Must(template.New("index").Parse(test.template))
			repo := MockRepository{}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, test.url, nil)

			cal := func() time.Time {
				return test.maxReleaseDate
			}
			Index(templ, repo, cal).ServeHTTP(w, r)

			assert.Equal(t, test.expectedBody, w.Body.String())
			assert.Equal(t, test.expectedQueryOptions, passedOptions)
		})
	}
}
