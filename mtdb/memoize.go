package mtdb

import (
	"github.com/xpmatteo/gomovies/domain"
	"log"
)

type Memoizer struct {
	Mtdb  *Mtdb
	pages map[string][]domain.Movie
}

func (m *Memoizer) Query(opts domain.QueryOptions) []domain.Movie {
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
		pages: make(map[string][]domain.Movie),
	}
}
