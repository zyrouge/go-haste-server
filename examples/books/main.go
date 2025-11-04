package main

import (
	"net/http"

	haste_router "zyrouge.me/haste/router"
)

func main() {
	router := haste_router.NewHasteRouter()
	router.Use(RequestLogMiddleware)
	router.HandleNamed("books", NewBooksRouter())
	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		panic(err)
	}
}
