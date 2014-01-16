package main

import (
	"encoding/json"
	"github.com/joarleth/spotify/track"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type searchParameters struct {
	title  string
	artist string
	album  string
}

type searcherInterface interface {
	FindClosestMatch(string, string, string) (track.Track, error)
}

func spotifySearchHandler(w http.ResponseWriter, r *http.Request) {
	s := track.NewSearcher("se")
	title, artist, album := getSearchParameters(r.URL)

	spotifySearchHelper(w, s, title, artist, album)
}

func spotifySearchHelper(w http.ResponseWriter, s searcherInterface, title, artist, album string) {
	t, err := s.FindClosestMatch(title, artist, album)
	if err != nil {
		if terr, isTrackError := err.(track.TrackError); isTrackError {
			switch terr.ErrorType {
			case track.ArgumentError:
				http.Error(w, "Bad Request: title and at least one of artist and album must be passed as query parameters.", http.StatusBadRequest)
				return
			case track.UnexpectedError:
				http.Error(w, "Internal Server Error: An unexpected error occurred.", http.StatusInternalServerError)
				return
			case track.RateLimitError:
				http.Error(w, "Service Unavailable: Too many requests to spotify at this time. Please come back another time.", http.StatusServiceUnavailable)
				return
			case track.ExternalServiceError:
				http.Error(w, "Service Unavailable: The Spotify service used by this API bahaves unexpectedly.", http.StatusServiceUnavailable)
				return
			}
		}
		http.Error(w, "Internal Server Error: This should never happen.", http.StatusInternalServerError)
		return
	}

	if reflect.DeepEqual(t, track.Track{}) {
		http.Error(w, "Not Found: Couldn't find a track matching your query.", http.StatusNotFound)
		return
	}

	json_result, merr := json.Marshal(t)

	if merr != nil {
		http.Error(w, "Internal Server Error: Unable to serialize track.", http.StatusInternalServerError)
		return
	}

	w.Write(json_result)
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
