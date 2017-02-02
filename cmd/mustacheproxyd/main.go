package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"

	"github.com/JalfResi/mustacheHandler"
	"github.com/JalfResi/regexphandler"
)

/*

Config file is CSV file with the following structure:

Guard RegExp URL, Target URL, Mustache Template

e.g.
/users/.*,https://ip-ranges.amazonaws.com/ip-ranges.json,<body><h1>Sync Token: {{syncToken}}</h1></body>

*/

var (
	hostAddr = flag.String("host", "127.0.0.1:12345", "Hostname and port")
	config   = flag.String("config", "", "Config CSV filename")
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

	log.Fatal(http.ListenAndServe(*hostAddr, reHandler))
}
