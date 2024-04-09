package grpc

import (
	"fmt"
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

func NewModule(prefix string) app.Module {
	return &Module{
		prefix: prefix,
	}
}

type Module struct {
	app.IsModule
	grpcServer *grpc.Server
	prefix     string
	c          *viper.Viper
}

func (m *Module) Init() (err error) {
	m.c, err = container.Resolve[*viper.Viper]()
	if err != nil {
		return
	}
	if m.prefix != "" {
		m.c = m.c.Sub(m.prefix)
	}

	m.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(simpleInterceptor),
	)

	healthCheck := health.NewServer()
	healthgrpc.RegisterHealthServer(m.grpcServer, healthCheck)

	err = container.Bind[*health.Server](healthCheck, container.AsSingleton())
	if err != nil {
		return
	}

	if m.c.GetBool("reflection") {
		reflection.Register(m.grpcServer)
	}

	err = container.Bind[*grpc.Server](m.grpcServer, container.AsSingleton())
	if err != nil {
		return
	}

	return
}

func (m *Module) Start() {
	var lis net.Listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", m.c.GetInt("port")))
	if err != nil {
		panic(errx.Wrap(err, "tcp listen port failed"))
	}

	err = m.grpcServer.Serve(lis)
	if err != nil {
		slog.Error("grpc server shutdown", "error", err)
	}
}
