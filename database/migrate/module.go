package migrate

import (
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
	err = container.Bind[database.Migrator](NewDefaultMigrator)
	if err != nil {
		return
	}

	return
}

func (m Module) Boot() (err error) {
	err = container.Invoke(func(m database.Migrator) (err error) {
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
