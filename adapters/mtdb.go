package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/xpmatteo/gomovies/model"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Mtdb struct {
}

func QueryString(o model.QueryOptions) string {
	params := make([]string, 0)
	if o.Genre != 0 {
		params = append(params, fmt.Sprintf("with_genres=%d", o.Genre))
	}
	if o.Page != 0 {
		params = append(params, fmt.Sprintf("page=%d", o.Page))
	}
	if !o.ReleaseDateMax.IsZero() {
		s := fmt.Sprintf("primary_release_date.lte=%d-%02d-%02d",
			o.ReleaseDateMax.Year(),
			o.ReleaseDateMax.Month(),
			o.ReleaseDateMax.Day())
		params = append(params, s)
	}
	return strings.Join(params, "&")
}

func (m *Mtdb) Query(opts model.QueryOptions) []model.Movie {
	apiKey := os.Getenv("TMDB_ACCESS_TOKEN")
	if apiKey == "" {
		panic("missing env var TMDB_ACCESS_TOKEN")
	}
	url := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc"
	url = "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&primary_release_date.lte=2024-03-24&sort_by=primary_release_date.desc"
	url = "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&primary_release_date.lte=2024-03-24&sort_by=primary_release_date.desc&vote_count.gte=100"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Print(err)
	}
	request.Header.Add("accept", "application/json")
	request.Header.Add("Authorization", "Bearer "+apiKey)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Print(err)
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
	}

	var data struct {
		Results []model.Movie `json:"Results"`
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Print("Error:", err)
	}

	return data.Results
}
