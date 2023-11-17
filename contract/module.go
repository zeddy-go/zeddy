package contract

type IRouter interface {
	Any(route string, handler any, middlewares ...any) IRouter
	GET(route string, handler any, middlewares ...any) IRouter
	POST(route string, handler any, middlewares ...any) IRouter
	DELETE(route string, handler any, middlewares ...any) IRouter
	PATCH(route string, handler any, middlewares ...any) IRouter
	PUT(route string, handler any, middlewares ...any) IRouter
	OPTIONS(route string, handler any, middlewares ...any) IRouter
	HEAD(route string, handler any, middlewares ...any) IRouter
	Group(prefix string, middlewares ...any) IRouter
}

type IHasSubModule interface {
	Register(subs ...IModule)
}

type IShouldBoot interface {
	Boot()
}

type IModule interface {
	Init()
}
