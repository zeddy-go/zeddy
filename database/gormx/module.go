//go:build !sqlite

package gormx

import (
	"database/sql"
	"github.com/sony/sonyflake"
	"github.com/spf13/viper"
	"github.com/zeddy-go/zeddy/app"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
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

func (m *Module) newGorm(c *viper.Viper) (*gorm.DB, error) {
	c = c.Sub(m.prefix)
	dsn := database.DSN(c.GetString("dsn"))
	return gorm.Open(mysql.Open(dsn.RemoveSchema()), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
}

func (m *Module) Init() (err error) {
	var (
		gormInstance *gorm.DB
		lock         sync.Mutex
	)
	err = container.Bind[*gorm.DB](func(c *viper.Viper) (db *gorm.DB, err error) {
		lock.Lock()
		defer lock.Unlock()
		if gormInstance != nil {
			db = gormInstance
			return
		}
		gormInstance, err = m.newGorm(c)
		if err != nil {
			return
		}
		db = gormInstance
		return
	})
	if err != nil {
		return
	}

	err = container.Bind[*sql.DB](func(c *viper.Viper) (db *sql.DB, err error) {
		lock.Lock()
		defer lock.Unlock()

		if gormInstance != nil {
			return gormInstance.DB()
		}
		gormInstance, err = m.newGorm(c)
		if err != nil {
			return
		}
		return gormInstance.DB()
	})
	if err != nil {
		return
	}

	err = container.Bind[*GormDBHolder](NewGormDBHolder)
	if err != nil {
		return
	}

	err = container.Bind[*sonyflake.Sonyflake](func() *sonyflake.Sonyflake {
		return sonyflake.NewSonyflake(sonyflake.Settings{})
	})
	if err != nil {
		return
	}

	return
}
