package gormx

import (
	"errors"
	"fmt"
	"github.com/zeddy-go/zeddy/database"
	"github.com/zeddy-go/zeddy/errx"
	"github.com/zeddy-go/zeddy/mapper"
	"gorm.io/gorm"
)

func WithM2E[PO any, Entity any](f func(dst *Entity, src *PO) error) func(*Repository[PO, Entity]) {
	return func(r *Repository[PO, Entity]) {
		if f == nil {
			r.m2e = defaultM2E[PO, Entity]
		} else {
			r.m2e = f
		}
	}
}

func WithE2M[PO any, Entity any](f func(dst *PO, src *Entity) error) func(*Repository[PO, Entity]) {
	return func(r *Repository[PO, Entity]) {
		if f == nil {
			r.e2m = defaultE2M[PO, Entity]
		} else {
			r.e2m = f
		}
	}
}

func defaultM2E[PO any, Entity any](dst *Entity, src *PO) error {
	return mapper.SimpleMap(dst, src)
}

func defaultE2M[PO any, Entity any](dst *PO, src *Entity) error {
	return mapper.SimpleMap(dst, src)
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
	m2e func(dst *Entity, src *PO) error
	e2m func(dst *PO, src *Entity) error
}

func (r *Repository[PO, Entity]) E2M(dst *PO, src *Entity) (err error) {
	if r.e2m != nil {
		return r.e2m(dst, src)
	} else {
		return defaultE2M(dst, src)
	}
}

func (r *Repository[PO, Entity]) M2E(dst *Entity, src *PO) (err error) {
	if r.m2e != nil {
		return r.m2e(dst, src)
	} else {
		return defaultM2E(dst, src)
	}
}

func (r *Repository[PO, Entity]) Create(entities ...*Entity) (err error) {
	pos := make([]*PO, 0, len(entities))
	for _, item := range entities {
		po := new(PO)
		err = r.E2M(po, item)
		if err != nil {
			return
		}
		pos = append(pos, po)
	}

	err = r.GetDB().Create(&pos).Error
	if err != nil {
		return
	}

	for index, item := range pos {
		err = r.M2E(entities[index], item)
		if err != nil {
			return
		}
	}
	return
}

// Update struct or map
func (r *Repository[PO, Entity]) Update(entity any, conditions ...any) (err error) {
	switch x := entity.(type) {
	case *Entity:
		po := new(PO)
		err = r.E2M(po, x)
		if err != nil {
			return
		}
		err = r.GetDB().Updates(po).Error
		if err != nil {
			return
		}
		err = r.M2E(x, po)
		if err != nil {
			return
		}
	case map[string]any:
		query := r.GetDB()
		if len(conditions) > 0 {
			query, err = apply(query, conditions...)
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

func (r *Repository[PO, Entity]) Delete(conditions ...any) (err error) {
	db, err := apply(r.GetDB(), conditions...)
	if err != nil {
		return
	}

	err = db.Delete(new(PO)).Error
	return
}

func (r *Repository[PO, Entity]) First(conditions ...any) (entity *Entity, err error) {
	db, err := apply(r.GetDB(), conditions...)
	if err != nil {
		return
	}

	po := new(PO)
	err = db.First(po).Error
	if err != nil {
		return
	}

	entity = new(Entity)
	err = r.M2E(entity, po)
	return
}

func (r *Repository[PO, Entity]) List(conditions ...any) (list []*Entity, err error) {
	db, err := apply(r.GetDB(), conditions...)
	if err != nil {
		return
	}

	var poList []PO
	err = db.Find(&poList).Error
	if err != nil {
		return
	}

	list = make([]*Entity, 0, len(poList))
	for _, item := range poList {
		var dst Entity
		err = r.M2E(&dst, &item)
		if err != nil {
			return
		}
		list = append(list, &dst)
	}

	return
}

func (r *Repository[PO, Entity]) Pagination(offset, limit int, conditions ...any) (total int64, list []*Entity, err error) {
	db, err := apply(r.GetDB(), conditions...)
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

	list = make([]*Entity, 0, len(poList))
	for _, item := range poList {
		var dst Entity
		err = r.M2E(&dst, &item)
		if err != nil {
			return
		}
		list = append(list, &dst)
	}

	return
}

func (r *Repository[PO, Entity]) ListInBatch(batchSize int, callback func(repo database.IRepository[Entity], list []*Entity) error) (err error) {
	return r.Transaction(func() (err error) {
		var list []PO
		return r.GetDB().FindInBatches(&list, batchSize, func(tx *gorm.DB, batch int) (err error) {
			entities := make([]*Entity, 0, len(list))
			for _, item := range list {
				var dst Entity
				err = r.M2E(&dst, &item)
				if err != nil {
					return
				}
				entities = append(entities, &dst)
			}
			repo := NewRepository[PO, Entity](tx, WithE2M(r.e2m), WithM2E(r.m2e))
			err = callback(repo, entities)
			return
		}).Error
	})
}

func apply(db *gorm.DB, conditions ...any) (newDB *gorm.DB, err error) {
	if len(conditions) == 0 {
		return db, nil
	}
	newDB = db
	for _, condition := range conditions {
		switch x := condition.(type) {
		case database.ConditionApplier[*gorm.DB]:
			newDB, err = x.Apply(newDB)
			if err != nil {
				return
			}
		case []any:
			newDB, err = applyCondition(newDB, x)
			if err != nil {
				return
			}
		case [][]any:
			newDB, err = applyCondition(newDB, x...)
			if err != nil {
				return
			}
		case database.Order:
			for key, value := range x {
				if key == "" {
					newDB = newDB.Order(value)
				} else {
					newDB = newDB.Order(key + " " + value)
				}
			}
		default:
			err = errx.New(fmt.Sprintf("unsupported condition type: %T", condition))
			return
		}
	}
	return
}
