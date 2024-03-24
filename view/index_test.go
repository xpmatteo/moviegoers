package view

import (
	"bytes"
	"encoding/xml"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xpmatteo/gomovies/model"
	"golang.org/x/net/html"
	"html/template"
	"io"
	"strings"
	"testing"
)

var testCases = []struct {
	name     string
	movies   []*model.Movie
	path     string
	selector string
	matches  []string
}{
	{
		name: "movie titles",
		movies: []*model.Movie{
			{Title: "Foobar"},
			{Title: "Zork"},
			{Title: "Blah"},
		},
		selector: "#movieGrid .movie h3",
		matches:  []string{"Foobar", "Zork", "Blah"},
	},
	{
		name: "movie overviews",
		movies: []*model.Movie{
			{Overview: "Something"},
			{Overview: "Something else"},
		},
		selector: "#movieGrid .movie .overview",
		matches:  []string{"Something", "Something else"},
	},
}

func Test_allDynamicFeatures(t *testing.T) {
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if test.movies == nil {
				test.movies = []*model.Movie{}
			}

			buf := renderTemplate(test.movies, test.path)

			assertWellFormedHTML(t, buf)
			document := parseHtml(t, buf)
			selection := document.Find(test.selector)
			require.Equal(t, len(test.matches), len(selection.Nodes), "unexpected # of matches")
			for i, node := range selection.Nodes {
				assert.Equal(t, test.matches[i], text(node))
			}
		})
	}
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

const indexFilename = "index.gotmpl"

func renderTemplate(movies []*model.Movie, path string) bytes.Buffer {
	templ := template.Must(template.ParseFiles(indexFilename))
	var buf bytes.Buffer
	data := map[string]any{
		"movies": movies,
		"path":   path,
	}
	err := templ.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf
}
