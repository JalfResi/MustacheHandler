package mustacheHandler

import (
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/hoisie/mustache"
)

type MustacheHandler struct {
	handler  http.Handler
	template string
}

func (h *MustacheHandler) Handler(template string, handler http.Handler) {
	h.handler = handler
	h.template = template
}

func (h *MustacheHandler) HandleFunc(template string, handler func(http.ResponseWriter, *http.Request)) {
	h.handler = http.HandlerFunc(handler)
	h.template = template
}

func (h *MustacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Capture the response
	rec := httptest.NewRecorder()
	h.handler.ServeHTTP(rec, r)

	// we copy the captured response headers to our new response
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}

	// grab the captured response
	originalData := rec.Body.Bytes()
	data := originalData

	if rec.Header().Get("Content-Type") == "application/json" {
		// unmarshall json
		data = mustache.Render(h.template, jsonObj)
	}

	// But the Content-Length might have been set already,
	// we should modify it by adding the length
	// of our own data.
	// Ignoring the error is fine here:
	// if Content-Length is empty or otherwise invalid,
	// Atoi() will return zero,
	// which is just what we'd want in that case.
	clen, _ := strconv.Atoi(rec.Header().Get("Content-Length"))
	clen += len(data)
	w.Header().Set("Content-Length", strconv.Itoa(clen))

	// write out our modified response
	w.Write([]byte(data))
}