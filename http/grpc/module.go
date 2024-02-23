package grpc

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/contract"
	"github.com/zeddy-go/zeddy/errx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
)

func NewModule(prefix string) contract.IModule {
	return &Module{
		prefix: prefix,
	}
}

type Module struct {
	grpcServer *grpc.Server
	prefix     string
	c          *viper.Viper
}

func (m *Module) Name() string {
	return "grpc"
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
