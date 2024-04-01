package web

import (
	"bytes"
	"encoding/xml"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xpmatteo/gomovies/domain"
	"golang.org/x/net/html"
	"html/template"
	"io"
	"strings"
	"testing"
)

var testCases = []struct {
	name                   string
	movies                 []domain.Movie
	options                domain.QueryOptions
	selector               string
	matches                []string
	attribute              string
	expectedAttributeValue string
}{
	{
		name: "movie titles",
		movies: []domain.Movie{
			{Title: "Foobar"},
			{Title: "Zork"},
			{Title: "Blah"},
		},
		selector: "#movieGrid .movie h3",
		matches:  []string{"Foobar", "Zork", "Blah"},
	},
	{
		name: "movie overviews",
		movies: []domain.Movie{
			{Overview: "Something"},
			{Overview: "Something else"},
		},
		selector: "#movieGrid .movie .overview",
		matches:  []string{"Something", "Something else"},
	},
	{
		name: "poster",
		movies: []domain.Movie{
			{PosterPath: "foo.jpg"},
		},
		selector:               "#movieGrid .movie img",
		attribute:              "src",
		expectedAttributeValue: "https://image.tmdb.org/t/p/w185/foo.jpg",
	},
	{
		name: "poster missing",
		movies: []domain.Movie{
			{PosterPath: ""},
		},
		selector:               "#movieGrid .movie img",
		attribute:              "src",
		expectedAttributeValue: "images/no_poster_available.jpg",
	},
}

func Test_allDynamicFeatures(t *testing.T) {
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if test.movies == nil {
				test.movies = []domain.Movie{}
			}
			if test.matches == nil {
				test.matches = []string{""}
			}

			buf := renderTemplate(test.movies, test.options)

			assertWellFormedHTML(t, buf)
			document := parseHtml(t, buf)
			selection := document.Find(test.selector)
			require.Equal(t, len(test.matches), len(selection.Nodes), "unexpected # of matches")
			for i, node := range selection.Nodes {
				assert.Equal(t, test.matches[i], text(node))
			}
			if test.attribute != "" {
				node := selection.Nodes[0]
				attributeValue := findAttributeValue(node, test.attribute)
				assert.Equal(t, test.expectedAttributeValue, attributeValue)
			}
		})
	}
}

func findAttributeValue(node *html.Node, attribute string) string {
	for _, a := range node.Attr {
		if a.Key == attribute {
			return a.Val
		}
	}
	return ""
}

func text(node *html.Node) string {
	// A little mess due to the fact that goquery has
	// a .Text() method on Selection but not on html.Node
	sel := goquery.Selection{Nodes: []*html.Node{node}}
	return strings.TrimSpace(sel.Text())
}

func parseHtml(t *testing.T, buf bytes.Buffer) *goquery.Document {
	document, err := goquery.NewDocumentFromReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("error rendering template %s", err)
	}
	return document
}

func assertWellFormedHTML(t *testing.T, buf bytes.Buffer) {
	decoder := xml.NewDecoder(bytes.NewReader(buf.Bytes()))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	decoder.Entity = xml.HTMLEntity

	for {
		_, err := decoder.Token()
		switch err {
		case io.EOF:
			return // We're done, it's valid!
		case nil:
			// do nothing
		default:
			t.Fatalf("error parsing html: %s", err)
		}
	}
}

const indexFilename = "index.tmpl"

func renderTemplate(movies []domain.Movie, opts domain.QueryOptions) bytes.Buffer {
	templ := template.Must(template.ParseFiles(indexFilename))
	var buf bytes.Buffer
	data := Model(movies, opts)
	err := templ.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf
}
