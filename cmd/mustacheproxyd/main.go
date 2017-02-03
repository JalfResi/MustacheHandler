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

Guard RegExp URL, Target URL, Mustache Template filename

e.g.
/services/(.*)/users/(.*),https://$1,example.com/users/$2,./templates/$1.mustache

Results in:

/services/reporter/users/ben -> https://reporter,example.com/users/ben
./templates/reporter.mustache

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

		re := regexp.MustCompile(record[0])

		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				target, err := url.Parse(re.ReplaceAllString(req.URL.String(), record[1]))
				if err != nil {
					log.Fatal(err)
				}
				log.Println(req.URL.String(), "->", target.String())

				req.Host = target.Host
				req.URL = target
			},
		}

		mHandler := &mustacheHandler.MustacheHandler{}
		mHandler.Handler(re, record[2], proxy)

		reHandler.Handler(re, mHandler)
	}
	f.Close()

	log.Fatal(http.ListenAndServe(*hostAddr, reHandler))
}
