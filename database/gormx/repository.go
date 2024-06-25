package gormx

import (
	"errors"
	"fmt"
	"github.com/zeddy-go/zeddy/container"
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

func NewRepository[PO any, Entity any](opts ...func(*Repository[PO, Entity])) *Repository[PO, Entity] {
	r := &Repository[PO, Entity]{
		GormDBHolder: container.MustResolve[*GormDBHolder](),
		m2e:          defaultM2E[PO, Entity],
		e2m:          defaultE2M[PO, Entity],
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

type Repository[PO any, Entity any] struct {
	*GormDBHolder
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
			query, err = Apply(query, conditions...)
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
	db, err := Apply(r.GetDB(), conditions...)
	if err != nil {
		return
	}

	err = db.Delete(new(PO)).Error
	return
}

func (r *Repository[PO, Entity]) First(conditions ...any) (entity *Entity, err error) {
	db, err := Apply(r.GetDB(), conditions...)
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
	db, err := Apply(r.GetDB(), conditions...)
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
	db, err := Apply(r.GetDB(), conditions...)
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

func Apply(db *gorm.DB, conditions ...any) (newDB *gorm.DB, err error) {
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
			xx := make([]database.Condition, 0, len(x))
			for _, item := range x {
				xx = append(xx, item)
			}
			newDB, err = applyCondition(newDB, xx...)
			if err != nil {
				return
			}
		case database.Condition:
			newDB, err = applyCondition(newDB, x)
			if err != nil {
				return
			}
		case []database.Condition:
			newDB, err = applyCondition(newDB, x...)
			if err != nil {
				return
			}
		case database.Order:
			for _, item := range x {
				newDB = newDB.Order(item)
			}
		default:
			err = errx.New(fmt.Sprintf("unsupported condition type: %T", condition))
			return
		}
	}
	return
}

type JoinWhere struct {
	Direction  string
	Table      string
	Conditions string
	Where      []any
}

func (j *JoinWhere) Apply(db *gorm.DB) (*gorm.DB, error) {
	return db.
		Joins(fmt.Sprintf("%s JOIN %s ON %s", j.Direction, j.Table, j.Conditions)).
		Where(j.Where[0], j.Where[1:]...), nil
}
