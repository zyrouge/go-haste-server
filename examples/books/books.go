package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	haste_router "zyrouge.me/haste/router"
)

type Book struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var books = []*Book{
	{0, "Lord of Mysteries"},
	{1, "Bottom-Tier Character Tomozaki"},
}

func NewBooksRouter() http.Handler {
	router := haste_router.NewHasteRouter()
	router.SetParamName("bookId")
	router.HandleGet(http.HandlerFunc(HandleBooksGet))
	router.HandleGetEntity(http.HandlerFunc(HandleBooksGetEntity))
	return router
}

func HandleBooksGet(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(books)
	w.Write(bytes)
}

func HandleBooksGetEntity(w http.ResponseWriter, r *http.Request) {
	routerState := haste_router.GetHasteRouterStateFromRequest(r)
	bookId, _ := strconv.Atoi(routerState.Params["bookId"])
	for _, x := range books {
		if bookId == x.Id {
			bytes, _ := json.Marshal(x)
			w.Write(bytes)
			return
		}
	}
	w.Write([]byte{})
}
