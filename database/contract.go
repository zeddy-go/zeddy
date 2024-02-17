package database

import "database/sql"

type ITransaction interface {
	Begin(sets ...func(*sql.TxOptions))
	Commit() error
	Rollback() error
	Transaction(f func() error) error
}

type IDBHolder[DB any] interface {
	GetDB() DB
	ITransaction
}

type IRepository[Entity any, DB any] interface {
	IDBHolder[DB]
	Create(...*Entity) error
	Update(structOrMap any, conditions ...Condition) error
	First(conditions ...Condition) (Entity, error)
	List(conditions ...Condition) ([]Entity, error)
	Delete(conditions ...Condition) error
	Pagination(offset, limit int, conditions ...Condition) (total int64, list []Entity, err error)
}
