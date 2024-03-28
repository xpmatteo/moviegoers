package adapters

import (
	"github.com/xpmatteo/gomovies/model"
	"log"
)

type Memoizer struct {
	Mtdb  *Mtdb
	pages map[string][]model.Movie
}

func (m *Memoizer) Query(opts model.QueryOptions) []model.Movie {
	url := QueryString(opts)
	movies, ok := m.pages[url]
	if !ok {
		log.Print("Not OK:", url)
		movies = m.Mtdb.Query(opts)
		m.pages[url] = movies
	}
	return movies
}

func NewMemoizer(m *Mtdb) *Memoizer {
	return &Memoizer{
		Mtdb:  m,
		pages: make(map[string][]model.Movie),
	}
}
