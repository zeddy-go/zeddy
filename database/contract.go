package database

import (
	"database/sql"
	"gorm.io/gorm"
)

type ITransaction interface {
	Begin(sets ...func(*sql.TxOptions))
	Commit() error
	Rollback() error
	Transaction(f func() error, sets ...func(*sql.TxOptions)) error
	BeginTx(sets ...func(*sql.TxOptions)) (tx *gorm.DB)
	TransactionTx(f func(tx *gorm.DB) error, sets ...func(*sql.TxOptions)) (err error)
}

type IDBHolder[DB any] interface {
	GetDB() DB
	ITransaction
}

type IRepository[Entity any, DB any] interface {
	IDBHolder[DB]
	Create(...*Entity) error
	Update(structOrMap any, conditions ...any) error
	First(conditions ...any) (*Entity, error)
	List(conditions ...any) ([]*Entity, error)
	Delete(conditions ...any) error
	Pagination(offset, limit int, conditions ...any) (total int64, list []*Entity, err error)
	ListInBatch(batchSize int, callback func(repo IRepository[Entity, DB], list []*Entity) error) (err error)
}

type Condition[DB any] interface {
	Apply(DB) (DB, error)
}
