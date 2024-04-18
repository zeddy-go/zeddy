//go:build sqlite

package gormx

import (
	"github.com/sony/sonyflake"
	"github.com/spf13/viper"
	"github.com/zeddy-go/zeddy/app"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/database"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func WithPrefix(prefix string) func(*Module) {
	return func(module *Module) {
		module.prefix = prefix
	}
}

func NewModule(opts ...func(*Module)) *Module {
	m := &Module{
		prefix: "database",
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

type Module struct {
	app.IsModule
	prefix string
}

func (m *Module) Init() (err error) {
	err = container.Bind[*gorm.DB](func(c *viper.Viper) (db *gorm.DB) {
		c = c.Sub(m.prefix)
		dsn := database.DSN(c.GetString("dsn"))
		db, err := gorm.Open(sqlite.Open(dsn.RemoveSchema()), &gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			}),
		})
		if err != nil {
			panic(err)
		}
		return
	}, container.AsSingleton())
	if err != nil {
		return
	}

	err = container.Bind[*DBHolder](NewDBHolder, container.AsSingleton())
	if err != nil {
		return
	}

	err = container.Bind[*sonyflake.Sonyflake](func() *sonyflake.Sonyflake {
		return sonyflake.NewSonyflake(sonyflake.Settings{})
	}, container.AsSingleton())
	if err != nil {
		return
	}

	return
}
