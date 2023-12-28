package wgorm

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

func NewDBHolder(db *gorm.DB) *DBHolder {
	return &DBHolder{
		root: db,
		txs:  make(map[int64]*gorm.DB),
	}
}

type DBHolder struct {
	root *gorm.DB
	txs  map[int64]*gorm.DB
	lock sync.Mutex
}

func (d *DBHolder) Begin(sets ...func(*sql.TxOptions)) {
	d.lock.Lock()
	defer d.lock.Unlock()

	db := d.pull()
	if db == nil {
		db = d.root.Session(&gorm.Session{
			SkipDefaultTransaction: true,
			Logger:                 d.root.Logger,
		})
		opts := &sql.TxOptions{}
		for _, set := range sets {
			set(opts)
		}
		db = db.Begin(opts)
		d.put(db)
	}
}

func (d *DBHolder) Commit() {
	d.lock.Lock()
	defer d.lock.Unlock()

	w := d.pull()
	w.Commit()
}

func (d *DBHolder) Rollback() {
	d.lock.Lock()
	defer d.lock.Unlock()

	w := d.pull()
	w.Rollback()
}

func (d *DBHolder) put(db *gorm.DB) {
	d.txs[routine.Goid()] = db
}

func (d *DBHolder) pull() *gorm.DB {
	defer delete(d.txs, routine.Goid())
	return d.txs[routine.Goid()]
}

func (d *DBHolder) GetDB() (db *gorm.DB) {
	d.lock.Lock()
	defer d.lock.Unlock()

	db = d.pull()
	if db == nil {
		db = d.root.Session(&gorm.Session{SkipDefaultTransaction: true, Logger: d.root.Logger})
	}

	return
}
