package main

import (
	"log"
	"net/http"
)

func RequestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("request uri: " + r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
