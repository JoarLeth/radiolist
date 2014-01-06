package main

import (
	"net/url"
	"testing"
)

func TestGetSearchParameters(t *testing.T) {
	expected_title := "human behaviour"
	expected_artist := "björk"
	expected_album := ""

	url_string := "http://localhost:8080/tracks/search?artist=" + url.QueryEscape(expected_artist) + "&title=" + url.QueryEscape(expected_title)

	url, _ := url.Parse(url_string)

	actual_title, actual_artist, actual_album := getSearchParameters(url)

	if actual_title != "human behaviour" {
		t.Errorf("Unexpected title. Expected: %s, Actual: %s", expected_title, actual_title)
	}

	if actual_artist != "björk" {
		t.Errorf("Unexpected artist. Expected: %s, Actual: %s", expected_artist, actual_artist)
	}

	if actual_album != "" {
		t.Errorf("Unexpected album. Expected: %s, Actual: %s", expected_album, actual_album)
	}
}
