package adapters

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/gomovies/model"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func Test_queryString(t *testing.T) {
	tests := []struct {
		name     string
		opts     model.QueryOptions
		expected string
	}{
		{
			name:     "empty",
			opts:     model.QueryOptions{},
			expected: "",
		},
		{
			name:     "genre",
			opts:     model.QueryOptions{Genre: 123},
			expected: "with_genres=123",
		},
		{
			name:     "page",
			opts:     model.QueryOptions{Page: 3},
			expected: "page=3",
		},
		{
			name:     "genre & page",
			opts:     model.QueryOptions{Genre: 44, Page: 2},
			expected: "with_genres=44&page=2",
		},
		{
			name:     "release date",
			opts:     model.QueryOptions{ReleaseDateMax: time.Date(2024, 1, 22, 0, 0, 0, 0, time.Local)},
			expected: "primary_release_date.lte=2024-01-22",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, QueryString(test.opts))
		})
	}
}

var m = Mtdb{Agent: http.DefaultClient}

func Test_live(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	movies := m.Query(model.QueryOptions{})

	assert.Len(t, movies, 20)
}

func Test_NoApiKey(t *testing.T) {
	if err := os.Unsetenv("TMDB_ACCESS_TOKEN"); err != nil {
		t.Fatal(err)
	}

	assert.Panics(t, func() { m.Query(model.QueryOptions{}) })
}
