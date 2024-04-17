//go:build !sqlite

package migrate

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"io/fs"
)

func NewDefaultMigrator(db *sql.DB) *DefaultMigrator {
	return &DefaultMigrator{
		db:             db,
		SourceInstance: NewFsDriver(),
	}
}

type DefaultMigrator struct {
	SourceInstance *EmbedDriver
	db             *sql.DB
}

func (d *DefaultMigrator) Up(stepNum int) error {
	//TODO implement me
	panic("implement me")
}

func (d *DefaultMigrator) Down(stepNum int) error {
	//TODO implement me
	panic("implement me")
}

func (d *DefaultMigrator) RegisterMigrates(ms ...any) (err error) {
	for _, m := range ms {
		d.SourceInstance.Add(m.(fs.FS))
	}

	return
}

func (d *DefaultMigrator) Migrate() (err error) {
	var m *migrate.Migrate

	ctx := context.Background()
	newConn, err := d.db.Conn(ctx)
	if err != nil {
		return
	}
	dbInstance, err := mysql.WithConnection(ctx, newConn, &mysql.Config{})
	if err != nil {
		return
	}
	m, err = migrate.NewWithInstance("", d.SourceInstance, "", dbInstance)
	if err != nil {
		return
	}
	defer func() {
		_, _ = m.Close()
	}()

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		err = nil
	}
	return
}
