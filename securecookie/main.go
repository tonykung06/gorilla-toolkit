package main

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

func main() {
	//in production, this should generate once and kept secret
	privateHashKey := securecookie.GenerateRandomKey(64)
	s := securecookie.New(privateHashKey, nil)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if signedEncodedCookie, err := s.Encode("mycookie", "myvalue"); err == nil {
			cookie := http.Cookie{
				Name:  "mycookie",
				Value: signedEncodedCookie, //signed and base64 encoded
				Path:  "/",
			}
			http.SetCookie(w, &cookie)
		}
		w.Write([]byte("Setting cookies"))
	})
	http.ListenAndServe(":3000", nil)
}
