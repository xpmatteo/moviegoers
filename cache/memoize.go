package cache

import (
	"github.com/patrickmn/go-cache"
	"github.com/xpmatteo/gomovies/domain"
	"github.com/xpmatteo/gomovies/mtdb"
	"time"
)

type Memoizer struct {
	repository domain.MovieRepository
	cache      *cache.Cache
}

func (m *Memoizer) Query(opts domain.QueryOptions) []domain.Movie {
	url := mtdb.QueryString(opts)
	movies, ok := m.cache.Get(url)
	if !ok {
		movies = m.repository.Query(opts)
		m.cache.Set(url, movies.([]domain.Movie), 60*time.Minute)
	}
	return movies.([]domain.Movie)
}

func NewMemoizer(q domain.MovieRepository) *Memoizer {
	return &Memoizer{
		repository: q,
		cache:      cache.New(60*time.Minute, 60*time.Minute),
	}
}
