package domain

type Movie struct {
	Title      string `json:"title"`
	Overview   string `json:"overview"`
	PosterPath string `json:"poster_path"`
}
