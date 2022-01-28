package main

import (
	"log"
	"net/http"
)

func genRootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		w.Write([]byte("Helloworld!"))
	})
}

func main() {
	log.Println("Hello!")
	http.Handle("/", genRootHandler())

	http.ListenAndServe("0.0.0.0:6673", nil)
}
