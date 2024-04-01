package cache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/gomovies/domain"
	"testing"
)

type FakeMovieRepository struct {
	receivedOptions domain.QueryOptions
	numCalls        int
}

func (f *FakeMovieRepository) Query(options domain.QueryOptions) []domain.Movie {
	f.receivedOptions = options
	f.numCalls++
	return []domain.Movie{
		{
			Title: fmt.Sprintf("Movie for genre %d", +options.Genre),
		},
	}
}

func Test_cacheMiss(t *testing.T) {
	repository := FakeMovieRepository{}
	memoizer := NewMemoizer(&repository)

	movies := memoizer.Query(domain.QueryOptions{Genre: 1234})

	assert.Equal(t, "Movie for genre 1234", movies[0].Title)
	assert.Equal(t, 1, repository.numCalls)
	assert.Equal(t, domain.QueryOptions{Genre: 1234}, repository.receivedOptions)
}

func Test_cacheHit(t *testing.T) {
	repository := FakeMovieRepository{}
	memoizer := NewMemoizer(&repository)
	memoizer.Query(domain.QueryOptions{Genre: 1234})

	movies := memoizer.Query(domain.QueryOptions{Genre: 1234})

	assert.Equal(t, "Movie for genre 1234", movies[0].Title)
	assert.Equal(t, 1, repository.numCalls)
	assert.Equal(t, domain.QueryOptions{Genre: 1234}, repository.receivedOptions)
}

func Test_multipleHits(t *testing.T) {
	repository := FakeMovieRepository{}
	memoizer := NewMemoizer(&repository)
	memoizer.Query(domain.QueryOptions{Genre: 111})
	memoizer.Query(domain.QueryOptions{Genre: 222})
	memoizer.Query(domain.QueryOptions{Genre: 333})

	movies111 := memoizer.Query(domain.QueryOptions{Genre: 111})
	movies222 := memoizer.Query(domain.QueryOptions{Genre: 222})
	movies333 := memoizer.Query(domain.QueryOptions{Genre: 333})

	assert.Equal(t, "Movie for genre 111", movies111[0].Title)
	assert.Equal(t, "Movie for genre 222", movies222[0].Title)
	assert.Equal(t, "Movie for genre 333", movies333[0].Title)
	assert.Equal(t, 3, repository.numCalls)
}
