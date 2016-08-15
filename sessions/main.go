package main

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(64))
var fileSystemStore = sessions.NewFilesystemStore("/test_session", securecookie.GenerateRandomKey(64))

func main() {
	// sessionDataStoredInCookie()
	// sessionDataStoredInFileSystem()
	flashMessages()
}

func sessionDataStoredInFileSystem() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, _ := fileSystemStore.Get(r, "session-name")
		session.Values["favorite-song"] = []string{"My Heart will go on", "Love"}
		session.Save(r, w)
		w.Write([]byte("<h1>Hello Sessions!</h1>"))
	})
	http.HandleFunc("/sessioncontents", func(w http.ResponseWriter, r *http.Request) {
		session, _ := fileSystemStore.Get(r, "session-name")
		songs := session.Values["favorite-song"].([]string)
		w.Header().Add("Content-Type", "text/html")
		for _, song := range songs {
			w.Write([]byte(song + "<br />"))
		}
	})
	//seems there are also some other places leaking memory
	http.ListenAndServe(":3000", context.ClearHandler(http.DefaultServeMux))
}

func sessionDataStoredInCookie() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		session.Values["favorite-song"] = []string{"My Heart will go on", "Love"}
		session.Save(r, w)
		w.Write([]byte("<h1>Hello Sessions!</h1>"))
	})
	http.HandleFunc("/sessioncontents", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		songs := session.Values["favorite-song"].([]string)
		w.Header().Add("Content-Type", "text/html")
		for _, song := range songs {
			w.Write([]byte(song + "<br />"))
		}
	})
	//seems there are also some other places leaking memory
	http.ListenAndServe(":3000", context.ClearHandler(http.DefaultServeMux))
}

func flashMessages() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		session.Values["favorite-song"] = []string{"My Heart will go on", "Love"}
		session.AddFlash("Flash messages only work once!")
		session.Save(r, w)
		w.Write([]byte("<h1>Hello Sessions!</h1>"))
	})
	http.HandleFunc("/sessioncontents", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		songs := session.Values["favorite-song"].([]string)
		w.Header().Add("Content-Type", "text/html")
		if flashes := session.Flashes(); len(flashes) > 0 {
			session.Save(r, w)
			for _, msg := range flashes {
				m := msg.(string)
				w.Write([]byte(m + "<br />"))
			}
		}
		for _, song := range songs {
			w.Write([]byte(song + "<br />"))
		}
	})
	//seems there are also some other places leaking memory
	http.ListenAndServe(":3000", context.ClearHandler(http.DefaultServeMux))
}
