package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
)

func main() {
	// authenticatedCookie()
	authenticatedAndEncryptedCookie()
}

func authenticatedCookie() {
	//in production, these keys should generate once and kept secret
	privateHashKey := securecookie.GenerateRandomKey(64) //32 or 64
	cookieCodec := securecookie.New(privateHashKey, nil)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//to base64 decode the signed base64 encoded cookie, on browser, do atob("MTQ3MTE4NTg4OHxDZ3dBQjIxNWRtRnNkV1U9fJpSombuiu7w0MjZ1UOJOj9qqABVsScMkvW4ujoD2o8b".replace("-", "+").replace("_", "/"))
		//we have two base64 encoding schema, the traditional one maps binary data to 64 ASCII characters that are safe to transmit without data loss --> a-z A-Z 0-9 + /
		//but since + and / have special meanings to browser, so they are replaced with - and _ respectively to give a url friendly schema
		//but browsers atob only understands the traditional one and golang is using the url friendly one, we need to do string replace
		if signedEncodedCookie, err := cookieCodec.Encode("mycookie", "myvalue"); err == nil {
			cookie := http.Cookie{
				Name:  "mycookie",
				Value: signedEncodedCookie, //signed and base64 encoded, which is comprised of three parts: a timestamp(for expiry checking) | cookie value | hash signature = f(name, value, private key, timestamp)
				Path:  "/",
			}
			http.SetCookie(w, &cookie)
		}
		w.Write([]byte("Setting signed cookies"))
	})
	http.ListenAndServe(":3000", nil)
}

func authenticatedAndEncryptedCookie() {
	//in production, these keys should generate once and kept secret
	privateHashKey := securecookie.GenerateRandomKey(64)
	blockKey := securecookie.GenerateRandomKey(32) //16, 24 or 32 corresponding to AES128, AES196 and AES256 encryption algorithms respectively

	s := securecookie.New(privateHashKey, blockKey)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if signedEncryptedEncodedCookie, err := s.Encode("mycookie", "myvalue"); err == nil {
			cookie := http.Cookie{
				Name:  "mycookie",
				Value: signedEncryptedEncodedCookie, //signed , encrypted and base64 encoded, which is comprised of three parts: a timestamp(for expiry checking) | cookie value | hash signature = f(name, value, hash private key, timestamp)
				Path:  "/",
			}

			//decoding, decrypting and authenticating the cookie
			var result string
			if err = s.Decode("mycookie", cookie.Value, &result); err == nil {
				fmt.Println(result)
			}

			http.SetCookie(w, &cookie)
		}
		w.Write([]byte("Setting signed and encrypted cookies"))
	})
	http.ListenAndServe(":3000", nil)
}
