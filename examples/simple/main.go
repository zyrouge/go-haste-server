package main

import (
	"net/http"

	haste_router "zyrouge.me/haste/router"
)

func HandlePingGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	router := haste_router.NewHasteRouter()
	pingRouter := haste_router.NewHasteRouter()
	pingRouter.HandleGet(http.HandlerFunc(HandlePingGet))
	router.HandleNamed("ping", pingRouter)
	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		panic(err)
	}
}
