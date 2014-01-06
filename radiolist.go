package main

import (
	"encoding/json"
	"fmt"
	"github.com/joarleth/spotify/track"
	"net/http"
	"net/url"
	"strings"
)

type searchParameters struct {
	title  string
	artist string
	album  string
}

func spotifySearchHandler(w http.ResponseWriter, r *http.Request) {
	s := track.NewSearcher("se")
	w.Write(getTrackJson(r.URL, s))
}

func getTrackJson(request_url *url.URL, track_searcher *track.Searcher) []byte {
	title, artist, album := getSearchParameters(request_url)
	track, _ := track_searcher.FindClosestMatch(title, artist, album)

	fmt.Print(track)

	json_track, _ := json.Marshal(track)

	return json_track
}

func getSearchParameters(url *url.URL) (title string, artist string, album string) {
	queries := url.Query()

	title = queries.Get("title")
	artist = queries.Get("artist")
	album = queries.Get("album")

	title = strings.TrimSpace(title)
	artist = strings.TrimSpace(artist)
	album = strings.TrimSpace(album)

	return
}

func main() {
	http.HandleFunc("/", spotifySearchHandler)
	http.ListenAndServe(":8080", nil)
}