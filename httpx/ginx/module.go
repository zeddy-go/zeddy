package ginx

import (
	"github.com/spf13/viper"
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

func NewModule(opts ...func(*Module)) *Module {
	m := &Module{
		subModules: make([]app.Module, 0),
	}

	for _, set := range opts {
		set(m)
	}

	if m.router == nil {
		m.router = gin.Default()
	}

	return m
}

type Module struct {
	app.IsModule
	prefix     string
	router     gin.IRouter
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
	return
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
	var c *viper.Viper
	if m.prefix != "" {
		c = viper.Sub(m.prefix)
	} else {
		c = viper.GetViper()
	}

	svr := http.Server{
		Handler: m.router.(http.Handler),
	}
	svr.Addr = c.GetString("addr")

	var err error
	if c.GetBool("lts") {
		err = svr.ListenAndServeTLS(c.GetString("certFile"), c.GetString("keyFile"))
	} else {
		err = svr.ListenAndServe()
	}
	if err != nil {
		slog.Info("server shutdown", "error", err)
	}
}
