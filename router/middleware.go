package haste_router

import "net/http"

type HasteRouterMiddlewareFunc func(next http.Handler) http.Handler
