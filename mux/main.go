package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// usingSchemeMatchers()
	// usingPathPrefixMatchers()
	// usingFullPathMatchers()
	// usingMethodMatchers()
	// usingQueryMatchers()
	usingSubrouters()
}

func usingSubrouters() {
	r := mux.NewRouter()
	func1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("func 1"))
	}
	func2 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("func 2"))
	}
	r.HandleFunc("/", func1)
	sr := r.PathPrefix("/foo").Subrouter()
	sr.HandleFunc("/", func2)
	http.Handle("/", r)
	go http.ListenAndServe(":4000", nil)
	fmt.Scanln()
}

func usingQueryMatchers() {
	r := mux.NewRouter()
	func1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("func 1"))
	}
	func2 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("func 2"))
	}
	r.Path("/").Queries("foo", "bar").HandlerFunc(func1)
	r.Path("/").Queries("bar", "{bar:[0-9]+}").HandlerFunc(func2)
	http.Handle("/", r)
	go http.ListenAndServe(":4000", nil)
	fmt.Scanln()
}

func usingMethodMatchers() {
	r := mux.NewRouter()
	func1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
				<head></head>
				<body>
					<form action='' method='POST'>
					The field <input name='field' />
					</form>
				</body>
			</html>
		`))
	}
	func2 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.FormValue("field")))
	}
	r.HandleFunc("/", func1).Methods("GET")
	r.HandleFunc("/", func2).Methods("POST")
	http.Handle("/", r)
	go http.ListenAndServe(":4000", nil)
	fmt.Scanln()
}

func usingFullPathMatchers() {
	r := mux.NewRouter()
	func1 := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Write([]byte(vars["id"]))
	}
	r.HandleFunc("/foo/{id:[0-9]+}", func1)
	http.Handle("/", r)
	go http.ListenAndServe(":4000", nil)
	fmt.Scanln()
}

func usingPathPrefixMatchers() {
	r := mux.NewRouter()
	func1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("function 1"))
	}
	func2 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("function 2"))
	}
	r.PathPrefix("/foo").HandlerFunc(func1)
	r.PathPrefix("/bar").HandlerFunc(func2)
	http.Handle("/", r)
	go http.ListenAndServe(":4000", nil)
	fmt.Scanln()
}

func usingSchemeMatchers() {
	//.Scheme() is disappointing
	r := mux.NewRouter()
	func1 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("function 1"))
	}
	func2 := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("function 2"))
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			func1(w, r)
		} else {
			func2(w, r)
		}
	}
	r.HandleFunc("/", handler)
	http.Handle("/", r)
	go http.ListenAndServe(":4000", nil)
	//see https://gist.github.com/denji/12b3a568f092ab951456
	go http.ListenAndServeTLS(":4443", "server.pem", "server.key", nil)
	fmt.Scanln()
}
