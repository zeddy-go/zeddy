package gormx

import (
	"database/sql"
	"sync"

	"github.com/timandy/routine"
	"gorm.io/gorm"
)

func WithIsolation(level sql.IsolationLevel) func(*sql.TxOptions) {
	return func(options *sql.TxOptions) {
		options.Isolation = level
	}
}

func WithReadOnly(readOnly bool) func(*sql.TxOptions) {
	return func(options *sql.TxOptions) {
		options.ReadOnly = readOnly
	}
}

func NewGormDBHolder(db *gorm.DB) *GormDBHolder {
	return &GormDBHolder{
		root: db,
		txs:  make(map[int64]*gorm.DB),
	}
}

type GormDBHolder struct {
	root *gorm.DB
	txs  map[int64]*gorm.DB
	lock sync.Mutex
}

func (d *GormDBHolder) BeginTx(sets ...func(*sql.TxOptions)) (tx *gorm.DB) {
	db := d.root.Session(&gorm.Session{
		SkipDefaultTransaction: true,
		Logger:                 d.root.Logger,
	})
	opts := &sql.TxOptions{}
	for _, set := range sets {
		set(opts)
	}
	return db.Begin(opts)
}

func (d *GormDBHolder) TransactionTx(f func(tx *gorm.DB) error, sets ...func(*sql.TxOptions)) (err error) {
	tx := d.BeginTx(sets...)
	err = f(tx)

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}

	return
}

func (d *GormDBHolder) Transaction(f func() error, sets ...func(*sql.TxOptions)) (err error) {
	d.Begin(sets...)
	err = f()
	if err != nil {
		d.Rollback()
	} else {
		err = d.Commit()
	}

	return
}

func (d *GormDBHolder) Begin(sets ...func(*sql.TxOptions)) {
	d.lock.Lock()
	defer d.lock.Unlock()

	db := d.get()
	if db == nil {
		db = d.BeginTx(sets...)
		d.put(db)
	}
}

func (d *GormDBHolder) Commit() error {
	d.lock.Lock()
	defer func() {
		delete(d.txs, routine.Goid())
		d.lock.Unlock()
	}()

	w := d.get()
	return w.Commit().Error
}

func (d *GormDBHolder) Rollback() error {
	d.lock.Lock()
	defer func() {
		delete(d.txs, routine.Goid())
		d.lock.Unlock()
	}()

	w := d.get()
	return w.Rollback().Error
}

func (d *GormDBHolder) put(db *gorm.DB) {
	d.txs[routine.Goid()] = db
}

func (d *GormDBHolder) get() *gorm.DB {
	return d.txs[routine.Goid()]
}

func (d *GormDBHolder) GetDB() (db *gorm.DB) {
	d.lock.Lock()
	defer d.lock.Unlock()

	db = d.get()
	if db == nil {
		db = d.root.Session(&gorm.Session{SkipDefaultTransaction: true, Logger: d.root.Logger})
	}

	return
}
