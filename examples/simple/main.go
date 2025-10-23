package main

import (
	"net/http"

	haste_router "zyrouge.me/haste/router"
)

type PingHandler struct{}

func (*PingHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	router := haste_router.NewHasteRouter(nil)
	router.SetNamedHandler("ping", haste_router.NewHasteRouter(&PingHandler{}))
	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		panic(err)
	}
}
