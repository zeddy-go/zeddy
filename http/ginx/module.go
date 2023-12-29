package ginx

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/contract"
	"github.com/zeddy-go/zeddy/module"
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
		BaseModule: module.NewBaseModule("ginx"),
		addr:       ":80",
		subModules: make([]contract.IModule, 0),
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
	*module.BaseModule
	router     gin.IRouter
	addr       string
	lts        bool
	certFile   string
	keyFile    string
	subModules []contract.IModule
}

func (m *Module) Register(subs ...contract.IModule) {
	for _, sub := range subs {
		m.subModules = append(m.subModules, sub)

		registerMethod := reflect.ValueOf(sub).MethodByName("RegisterRoute")

		if registerMethod.IsValid() && !registerMethod.IsNil() {
			container.Default().Invoke(registerMethod)
		}
	}
}

func (m *Module) Init() (err error) {
	err = container.Bind[contract.IRouter](m, container.AsSingleton())
	if err != nil {
		return
	}

	return
}

func (m *Module) Any(route string, handler any, middlewares ...any) contract.IRouter {
	m.router.Any(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) GET(route string, handler any, middlewares ...any) contract.IRouter {
	m.router.GET(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) POST(route string, handler any, middlewares ...any) contract.IRouter {
	m.router.POST(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) DELETE(route string, handler any, middlewares ...any) contract.IRouter {
	m.router.DELETE(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) PATCH(route string, handler any, middlewares ...any) contract.IRouter {
	m.router.PATCH(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) PUT(route string, handler any, middlewares ...any) contract.IRouter {
	m.router.PUT(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) OPTIONS(route string, handler any, middlewares ...any) contract.IRouter {
	m.router.OPTIONS(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) HEAD(route string, handler any, middlewares ...any) contract.IRouter {
	m.router.HEAD(route, m.wrap(handler, middlewares...)...)
	return m
}

func (m *Module) Group(prefix string, middlewares ...any) contract.IRouter {
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
