package gormx

import (
	"fmt"
	"github.com/zeddy-go/zeddy/convert"
	"reflect"
	"strings"

	"github.com/zeddy-go/zeddy/errx"
	"gorm.io/gorm"
)

func contains(target []string, str string) bool {
	for _, s := range target {
		if strings.Contains(str, s) {
			return true
		}
	}

	return false
}

func quote(field string) string {
	var arr []string
	for _, item := range strings.Split(field, ".") {
		arr = append(arr, "`"+item+"`")
	}
	return strings.Join(arr, ".")
}

func applyCondition(db *gorm.DB, conditions ...[]any) (newDB *gorm.DB, err error) {
	newDB = db
	for _, c := range conditions {
		if len(c) < 2 {
			return db, errx.New("condition require at least 2 params")
		}

		list := []string{
			" and ",
			" or ",
			"?",
			" not ",
			" between ",
			" like ",
			" is ",
		}
		if s, ok := c[0].(string); ok && contains(list, strings.ToLower(s)) {
			newDB = newDB.Where(s, c[1:]...)
		} else {
			switch len(c) {
			case 2:
				v := reflect.ValueOf(c[1])
				if v.Kind() == reflect.Slice {
					newDB = newDB.Where(fmt.Sprintf("%s IN ?", quote(c[0].(string))), c[1])
				} else {
					newDB = newDB.Where(fmt.Sprintf("%s = ?", quote(c[0].(string))), c[1])
				}
			case 3:
				switch c[1] {
				case "like":
					var value string
					value, err = convert.To[string](c[2])
					if err != nil {
						return
					}
					newDB = db.Where(fmt.Sprintf("%s LIKE (?)", quote(c[0].(string))), "%"+value+"%")
				default:
					newDB = newDB.Where(fmt.Sprintf("%s %s (?)", quote(c[0].(string)), c[1]), c[2])
				}
			default:
				err = errx.New("condition params is too many")
				return
			}
		}
	}

	return
}
