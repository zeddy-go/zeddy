package gormx

import "gorm.io/gorm"

func NewSort(field, mode string) *Sort {
	return &Sort{field: field, mode: mode}
}

type Sort struct {
	field string
	mode  string
}

func (s *Sort) Apply(db *gorm.DB) (*gorm.DB, error) {
	return db.Order(s.field + " " + s.mode), nil
}
