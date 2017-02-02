package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"

	"regexp"

	"net/url"

	"flag"

	"github.com/JalfResi/mustacheHandler"
	"github.com/JalfResi/regexphandler"
)

var (
	config = flag.String("config", "", "Config CSV filename")
)

func sameHost(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = r.URL.Host
		log.Println("source: ", r.URL.String())
		handler.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()

	f, err := os.Open(*config)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(bufio.NewReader(f))

	reHandler := &regexphandler.RegexpHandler{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if len(record) != 3 {
			log.Printf("Config entry must consist of [ regexp-url-match, target-url, mustache-template ]")
			break
		}

		target, _ := url.Parse(record[1])
		proxy := httputil.NewSingleHostReverseProxy(target)

		mHandler := &mustacheHandler.MustacheHandler{}
		mHandler.Handler(record[2], logger("reverseproxy", proxy))

		re := regexp.MustCompile(record[0])
		reHandler.Handler(re, logger("mustache", mHandler))
	}
	f.Close()

	log.Fatal(http.ListenAndServe(":12345", logger("rehandler", sameHost(reHandler))))
}

func logger(prefix string, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Save a copy of this request for debugging.
		requestDump, err := httputil.DumpRequest(r, false)
		if err != nil {
			log.Println(err)
		}
		log.Println(prefix, string(requestDump))

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)

		dump, err := httputil.DumpResponse(rec.Result(), true)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(prefix, string(dump))

		// we copy the captured response headers to our new response
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}

		// grab the captured response body
		data := rec.Body.Bytes()

		w.Write([]byte(data))
	}
}
