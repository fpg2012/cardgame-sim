package main

import (
	"fmt"
	"log"
	"net/http"
)

func genRootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte("Helloworld!"))
	})
}

func genGetTokenHandler(tm *TokenManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		username := r.URL.Query().Get("username")
		if username == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token, err := tm.AllocateToken(username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprint(token)))
	})
}

func main() {
	log.Println("Hello!")

	tm := NewTokenManager()

	http.Handle("/", genRootHandler())
	http.Handle("/token", genGetTokenHandler(tm))

	http.ListenAndServe("0.0.0.0:6673", nil)
}
