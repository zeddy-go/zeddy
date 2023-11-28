package migrate

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source"
)

type IMigrator interface {
	SetSourceURL(url string) IMigrator
	SetSourceInstance(instance source.Driver) IMigrator
	GetSourceInstance() source.Driver
	SetDatabaseURL(dsn string) IMigrator
	Migrate() error
}

type DefaultMigrator struct {
	SourceUrl      string
	SourceInstance source.Driver
	DatabaseUrl    string
}

func (d *DefaultMigrator) SetSourceURL(url string) IMigrator {
	d.SourceUrl = url
	return d
}

func (d *DefaultMigrator) SetSourceInstance(instance source.Driver) IMigrator {
	d.SourceInstance = instance
	return d
}

func (d *DefaultMigrator) GetSourceInstance() source.Driver {
	return d.SourceInstance
}

func (d *DefaultMigrator) SetDatabaseURL(dsn string) IMigrator {
	d.DatabaseUrl = dsn
	return d
}

func (d DefaultMigrator) Migrate() (err error) {
	var m *migrate.Migrate

	if d.SourceInstance != nil {
		m, err = migrate.NewWithSourceInstance("", d.SourceInstance, d.DatabaseUrl)
	} else {
		m, err = migrate.New(d.SourceUrl, d.DatabaseUrl)
	}
	if err != nil {
		return
	}
	defer func() {
		_, _ = m.Close()
	}()

	err = m.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}
	return
}
