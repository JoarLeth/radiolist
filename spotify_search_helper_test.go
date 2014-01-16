package main

import (
	"encoding/json"
	"errors"
	"github.com/joarleth/spotify/track"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpotifySearchHelperReturnsJsonTrack(t *testing.T) {
	recorder := httptest.NewRecorder()
	track := track.Track{
		Name:        "Test Name",
		Artists:     []string{"Test Artist"},
		Album:       "Test Album",
		Href:        "spotify:track:foo",
		Territories: "SE",
	}
	s := mock_searcher{err: nil, track: track}

	spotifySearchHelper(recorder, s, "test", "test", "test")

	expectedBodyBytes, _ := json.Marshal(track)

	expectedBody := string(expectedBodyBytes)
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

func TestSpotifySearchHelperUnexpectedError(t *testing.T) {
	recorder := httptest.NewRecorder()
	unexpected_error := track.TrackError{Msg: "Foo", ErrorType: track.UnexpectedError}
	track := track.Track{}
	s := mock_searcher{err: unexpected_error, track: track}

	spotifySearchHelper(recorder, s, "test", "test", "test")

	expectedBody := "Internal Server Error: An unexpected error occurred.\n"
	actualBody := recorder.Body.String()

	if expectedBody != actualBody {
		t.Errorf("Unexpected content in response body.\nExpected: %#v\nActual: %#v", expectedBody, actualBody)
	}

	expectedStatusCode := http.StatusInternalServerError
	actualStatusCode := recorder.Code

	if expectedStatusCode != actualStatusCode {
		t.Errorf("Unexpected status code. Expected: %v, got %v", expectedStatusCode, actualStatusCode)
	}
}

func TestSpotifySearchHelperRateLimitError(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := track.TrackError{Msg: "Foo", ErrorType: track.RateLimitError}
	track := track.Track{}
	s := mock_searcher{err: err, track: track}

	spotifySearchHelper(recorder, s, "test", "test", "test")

	expectedBody := "Service Unavailable: Too many requests to spotify at this time. Please come back another time.\n"
	actualBody := recorder.Body.String()

	if expectedBody != actualBody {
		t.Errorf("Unexpected content in response body.\nExpected: %#v\nActual: %#v", expectedBody, actualBody)
	}

	expectedStatusCode := http.StatusServiceUnavailable
	actualStatusCode := recorder.Code

	if expectedStatusCode != actualStatusCode {
		t.Errorf("Unexpected status code. Expected: %v, got %v", expectedStatusCode, actualStatusCode)
	}
}

func TestSpotifySearchHelperExternalServiceError(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := track.TrackError{Msg: "Foo", ErrorType: track.ExternalServiceError}
	track := track.Track{}
	s := mock_searcher{err: err, track: track}

	spotifySearchHelper(recorder, s, "test", "test", "test")

	expectedBody := "Service Unavailable: The Spotify service used by this API bahaves unexpectedly.\n"
	actualBody := recorder.Body.String()

	if expectedBody != actualBody {
		t.Errorf("Unexpected content in response body.\nExpected: %#v\nActual: %#v", expectedBody, actualBody)
	}

	expectedStatusCode := http.StatusServiceUnavailable
	actualStatusCode := recorder.Code

	if expectedStatusCode != actualStatusCode {
		t.Errorf("Unexpected status code. Expected: %v, got %v", expectedStatusCode, actualStatusCode)
	}
}

func TestSpotifySearchHelperArgumentError(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := track.TrackError{Msg: "Foo", ErrorType: track.ArgumentError}
	track := track.Track{}
	s := mock_searcher{err: err, track: track}

	spotifySearchHelper(recorder, s, "this", "is", "irrelevant")

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

func TestSpotifySearchHelperNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()
	track := track.Track{}
	s := mock_searcher{track: track}

	spotifySearchHelper(recorder, s, "this", "is", "irrelevant")

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

func TestSpotifySearchHelperNoErrorType(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := track.TrackError{Msg: "Foo"}
	track := track.Track{}
	s := mock_searcher{err: err, track: track}

	spotifySearchHelper(recorder, s, "this", "is", "irrelevant")

	expectedBody := "Internal Server Error: This should never happen.\n"
	actualBody := recorder.Body.String()

	if expectedBody != actualBody {
		t.Errorf("Unexpected content in response body.\nExpected: %#v\nActual: %#v", expectedBody, actualBody)
	}

	expectedStatusCode := http.StatusInternalServerError
	actualStatusCode := recorder.Code

	if expectedStatusCode != actualStatusCode {
		t.Errorf("Unexpected status code. Expected: %v, got %v", expectedStatusCode, actualStatusCode)
	}
}

func TestSpotifySearchHelperErrorString(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := errors.New("This error is really unexpected.")
	track := track.Track{}
	s := mock_searcher{err: err, track: track}

	spotifySearchHelper(recorder, s, "test", "test", "test")

	expectedBody := "Internal Server Error: This should never happen.\n"
	actualBody := recorder.Body.String()

	if expectedBody != actualBody {
		t.Errorf("Unexpected content in response body.\nExpected: %#v\nActual: %#v", expectedBody, actualBody)
	}

	expectedStatusCode := http.StatusInternalServerError
	actualStatusCode := recorder.Code

	if expectedStatusCode != actualStatusCode {
		t.Errorf("Unexpected status code. Expected: %v, got %v", expectedStatusCode, actualStatusCode)
	}
}

type mock_searcher struct {
	err   error
	track track.Track
}

func (ms mock_searcher) FindClosestMatch(string, string, string) (track.Track, error) {
	return ms.track, ms.err
}
