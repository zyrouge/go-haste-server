package haste_router

type HasteRouterContextKey string

var (
	HasteRouterStateContextKey             = HasteRouterContextKey("haste-state")
	HasteRouterMiddlewareInvokedContextKey = HasteRouterContextKey("haste-middleware-invoked")
)
