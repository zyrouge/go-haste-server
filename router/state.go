package haste_router

import (
	"net/http"
)

type HasteRouterState struct {
	PathTraveller *HasteRouterPathTraveller
	Params        map[string]string
}

func NewHasteRouterState(r *http.Request) *HasteRouterState {
	pathTraveller := NewHasteRouterPathTraveller(r.URL.Path)
	return &HasteRouterState{
		PathTraveller: pathTraveller,
		Params:        map[string]string{},
	}
}

func GetHasteRouterStateFromRequest(r *http.Request) *HasteRouterState {
	return r.Context().Value(HasteRouterStateContextKey).(*HasteRouterState)
}
