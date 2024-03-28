// This is an almost end-to-end acceptance test.
//
// We treat the application as an hexagon with an input port for the HTTP input request,
// and an output port for the HTTP request we send to the MTDB.
//
// Given a certain HTTP request, we check
//  (a) the query parameters that we send to MTDB
//  (b) the html that we produce in response to the user request
//
// We avoid making real HTTP requests, in the interest of speed

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
	"time"
)

const baseQuery = "include_adult=false&include_video=false&language=en-US&sort_by=popularity.desc"

// FakeCalendar returns a CalendarFunc that always returns the same date
func FakeCalendar(Year, Month, Day int) handlers.CalendarFunc {
	return func() time.Time {
		return time.Date(Year, time.Month(Month), Day, 0, 0, 0, 0, time.UTC)
	}
}

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
			assert.Equal(t, "/?page=2&genre=0", attribute(document, "a#moreMovies", "href"))
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
			assert.Equal(t, "/?page=3&genre=0", attribute(document, "a#moreMovies", "href"))
		},
	},
	{
		name:        "genre specified",
		requestPath: "/?genre=123",
		expectedTmdbQuery: baseQuery +
			"&with_genres=123" +
			"&page=1" +
			"&primary_release_date.lte=2024-03-02",
		returnedMovies: make([]model.Movie, 20),
		assertions: func(t *testing.T, document *goquery.Document) {
			assert.Equal(t, "/?page=2&genre=123", attribute(document, "#movieGrid .movie:nth-of-type(20)", "data-hx-get"))
			assert.Equal(t, "/?page=2&genre=123", attribute(document, "a#moreMovies", "href"))
		},
	},
}

// attribute returns the value of a specified attribute, of an HTML element that
// is found through the css selector.  If no element is found, the empty string
// is returned
func attribute(document *goquery.Document, selector, attribute string) string {
	val, exists := document.Find(selector).Attr(attribute)
	if !exists {
		val = ""
	}
	return val
}

// mockMtdb is a fake http client that we use to snoop the url that we would be
// passing to the real MTDB, and returns the json for an arbitrary set of Movies
type mockMtdb struct {
	expectedQuery string
	toBeReturned  []model.Movie
	t             *testing.T
}

// Do makes mockMtdb satisfy the HttpAgent interface
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
			handlers.Index(templ, &adapters.Mtdb{Agent: &mockAgent}, FakeCalendar(2024, 3, 2)).ServeHTTP(w, r)
			document, err := goquery.NewDocumentFromReader(w.Body)
			if err != nil {
				t.Error(err)
			}
			test.assertions(t, document)
		})
	}
}
