package haste_router

import (
	"context"
	"net/http"
)

type HasteRouter struct {
	Handler       any
	NamedHandlers map[string]http.Handler
}

type HasteRouterPostHandler interface {
	HandlePost(w http.ResponseWriter, r *http.Request)
}

type HasteRouterGetHandler interface {
	HandleGet(w http.ResponseWriter, r *http.Request)
}

type HasteRouterDeleteHandler interface {
	HandleDelete(w http.ResponseWriter, r *http.Request)
}

type HasteRouterParamMethod interface {
	GetParamName() string
}

type HasteRouterGetEntityHandler interface {
	HasteRouterParamMethod
	HandleGetEntity(w http.ResponseWriter, r *http.Request)
}

type HasteRouterPostEntityHandler interface {
	HasteRouterParamMethod
	HandlePostEntity(w http.ResponseWriter, r *http.Request)
}

type HasteRouterDeleteEntityHandler interface {
	HasteRouterParamMethod
	HandleDeleteEntity(w http.ResponseWriter, r *http.Request)
}

type HasteRouterDynamicHandler interface {
	HasteRouterParamMethod
	HandleDynamic(w http.ResponseWriter, r *http.Request)
}

type HasteRouterFallbackHandler interface {
	HandleFallback(w http.ResponseWriter, r *http.Request)
}

func NewHasteRouter(handler any) *HasteRouter {
	router := &HasteRouter{
		Handler:       handler,
		NamedHandlers: map[string]http.Handler{},
	}
	return router
}

func (router *HasteRouter) SetNamedHandler(name string, handler http.Handler) {
	router.NamedHandlers[name] = handler
}

func (router *HasteRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state, ok := r.Context().Value(HasteRouterStateContextKey).(*HasteRouterState)
	if !ok {
		state = NewHasteRouterState(r)
		r = r.WithContext(context.WithValue(r.Context(), HasteRouterStateContextKey, state))
	}
	if !state.PathTraveller.HasNext() {
		switch r.Method {
		case http.MethodGet:
			if handler, ok := router.Handler.(HasteRouterGetHandler); ok {
				handler.HandleGet(w, r)
				return
			}

		case http.MethodPost:
			if handler, ok := router.Handler.(HasteRouterPostHandler); ok {
				handler.HandlePost(w, r)
				return
			}

		case http.MethodDelete:
			if handler, ok := router.Handler.(HasteRouterDeleteHandler); ok {
				handler.HandleDelete(w, r)
				return
			}
		}
	}
	if state.PathTraveller.Remaining() == 1 {
		next, _ := state.PathTraveller.Peek()
		switch r.Method {
		case http.MethodGet:
			if handler, ok := router.Handler.(HasteRouterGetEntityHandler); ok {
				state.PathTraveller.Next()
				paramName := handler.GetParamName()
				state.Params[paramName] = next
				handler.HandleGetEntity(w, r)
				return
			}

		case http.MethodPut:
			if handler, ok := router.Handler.(HasteRouterPostEntityHandler); ok {
				state.PathTraveller.Next()
				paramName := handler.GetParamName()
				state.Params[paramName] = next
				handler.HandlePostEntity(w, r)
				return
			}

		case http.MethodDelete:
			if handler, ok := router.Handler.(HasteRouterDeleteEntityHandler); ok {
				state.PathTraveller.Next()
				paramName := handler.GetParamName()
				state.Params[paramName] = next
				handler.HandleDeleteEntity(w, r)
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
		if handler, ok := router.Handler.(HasteRouterDynamicHandler); ok {
			state.PathTraveller.Next()
			paramName := handler.GetParamName()
			state.Params[paramName] = next
			handler.HandleDynamic(w, r)
			return
		}
	}
	if handler, ok := router.Handler.(HasteRouterFallbackHandler); ok {
		handler.HandleFallback(w, r)
		return
	}
	// TODO
	w.WriteHeader(404)
	w.Write([]byte(http.StatusText(404)))
}
