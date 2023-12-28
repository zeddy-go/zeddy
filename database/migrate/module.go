package migrate

import (
	"github.com/spf13/viper"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/database"
	"github.com/zeddy-go/zeddy/module"
)

func NewModule() *Module {
	m := &Module{
		BaseModule: module.NewBaseModule("migrate"),
	}

	driver := NewFsDriver()
	container.Bind[*EmbedDriver](driver, container.AsSingleton())
	container.Bind[IMigrator](func(conf *viper.Viper) IMigrator {
		return &DefaultMigrator{
			DatabaseUrl:    database.DSN(conf.GetString("database.dsn")).Encode(),
			SourceInstance: driver,
		}
	}, container.AsSingleton())

	return m
}

type Module struct {
	*module.BaseModule
}

func (m Module) Boot() {
	container.Invoke(func(m IMigrator) {
		err := m.Migrate()
		if err != nil {
			panic(err)
		}
	})
}
