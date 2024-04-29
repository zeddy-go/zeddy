package database

import (
	"database/sql"
)

type Migrator interface {
	RegisterMigrates(...any) error
	Migrate() error
	Up(stepNum int) error
	Down(stepNum int) error
}

type ITransaction interface {
	Begin(sets ...func(*sql.TxOptions))
	Commit() error
	Rollback() error
	Transaction(f func() error, sets ...func(*sql.TxOptions)) error
}

type IDBHolder[DB any] interface {
	GetDB() DB
	ITransaction
}

type IRepository[Entity any] interface {
	Create(...*Entity) error
	Update(structOrMap any, conditions ...any) error
	First(conditions ...any) (*Entity, error)
	List(conditions ...any) ([]*Entity, error)
	Delete(conditions ...any) error
	Pagination(offset, limit int, conditions ...any) (total int64, list []*Entity, err error)
	ListInBatch(batchSize int, callback func(repo IRepository[Entity], list []*Entity) error) (err error)
}

type ConditionApplier[DB any] interface {
	Apply(DB) (DB, error)
}

type Order map[string]string
