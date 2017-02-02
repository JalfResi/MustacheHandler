package mustacheHandler

import (
	"encoding/json"
	"log"
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
	data := rec.Body.Bytes()

	// Response content-type is json, process mustache template
	if rec.Header().Get("Content-Type") == "application/json" {
		var jsonData map[string]interface{}
		err := json.Unmarshal(data, &jsonData)
		if err != nil {
			log.Fatal(err)
		}
		bd := mustache.Render(h.template, jsonData)
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
