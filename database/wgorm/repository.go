package wgorm

import (
	"errors"
	"github.com/zeddy-go/zeddy/database"
	"github.com/zeddy-go/zeddy/mapper"
	"gorm.io/gorm"
)

func WithM2E[PO any, Entity any](f func(dst *Entity, src *PO)) func(*Repository[PO, Entity]) {
	return func(r *Repository[PO, Entity]) {
		r.m2e = f
	}
}

func WithE2M[PO any, Entity any](f func(dst *PO, src *Entity)) func(*Repository[PO, Entity]) {
	return func(r *Repository[PO, Entity]) {
		r.e2m = f
	}
}

func defaultM2E[PO any, Entity any](dst *Entity, src *PO) {
	mapper.MustSimpleMap(dst, src)
}

func defaultE2M[PO any, Entity any](dst *PO, src *Entity) {
	mapper.MustSimpleMap(dst, src)
}

func NewRepository[PO any, Entity any](db *gorm.DB, opts ...func(*Repository[PO, Entity])) *Repository[PO, Entity] {
	r := &Repository[PO, Entity]{
		IDBHolder: NewDBHolder(db),
		m2e:       defaultM2E[PO, Entity],
		e2m:       defaultE2M[PO, Entity],
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

type Repository[PO any, Entity any] struct {
	database.IDBHolder[*gorm.DB]
	m2e func(dst *Entity, src *PO)
	e2m func(dst *PO, src *Entity)
}

func (r *Repository[PO, Entity]) Create(entities ...*Entity) (err error) {
	pos := make([]*PO, 0, len(entities))
	for _, item := range entities {
		po := new(PO)
		r.e2m(po, item)
		pos = append(pos, po)
	}

	err = r.GetDB().Create(&pos).Error
	if err != nil {
		return
	}

	for index, item := range pos {
		r.m2e(entities[index], item)
	}
	return
}

// Update struct or map
func (r *Repository[PO, Entity]) Update(entity any, conditions ...database.Condition) (err error) {
	switch x := entity.(type) {
	case *Entity:
		po := new(PO)
		r.e2m(po, x)
		err = r.GetDB().Updates(po).Error
		if err != nil {
			return
		}
		r.m2e(x, po)
	case map[string]any:
		query := r.GetDB()
		if len(conditions) > 0 {
			query, err = applyConditions(query, conditions)
			if err != nil {
				return
			}
		}
		err = query.Model(new(PO)).Updates(entity).Error
	default:
		err = errors.New("only supported struct or map")
	}

	return
}

func (r *Repository[PO, Entity]) Delete(conditions ...database.Condition) (err error) {
	db, err := applyConditions(r.GetDB(), conditions)
	if err != nil {
		return
	}

	err = db.Delete(new(PO)).Error
	return
}

func (r *Repository[PO, Entity]) First(conditions ...database.Condition) (entity Entity, err error) {
	db, err := applyConditions(r.GetDB(), conditions)
	if err != nil {
		return
	}

	po := new(PO)
	err = db.First(po).Error
	if err != nil {
		return
	}

	r.m2e(&entity, po)
	return
}

func (r *Repository[PO, Entity]) List(conditions ...database.Condition) (list []Entity, err error) {
	db, err := applyConditions(r.GetDB(), conditions)
	if err != nil {
		return
	}

	var poList []PO
	err = db.Find(&poList).Error
	if err != nil {
		return
	}

	list = make([]Entity, 0, len(poList))
	for _, item := range poList {
		var dst Entity
		r.m2e(&dst, &item)
		list = append(list, dst)
	}

	return
}

func (r *Repository[PO, Entity]) Pagination(offset, limit int, conditions ...database.Condition) (total int64, list []Entity, err error) {
	db, err := applyConditions(r.GetDB(), conditions)
	if err != nil {
		return
	}
	err = db.Model(new(PO)).Count(&total).Error
	if err != nil {
		return
	}

	var poList []PO
	err = db.Offset(offset).Limit(limit).Find(&poList).Error
	if err != nil {
		return
	}

	list = make([]Entity, 0, len(poList))
	for _, item := range poList {
		var dst Entity
		r.m2e(&dst, &item)
		list = append(list, dst)
	}

	return
}

func (r *Repository[PO, Entity]) ListInBatch(batchSize int, callback func(repo database.IRepository[Entity, *gorm.DB], list []Entity) error) (err error) {
	return r.TransactionTx(func(tx *gorm.DB) (err error) {
		var list []PO
		return tx.FindInBatches(&list, batchSize, func(tx *gorm.DB, batch int) (err error) {
			entities := make([]Entity, 0, len(list))
			for _, item := range list {
				var dst Entity
				r.m2e(&dst, &item)
				entities = append(entities, dst)
			}
			repo := NewRepository[PO, Entity](tx, WithE2M(r.e2m), WithM2E(r.m2e))
			err = callback(repo, entities)
			return
		}).Error
	})
}
