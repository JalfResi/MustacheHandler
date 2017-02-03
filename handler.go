package mustacheHandler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"

	"github.com/hoisie/mustache"
)

type MustacheHandler struct {
	handler  http.Handler
	pattern  *regexp.Regexp
	template string
}

func (h *MustacheHandler) Handler(pattern *regexp.Regexp, template string, handler http.Handler) {
	h.handler = handler
	h.pattern = pattern
	h.template = template
}

func (h *MustacheHandler) HandleFunc(pattern *regexp.Regexp, template string, handler func(http.ResponseWriter, *http.Request)) {
	h.handler = http.HandlerFunc(handler)
	h.pattern = pattern
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
	data := rec.Body.Bytes()

	// Response content-type is json, process mustache template
	if rec.Header().Get("Content-Type") == "application/json" {
		var jsonData map[string]interface{}
		err := json.Unmarshal(data, &jsonData)
		if err != nil {
			log.Fatal(err)
		}

		filename := h.pattern.ReplaceAllString(r.URL.String(), h.template)
		bd := mustache.RenderFile(filename, jsonData)
		data = []byte(bd)

		w.Header().Set("Content-Type", "text/html")

		// Ignoring the error is fine here:
		// if Content-Length is empty or otherwise invalid,
		// Atoi() will return zero,
		// which is just what we'd want in that case.
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	}

	// write out our modified response
	w.Write([]byte(data))
}
