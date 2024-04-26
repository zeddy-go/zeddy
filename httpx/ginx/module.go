package ginx

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zeddy-go/zeddy/app"
	"github.com/zeddy-go/zeddy/container"
	"log/slog"
	"net/http"
	"time"
)

func WithCustomEngine(e *gin.Engine) func(*Module) {
	return func(module *Module) {
		module.router = e
	}
}

func WithPrefix(prefix string) func(*Module) {
	return func(module *Module) {
		module.prefix = prefix
	}
}

func NewModule(opts ...func(*Module)) *Module {
	m := &Module{}

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
	prefix string
	router gin.IRouter
	svr    *http.Server
}

func (m *Module) Init() (err error) {
	var c *viper.Viper
	if m.prefix != "" {
		c = viper.Sub(m.prefix)
	} else {
		c = viper.GetViper()
	}

	if c.GetBool("cors") {
		m.router.Use(CORS)
	}

	err = container.Bind[Router](m)
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

	m.svr = &http.Server{
		Handler: m.router.(http.Handler),
	}
	m.svr.Addr = c.GetString("addr")

	var err error
	if c.GetBool("lts") {
		err = m.svr.ListenAndServeTLS(c.GetString("certFile"), c.GetString("keyFile"))
	} else {
		err = m.svr.ListenAndServe()
	}
	if err != nil {
		slog.Info("server shutdown", "error", err)
	}
}

func (m *Module) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = m.svr.Shutdown(ctx)
}
