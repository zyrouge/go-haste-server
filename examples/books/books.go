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

type BooksMethodHandler struct{}

func NewBooksRouter() http.Handler {
	return haste_router.NewHasteRouter(&BooksMethodHandler{})
}

func (*BooksMethodHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(books)
	w.Write(bytes)
}

func (*BooksMethodHandler) GetParamName() string {
	return "bookId"
}

func (*BooksMethodHandler) HandleGetEntity(w http.ResponseWriter, r *http.Request) {
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
