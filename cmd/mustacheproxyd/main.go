package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	"regexp"

	"net/url"

	"github.com/JalfResi/mustacheHandler"
	"github.com/JalfResi/regexphandler"
)

func main() {
	target, _ := url.Parse("http://service.example.com/")
	proxy := httputil.NewSingleHostReverseProxy(target)

	mHandler := &mustacheHandler.MustacheHandler{}
	mHandler.Handler("<html><body>{{message}}</body></html>", proxy)

	re := regexp.MustCompile("/users")
	reHandler := &regexphandler.RegexpHandler{}
	reHandler.Handler(re, mHandler)

	log.Fatal(http.ListenAndServe(":12345", reHandler))
}
