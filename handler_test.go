package mustacheHandler

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHostMatcherHandler(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	// Create a request to pass to our handler. We dont have any query parameters
	// for now, so we'll pass 'nil' as the third parameter
	req, err := http.NewRequest("GET", "/users/ben", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	handler := &MustacheHandler{}

	// original proxy request response
	handler.HandleFunc("<html><body><dl><dt>user</dt><dd>{{user}}</dd></dl></body></html>", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"user": "ben"}`)
	})

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `<html><body><dl><dt>user</dt><dd>ben</dd></dl></body></html>`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
