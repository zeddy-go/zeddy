package wgorm

import (
	"fmt"
	"github.com/zeddy-go/zeddy/slicex"
	"reflect"
	"strings"

	"github.com/zeddy-go/zeddy/errx"
	"gorm.io/gorm"
)

type Condition []any

func (c Condition) Apply(db *gorm.DB) (newDB *gorm.DB, err error) {
	if len(c) < 2 {
		return db, errx.New("condition require at least 2 params")
	}
	newDB = db

	list := []string{
		" and ",
		" or ",
		"?",
		" not ",
		" between ",
		" like ",
		" is ",
	}
	if s, ok := c[0].(string); ok && slicex.Contains(strings.ToLower(s), list) {
		newDB = newDB.Where(s, c[1:]...)
	} else {
		switch len(c) {
		case 2:
			v := reflect.ValueOf(c[1])
			if v.Kind() == reflect.Slice {
				newDB = newDB.Where(fmt.Sprintf("%s IN ?", c[0]), c[1])
			} else {
				newDB = newDB.Where(fmt.Sprintf("%s = ?", c[0]), c[1])
			}
		case 3:
			newDB = newDB.Where(fmt.Sprintf("%s %s (?)", c[0], c[1]), c[2])
		default:
			err = errx.New("condition params is too many")
			return
		}
	}

	return
}

func sliceAnyToConditions(conditions ...[]any) (r Conditions) {
	r = make(Conditions, len(conditions))
	for i := range conditions {
		r[i] = conditions[i]
	}
	return
}

type Conditions []Condition

func (cs Conditions) Apply(db *gorm.DB) (newDB *gorm.DB, err error) {
	newDB = db
	for _, c := range cs {
		newDB, err = Condition(c).Apply(newDB)
		if err != nil {
			return
		}
	}
	return
}
