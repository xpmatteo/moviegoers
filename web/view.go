package web

import (
	"github.com/xpmatteo/gomovies/domain"
)

func Model(movies []domain.Movie, options domain.QueryOptions) map[string]any {
	return map[string]any{
		"movies":        movies,
		"nextPage":      options.Page + 1,
		"genres":        domain.Genres,
		"selectedGenre": options.Genre,
	}
}
