package migrate

import (
	"github.com/spf13/viper"
	"github.com/zeddy-go/zeddy/app"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/database"
)

func NewModule() *Module {
	m := &Module{}

	return m
}

type Module struct {
	app.IsModule
}

func (m Module) Init() (err error) {
	driver := NewFsDriver()
	err = container.Bind[*EmbedDriver](driver, container.AsSingleton())
	if err != nil {
		return
	}

	err = container.Bind[IMigrator](func(conf *viper.Viper) IMigrator {
		return &DefaultMigrator{
			DatabaseUrl:    database.DSN(conf.GetString("database.dsn")).Encode(),
			SourceInstance: driver,
		}
	}, container.AsSingleton())
	if err != nil {
		return
	}

	return
}

func (m Module) Boot() (err error) {
	err = container.Invoke(func(m IMigrator) (err error) {
		err = m.Migrate()
		if err != nil {
			return
		}
		return
	})
	if err != nil {
		return
	}
	return
}
