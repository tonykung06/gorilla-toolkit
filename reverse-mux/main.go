package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/reverse"
)

func main() {
	// basicReverseMatching()
	compositeMatchers()
}

func compositeMatchers() {
	u := &url.URL{
		Scheme:   "http",
		Host:     "localhost:9999",
		Path:     "/foo/33",
		RawQuery: "buz=42",
	}
	r := &http.Request{URL: u}
	p, _ := reverse.NewRegexpPath("/foo/[0-9]+")
	q := reverse.NewQuery(map[string]string{"buz": ""})
	//and
	a := reverse.NewAll([]reverse.Matcher{p, q})

	//or
	o := reverse.NewOne([]reverse.Matcher{p, q})
	fmt.Println("All Match?", a.Match(r))
	fmt.Println("Or Match?", o.Match(r))
}

func basicReverseMatching() {
	u := &url.URL{
		Scheme:   "http",
		Host:     "localhost:9999",
		Path:     "/foo/33",
		RawQuery: "buz=42",
	}
	r := &http.Request{URL: u}
	p, _ := reverse.NewGorillaPath("/foo/{param:[0-9]+}", false)
	fmt.Println("Match?", p.Match(r))
}
