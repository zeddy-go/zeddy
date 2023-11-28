package wgorm

import (
	"github.com/sony/sonyflake"
	"github.com/spf13/viper"
	"github.com/zeddy-go/database"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/module"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewModule() *Module {
	m := &Module{
		BaseModule: module.NewBaseModule("wgorm"),
	}

	container.Register(func(c *viper.Viper) (db *gorm.DB) {
		dsn := database.DSN(c.GetString("database.dsn"))
		db, err := gorm.Open(mysql.Open(dsn.RemoveSchema()))
		if err != nil {
			panic(err)
		}
		return
	}, container.WithSingleton())

	container.Register(NewDBHolder, container.WithSingleton())

	container.Register(func() *sonyflake.Sonyflake {
		return sonyflake.NewSonyflake(sonyflake.Settings{})
	}, container.WithSingleton())

	return m
}

type Module struct {
	*module.BaseModule
}
