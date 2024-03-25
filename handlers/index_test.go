package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/gomovies/model"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockRepository struct {
	willReturnMovies []model.Movie
}

func (repo MockRepository) Query(opts model.QueryOptions) []model.Movie {
	return repo.willReturnMovies
}

func Test_indexx(t *testing.T) {
	tests := []struct {
		name     string
		template string
		url      string
		expected string
	}{
		{
			name:     "movies",
			template: "{{ .movies }}",
			url:      "/",
			expected: "[]",
		},
		{
			name:     "nextPage",
			template: "{{.nextPage}}",
			url:      "/?page=7",
			expected: "8",
		},
		{
			name:     "nextPage default",
			template: "{{.nextPage}}",
			url:      "/",
			expected: "2",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			templ := template.Must(template.New("index").Parse(test.template))
			repo := MockRepository{}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, test.url, nil)

			Index(templ, repo).ServeHTTP(w, r)

			assert.Equal(t, test.expected, w.Body.String())
		})
	}
}
