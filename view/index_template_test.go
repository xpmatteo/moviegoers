package web

import (
	"bytes"
	"demo-template-test/todo"
	"encoding/xml"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"html/template"
	"io"
	"strings"
	"testing"
)

var testCases = []struct {
	name     string
	model    *todo.List
	path     string
	selector string
	matches  []string
}{
	{
		name: "all todo items are shown",
		model: todo.NewList().
			Add("Foo").
			Add("Bar"),
		selector: "ul.todo-list li",
		matches:  []string{"Foo", "Bar"},
	},
	{
		name: "completed items get the 'completed' class",
		model: todo.NewList().
			Add("Foo").
			AddCompleted("Bar"),
		selector: "ul.todo-list li.completed",
		matches:  []string{"Bar"},
	},
	{
		name: "items left",
		model: todo.NewList().
			Add("One").
			Add("Two").
			AddCompleted("Three"),
		selector: "span.todo-count",
		matches:  []string{"2 items left"},
	},
	{
		name:     "highlighted navigation link: All",
		path:     "/",
		selector: "ul.filters a.selected",
		matches:  []string{"All"},
	},
	{
		name:     "highlighted navigation link: Active",
		path:     "/active",
		selector: "ul.filters a.selected",
		matches:  []string{"Active"},
	},
	{
		name:     "highlighted navigation link: Completed",
		path:     "/completed",
		selector: "ul.filters a.selected",
		matches:  []string{"Completed"},
	},
}

func Test_allDynamicFeatures(t *testing.T) {
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if test.model == nil {
				test.model = todo.NewList()
			}

			buf := renderTemplate(test.model, test.path)

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

func renderTemplate(model *todo.List, path string) bytes.Buffer {
	templ := template.Must(template.ParseFiles("index.gotmpl"))
	var buf bytes.Buffer
	data := map[string]any{
		"model": model,
		"path":  path,
	}
	err := templ.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf
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
