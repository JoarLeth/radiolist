package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	slowNetTestMsg = "Skipping test in short mode since it makes requests over the internet."
)

func TestSpotifySearchHandlerBadRequest(t *testing.T) {
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://example.com/?artist=nirvana", nil)

	spotifySearchHandler(recorder, request)

	expectedBody := "Bad Request: title and at least one of artist and album must be passed as query parameters.\n"
	actualBody := recorder.Body.String()

	if expectedBody != actualBody {
		t.Errorf("Unexpected content in response body.\nExpected: %#v\nActual: %#v", expectedBody, actualBody)
	}

	expectedStatusCode := http.StatusBadRequest
	actualStatusCode := recorder.Code

	if expectedStatusCode != actualStatusCode {
		t.Errorf("Unexpected status code. Expected: %v, got %v", expectedStatusCode, actualStatusCode)
	}

}

func TestSpotifySearchHandlerNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip(slowNetTestMsg)
	}

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://example.com/?artist=nirvana&title=qwerasdfzxcv", nil)

	spotifySearchHandler(recorder, request)

	expectedBody := "Not Found: Couldn't find a track matching your query.\n"
	actualBody := recorder.Body.String()

	if expectedBody != actualBody {
		t.Errorf("Unexpected content in response body.\nExpected: %#v\nActual: %#v", expectedBody, actualBody)
	}

	expectedStatusCode := http.StatusNotFound
	actualStatusCode := recorder.Code

	if expectedStatusCode != actualStatusCode {
		t.Errorf("Unexpected status code. Expected: %v, got %v", expectedStatusCode, actualStatusCode)
	}
}

func TestSpotifySearchHandlerOK(t *testing.T) {
	if testing.Short() {
		t.Skip(slowNetTestMsg)
	}

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://example.com/?title=come&artist=nirvana", nil)

	spotifySearchHandler(recorder, request)

	expectedBody := "{\"Name\":\"Come As You Are\",\"Artists\":[\"Nirvana\"],\"Album\":\"Nirvana\",\"Href\":\"spotify:track:5r35Zd5Onw3aV3Gm9XdgtI\",\"Territories\":\"AD AR AT AU BE BG BO BR CH CL CO CR CY CZ DE DK DO EC EE ES FI FR GB GR GT HK HN HU IE IS IT LI LT LU LV MC MT MY NI NL NO NZ PA PE PH PL PT PY RO SE SG SI SK SV TR TW UY\"}"
	actualBody := recorder.Body.String()

	if expectedBody != actualBody {
		t.Errorf("Unexpected content in response body.\nExpected: %#v\nActual: %#v", expectedBody, actualBody)
	}

	expectedStatusCode := http.StatusOK
	actualStatusCode := recorder.Code

	if expectedStatusCode != actualStatusCode {
		t.Errorf("Unexpected status code. Expected: %v, got %v", expectedStatusCode, actualStatusCode)
	}
}
