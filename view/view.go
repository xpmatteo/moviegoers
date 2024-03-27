package view

import (
	"github.com/xpmatteo/gomovies/model"
)

func Model(movies []model.Movie, options model.QueryOptions) map[string]any {
	return map[string]any{
		"movies":        movies,
		"nextPage":      options.Page + 1,
		"genres":        model.Genres,
		"selectedGenre": options.Genre,
	}
}
