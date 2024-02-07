package wgorm

import (
	"errors"
	"github.com/zeddy-go/zeddy/database"
	"github.com/zeddy-go/zeddy/mapper"
	"gorm.io/gorm"
)

type Repository[PO any, Entity any] struct {
	database.IDBHolder[*gorm.DB]
}

func (r *Repository[PO, Entity]) Create(entities ...*Entity) (err error) {
	pos := make([]*PO, 0, len(entities))
	for _, item := range entities {
		po := new(PO)
		mapper.MustSimpleMap(po, item)
		pos = append(pos, po)
	}

	err = r.GetDB().Create(&pos).Error
	if err != nil {
		return
	}

	for index, item := range pos {
		mapper.MustSimpleMap(entities[index], item)
	}
	return
}

// Update struct or map
func (r *Repository[PO, Entity]) Update(entity any, conditions ...database.Condition) (err error) {
	switch entity.(type) {
	case *Entity:
		po := new(PO)
		mapper.MustSimpleMap(po, entity)
		err = r.GetDB().Updates(po).Error
		if err != nil {
			return
		}
		mapper.MustSimpleMap(entity, po)
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

	mapper.MustSimpleMap(&entity, po)

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

	mapper.MustSimpleMapSlice(&list, poList)

	return
}

func (r *Repository[PO, Entity]) Pagination(offset, limit int, conditions ...database.Condition) (total int64, list []Entity, err error) {
	db, err := applyConditions(r.GetDB(), conditions)
	if err != nil {
		return
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	var poList []PO
	err = db.Offset(offset).Limit(limit).Find(&poList).Error
	if err != nil {
		return
	}

	mapper.MustSimpleMapSlice(&list, poList)

	return
}
