package main

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/gomovies/adapters"
	"github.com/xpmatteo/gomovies/handlers"
	"github.com/xpmatteo/gomovies/model"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const baseQuery = "include_adult=false&include_video=false&language=en-US&sort_by=popularity.desc"

var testCases = []struct {
	name              string
	requestPath       string
	expectedTmdbQuery string
	returnedMovies    []model.Movie
	assertions        func(*testing.T, *goquery.Document)
}{
	{
		name:        "default",
		requestPath: "/",
		expectedTmdbQuery: baseQuery +
			"&page=1" +
			"&primary_release_date.lte=2024-03-02",
		returnedMovies: make([]model.Movie, 20),
		assertions: func(t *testing.T, document *goquery.Document) {
			assert.Equal(t, "/?page=2", attribute(document, "a#moreMovies", "href"))
		},
	},
	{
		name:        "all genres page 2",
		requestPath: "/?page=2",
		expectedTmdbQuery: baseQuery +
			"&page=2" +
			"&primary_release_date.lte=2024-03-02",
		returnedMovies: make([]model.Movie, 20),
		assertions: func(t *testing.T, document *goquery.Document) {
			assert.Equal(t, "/?page=3", attribute(document, "a#moreMovies", "href"))
		},
	},
}

func attribute(document *goquery.Document, selector, attribute string) string {
	val, exists := document.Find(selector).Attr(attribute)
	if !exists {
		val = ""
	}
	return val
}

type mockMtdb struct {
	expectedQuery string
	toBeReturned  []model.Movie
	t             *testing.T
}

func (m mockMtdb) Do(req *http.Request) (*http.Response, error) {
	assert.Equal(m.t, "/3/discover/movie", req.URL.Path)
	assert.Equal(m.t, m.expectedQuery, req.URL.RawQuery)
	var data = struct {
		Results []model.Movie `json:"Results"`
	}{m.toBeReturned}
	body, err := json.Marshal(data)
	if err != nil {
		m.t.Error(err)
	}
	return &http.Response{
		Status:        "OK",
		StatusCode:    200,
		Body:          io.NopCloser(bytes.NewReader(body)),
		ContentLength: 0,
	}, nil
}

func TestEndToEnd(t *testing.T) {
	var templ = template.Must(template.ParseFiles("view/index.tmpl"))
	if err := os.Setenv("TMDB_ACCESS_TOKEN", "anything"); err != nil {
		panic(err)
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, test.requestPath, nil)
			mockAgent := mockMtdb{
				expectedQuery: test.expectedTmdbQuery,
				toBeReturned:  test.returnedMovies,
				t:             t,
			}
			handlers.Index(templ, &adapters.Mtdb{Agent: &mockAgent}, handlers.FakeCalendar{2024, 3, 2}).ServeHTTP(w, r)
			document, err := goquery.NewDocumentFromReader(w.Body)
			if err != nil {
				t.Error(err)
			}
			test.assertions(t, document)
		})
	}
}
