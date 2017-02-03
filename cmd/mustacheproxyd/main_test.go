package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"regexp"
	"testing"

	"github.com/JalfResi/mustacheHandler"
	"github.com/JalfResi/regexphandler"
)

func TestHandlers(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	// Create a request to pass to our handler. We dont have any query parameters
	// for now, so we'll pass 'nil' as the third parameter
	req, err := http.NewRequest("GET", "/users/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// proxy target server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"message": "Hello, world!"}`)
	}))
	defer ts.Close()

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	// configure handler chain
	target, _ := url.Parse(ts.URL)
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = target.Host
			req.URL = target
		},
	}

	re := regexp.MustCompile("/users/(.*)")

	mHandler := &mustacheHandler.MustacheHandler{}
	mHandler.Handler(re, "./templates/$1.mustache", proxy)

	reHandler := &regexphandler.RegexpHandler{}
	reHandler.Handler(re, mHandler)

	// Call handler with request and response recorder
	reHandler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `<html><body>Hello, world!</body></html>`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
