package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/zeddy-go/zeddy/app"
	"github.com/zeddy-go/zeddy/container"
)

func WithPrefix(prefix string) func(*Module) {
	return func(module *Module) {
		module.prefix = prefix
	}
}

func NewModule(opts ...func(*Module)) *Module {
	m := &Module{
		prefix: "redis",
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

type Module struct {
	app.Module
	prefix string
}

func (m Module) Init() (err error) {
	err = container.Bind[*redis.Client](func(c *viper.Viper) *redis.Client {
		c = c.Sub(m.prefix)
		return redis.NewClient(&redis.Options{
			Addr:     c.GetString("addr"),
			Password: c.GetString("password"),
			DB:       c.GetInt("db"),
		})
	})
	if err != nil {
		return
	}

	return
}
