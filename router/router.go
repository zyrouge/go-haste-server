package haste_router

import (
	"context"
	"net/http"
)

type HasteRouter struct {
	ParamName           string
	Middleware          http.Handler
	GetHandler          http.Handler
	PostHandler         http.Handler
	DeleteHandler       http.Handler
	GetEntityHandler    http.Handler
	PostEntityHandler   http.Handler
	DeleteEntityHandler http.Handler
	DynamicHandler      http.Handler
	FallbackHandler     http.Handler
	NamedHandlers       map[string]http.Handler
}

func NewHasteRouter() *HasteRouter {
	router := &HasteRouter{
		NamedHandlers: map[string]http.Handler{},
	}
	return router
}

func (router *HasteRouter) SetParamName(name string) {
	router.ParamName = name
}

func (router *HasteRouter) Use(middleware HasteRouterMiddlewareFunc) {
	next := router.Middleware
	if next == nil {
		next = router
	}
	router.Middleware = middleware(next)
}

func (router *HasteRouter) HandleGet(handler http.Handler) {
	router.GetHandler = handler
}

func (router *HasteRouter) HandlePost(handler http.Handler) {
	router.PostHandler = handler
}

func (router *HasteRouter) HandleDelete(handler http.Handler) {
	router.DeleteHandler = handler
}

func (router *HasteRouter) HandleGetEntity(handler http.Handler) {
	router.GetEntityHandler = handler
}

func (router *HasteRouter) HandlePostEntity(handler http.Handler) {
	router.PostEntityHandler = handler
}

func (router *HasteRouter) HandleDeleteEntity(handler http.Handler) {
	router.DeleteEntityHandler = handler
}

func (router *HasteRouter) HandleDynamic(handler http.Handler) {
	router.DynamicHandler = handler
}

func (router *HasteRouter) HandleFallback(handler http.Handler) {
	router.FallbackHandler = handler
}

func (router *HasteRouter) HandleNamed(name string, handler http.Handler) {
	router.NamedHandlers[name] = handler
}

func (router *HasteRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if router.Middleware != nil {
		done := ctx.Value(HasteRouterMiddlewareInvokedContextKey)
		if done == nil {
			ctx = context.WithValue(ctx, HasteRouterMiddlewareInvokedContextKey, true)
			r = r.WithContext(ctx)
			router.Middleware.ServeHTTP(w, r)
			return
		}
	}
	state, ok := ctx.Value(HasteRouterStateContextKey).(*HasteRouterState)
	if !ok {
		state = NewHasteRouterState(r)
		ctx = context.WithValue(ctx, HasteRouterStateContextKey, state)
		r = r.WithContext(ctx)
	}
	if !state.PathTraveller.HasNext() {
		switch r.Method {
		case http.MethodGet:
			if router.GetHandler != nil {
				router.GetHandler.ServeHTTP(w, r)
				return
			}

		case http.MethodPost:
			if router.PostHandler != nil {
				router.PostHandler.ServeHTTP(w, r)
				return
			}

		case http.MethodDelete:
			if router.DeleteHandler != nil {
				router.DeleteHandler.ServeHTTP(w, r)
				return
			}
		}
	}
	if state.PathTraveller.Remaining() == 1 {
		next, _ := state.PathTraveller.Peek()
		switch r.Method {
		case http.MethodGet:
			if router.GetEntityHandler != nil {
				state.PathTraveller.Next()
				state.Params[router.ParamName] = next
				router.GetEntityHandler.ServeHTTP(w, r)
				return
			}

		case http.MethodPost:
			if router.PostEntityHandler != nil {
				state.PathTraveller.Next()
				state.Params[router.ParamName] = next
				router.PostEntityHandler.ServeHTTP(w, r)
				return
			}

		case http.MethodDelete:
			if router.DeleteEntityHandler != nil {
				state.PathTraveller.Next()
				state.Params[router.ParamName] = next
				router.DeleteEntityHandler.ServeHTTP(w, r)
				return
			}
		}
	}
	if state.PathTraveller.HasNext() {
		next, _ := state.PathTraveller.Peek()
		namedHandler, ok := router.NamedHandlers[next]
		if ok {
			state.PathTraveller.Next()
			namedHandler.ServeHTTP(w, r)
			return
		}
		if router.DynamicHandler != nil {
			state.PathTraveller.Next()
			state.Params[router.ParamName] = next
			router.DynamicHandler.ServeHTTP(w, r)
			return
		}
	}
	if router.FallbackHandler != nil {
		router.FallbackHandler.ServeHTTP(w, r)
		return
	}
	// TODO
	w.WriteHeader(404)
	w.Write([]byte(http.StatusText(404)))
}
