package wgorm

import (
	"github.com/sony/sonyflake"
	"github.com/spf13/viper"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/database"
	"github.com/zeddy-go/zeddy/module"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewModule() *Module {
	m := &Module{
		BaseModule: module.NewBaseModule("wgorm"),
	}
	return m
}

type Module struct {
	*module.BaseModule
}

func (m *Module) Init() (err error) {
	err = container.Bind[*gorm.DB](func(c *viper.Viper) (db *gorm.DB) {
		dsn := database.DSN(c.GetString("database.dsn"))
		db, err := gorm.Open(mysql.Open(dsn.RemoveSchema()))
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
