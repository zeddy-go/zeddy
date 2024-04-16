package grpcx

import (
	"github.com/spf13/viper"
	"github.com/zeddy-go/zeddy/app"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/errx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
)

func WithPrefix(prefix string) func(*Module) {
	return func(module *Module) {
		module.prefix = prefix
	}
}

func NewModule(opts ...func(*Module)) *Module {
	m := &Module{
		prefix: "grpc",
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

type Module struct {
	app.IsModule
	grpcServer *grpc.Server
	prefix     string
}

func (m *Module) Init() (err error) {
	c := viper.Sub(m.prefix)

	m.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(simpleInterceptor),
	)

	healthCheck := health.NewServer()
	healthgrpc.RegisterHealthServer(m.grpcServer, healthCheck)

	err = container.Bind[*health.Server](healthCheck)
	if err != nil {
		return
	}

	if c.GetBool("reflection") {
		reflection.Register(m.grpcServer)
	}

	err = container.Bind[*grpc.Server](m.grpcServer)
	if err != nil {
		return
	}

	return
}

func (m *Module) Start() {
	c := viper.Sub(m.prefix)

	var lis net.Listener
	lis, err := net.Listen("tcp", c.GetString("addr"))
	if err != nil {
		panic(errx.Wrap(err, "tcp listen port failed"))
	}

	err = m.grpcServer.Serve(lis)
	if err != nil {
		slog.Error("grpc server shutdown", "error", err)
	}
}

func (m *Module) Stop() {
	m.grpcServer.Stop()
}
