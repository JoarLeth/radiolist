package main

import (
	"encoding/json"
	"github.com/joarleth/spotify/track"
	"net/url"
	"reflect"
	"testing"
)

func TestGetTrackJsonUsingInterface(t *testing.T) {
	s := mock_searcher{}

	expected_track := track.Track{
		Name:        "Uncover",
		Artists:     []string{"Zara Larsson"},
		Album:       "Introducing",
		Href:        "spotify:track:131l5GkXPIk81bxihGypPt",
		Territories: "SE",
	}
	url, _ := url.Parse("http://localhost:8080/tracks/search?foo=bar")

	expected, _ := json.Marshal(expected_track)
	actual, _ := getTrackJson(url, s)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Resulting track json not matching expected.\nExpected: %#v\nActual: %#v", string(expected), string(actual))
	}
}

type mock_searcher struct{}

func (m mock_searcher) FindClosestMatch(title, artist, album string) (track.Track, error) {
	track := track.Track{
		Name:        "Uncover",
		Artists:     []string{"Zara Larsson"},
		Album:       "Introducing",
		Href:        "spotify:track:131l5GkXPIk81bxihGypPt",
		Territories: "SE",
	}

	return track, nil
}
