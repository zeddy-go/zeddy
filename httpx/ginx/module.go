package ginx

import (
	"github.com/zeddy-go/zeddy/app"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/zeddy-go/zeddy/container"
)

func WithCustomEngine(e *gin.Engine) func(*Module) {
	return func(module *Module) {
		module.router = e
	}
}

func WithLts(certFile string, keyFile string) func(*Module) {
	return func(module *Module) {
		module.lts = true
		module.certFile = certFile
		module.keyFile = keyFile
	}
}

func WithAddr(addr string) func(*Module) {
	return func(module *Module) {
		module.addr = addr
	}
}

func NewModule(sets ...func(*Module)) *Module {
	m := &Module{
		addr:       ":8080",
		subModules: make([]app.Module, 0),
	}

	for _, set := range sets {
		set(m)
	}

	if m.router == nil {
		m.router = gin.Default()
	}

	return m
}

type Module struct {
	app.IsModule
	router     gin.IRouter
	addr       string
	lts        bool
	certFile   string
	keyFile    string
	subModules []app.Module
}

func (m *Module) Register(subs ...app.Module) (err error) {
	for _, sub := range subs {
		m.subModules = append(m.subModules, sub)

		registerMethod := reflect.ValueOf(sub).MethodByName("RegisterRoute")

		if registerMethod.IsValid() && !registerMethod.IsNil() {
			_, err = container.Default().Invoke(registerMethod)
			if err != nil {
				return
			}
		}
	}
}

func (m *Module) Init() (err error) {
	err = container.Bind[Router](m, container.AsSingleton())
	if err != nil {
		return
	}

	return
}

func (m *Module) Any(route string, handler any, middlewares ...any) Router {
	m.router.Any(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) GET(route string, handler any, middlewares ...any) Router {
	m.router.GET(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) POST(route string, handler any, middlewares ...any) Router {
	m.router.POST(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) DELETE(route string, handler any, middlewares ...any) Router {
	m.router.DELETE(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) PATCH(route string, handler any, middlewares ...any) Router {
	m.router.PATCH(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) PUT(route string, handler any, middlewares ...any) Router {
	m.router.PUT(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) OPTIONS(route string, handler any, middlewares ...any) Router {
	m.router.OPTIONS(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) HEAD(route string, handler any, middlewares ...any) Router {
	m.router.HEAD(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) Group(prefix string, middlewares ...any) Router {
	group := m.router.Group(prefix, m.wrap(nil, middlewares...)...)
	return &Module{
		router: group,
	}
}

func (m *Module) wrap(handler any, middlewares ...any) (handlers []gin.HandlerFunc) {
	handlers = make([]gin.HandlerFunc, 0, len(middlewares)+1)
	for _, item := range middlewares {
		handlers = append(handlers, GinMiddleware(item))
	}

	if handler != nil {
		handlers = append(handlers, GinHandler(handler))
	}

	return
}

func (m *Module) Start() {
	svr := http.Server{
		Addr:    m.addr,
		Handler: m.router.(http.Handler),
	}

	var err error
	if m.lts {
		err = svr.ListenAndServeTLS(m.certFile, m.keyFile)
	} else {
		err = svr.ListenAndServe()
	}
	if err != nil {
		slog.Info("server shutdown", "error", err)
	}
}
