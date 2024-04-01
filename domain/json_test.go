package domain

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_singleMovie(t *testing.T) {
	const singleMovie = `{
      "overview": "Po is gearing up",
      "title": "Kung Fu Panda 4"
    }`
	var movie Movie

	err := json.Unmarshal([]byte(singleMovie), &movie)
	if err != nil {
		t.Error(err)
	}

	expected := Movie{
		Title:    "Kung Fu Panda 4",
		Overview: "Po is gearing up",
	}
	assert.Equal(t, expected, movie)
}
