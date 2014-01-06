package main

import (
	"encoding/json"
	"github.com/joarleth/spotify/track"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func TestGetTrackJson(t *testing.T) {
	xml_data := getTextFileData(t, "tracks2.xml")

	mockserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write(xml_data)
	}))

	s := new_mock_searcher("SE", mockserver.URL)

	expected_track := track.Track{
		Name:        "Uncover",
		Artists:     []string{"Zara Larsson"},
		Album:       "Introducing",
		Href:        "spotify:track:131l5GkXPIk81bxihGypPt",
		Territories: "SE",
	}

	expected, _ := json.Marshal(expected_track)

	url, _ := url.Parse("http://localhost:8080/tracks/search?artist=" + url.QueryEscape("Zara Larsson") + "&title=" + url.QueryEscape("Uncover"))

	actual := getTrackJson(url, s)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Resulting track json not matching expected.\nExpected: %#v\nActual: %#v", string(expected), string(actual))
	}
}

func getTextFileData(t *testing.T, filename string) []byte {
	file, open_file_error := os.Open(filename)
	defer file.Close()

	if open_file_error != nil {
		t.Fatalf("Failed to open file. Error: %v", open_file_error.Error())
	}

	data, read_file_error := ioutil.ReadAll(file)

	if read_file_error != nil {
		t.Fatalf("Failed to read file. Error: %v", read_file_error.Error())
	}

	return data
}

func new_mock_searcher(territory, search_url string) *track.Searcher {
	return &track.Searcher{
		Territory:             territory,
		Track_search_base_url: search_url,
	}
}
