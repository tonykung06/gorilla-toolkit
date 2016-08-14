package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/reverse"
)

func main() {
	// basicReverseMatching()
	// compositeMatchers()
	// reverseRegex()
	// routeTemplates()
	htmlTemplate()
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

//formatting a url from regex with args
func reverseRegex() {
	regex, _ := reverse.CompileRegexp(`/foo/(\d+)`)
	r, _ := regex.Revert(url.Values{"": {"42"}})

	namedGroupRegex, _ := reverse.CompileRegexp(`/foo/(?P<bar>\d+)`)
	r2, _ := namedGroupRegex.Revert(url.Values{"bar": {"44"}})

	fmt.Println(r, r2)
}

//a route template that can be used to generate many gorilla routes
func routeTemplates() {
	routeTemplate, _ := reverse.CompileRegexp(`/foo/(?P<bar>.+)`)
	muxRoute1, _ := routeTemplate.Revert(url.Values{"bar": {"{bar:[0-9]+}"}})
	muxRoute2, _ := routeTemplate.Revert(url.Values{"bar": {"42"}})
	fmt.Println(muxRoute1)
	fmt.Println(muxRoute2)
}

func htmlTemplate() {
	routeTemplate, _ := reverse.CompileRegexp(`/foo/(?P<bar>.+)`)
	reverter := func(regex *reverse.Regexp) func(...string) (string, error) {
		return func(params ...string) (string, error) {
			values := url.Values{}
			for i := 0; i < len(params); i += 2 {
				values[params[i]] = []string{params[i+1]}
			}
			return regex.Revert(values)
		}
	}
	fm := template.FuncMap{"pathGen": reverter(routeTemplate)}
	t, _ := template.New("").Funcs(fm).Parse(`{{pathGen "bar" "42"}}`)
	t.Execute(os.Stdout, nil)
}
