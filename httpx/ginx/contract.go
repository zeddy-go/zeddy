package ginx

type Router interface {
	Any(route string, handler any, middlewares ...any) Router
	GET(route string, handler any, middlewares ...any) Router
	POST(route string, handler any, middlewares ...any) Router
	DELETE(route string, handler any, middlewares ...any) Router
	PATCH(route string, handler any, middlewares ...any) Router
	PUT(route string, handler any, middlewares ...any) Router
	OPTIONS(route string, handler any, middlewares ...any) Router
	HEAD(route string, handler any, middlewares ...any) Router
	Group(prefix string, middlewares ...any) Router
	Use(middlewares ...any) Router
}
