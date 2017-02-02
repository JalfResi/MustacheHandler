package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"net/http"
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

		target, err := url.Parse(record[1])
		if err != nil {
			log.Fatal(err)
		}

		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.Host = target.Host
				req.URL = target
			},
		}

		mHandler := &mustacheHandler.MustacheHandler{}
		mHandler.Handler(record[2], proxy)

		re := regexp.MustCompile(record[0])
		reHandler.Handler(re, mHandler)
	}
	f.Close()

	log.Fatal(http.ListenAndServe(":12345", reHandler))
}
