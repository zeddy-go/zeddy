package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stoewer/go-strcase"
	"github.com/zeddy-go/zeddy/stringx"
	"gorm.io/gorm"
	"strings"
)

func NewPageFromCtx(ctx *gin.Context, defaultSize int) *Page {
	m := &Page{}
	ctx.ShouldBindQuery(m)
	if m.Page == 0 {
		m.Page = 1
	}
	if m.Size == 0 {
		m.Size = defaultSize
	}
	return m
}

type Page struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func (p Page) Apply(db *gorm.DB) (newDB *gorm.DB, err error) {
	newDB = db.Offset(p.Offset()).Limit(p.Size)
	return
}

func (p Page) Offset() int {
	return (p.Page - 1) * p.Size
}

var operaters = [...]string{
	">=", //大于等于
	"<=", //小于等于
	"!",  //不等于
	"~",  //like
	">",  //大于
	"<",  //小于
}

func NewFiltersFromCtx(ctx *gin.Context) *Filters {
	tmp := ctx.QueryMap("filters")
	for k, v := range tmp {
		delete(tmp, k)
		if v != "" {
			tmp[strcase.SnakeCase(k)] = v
		}
	}
	return &Filters{
		m: tmp,
	}
}

type CustomerParser struct {
	key   string
	parse func(key, value string) any
}

type Filters struct {
	m map[string]string
}

func (f Filters) ParseAll(customerParsers ...CustomerParser) (results []any) {
	results = make([]any, 0, len(f.m))
BIGLOOP:
	for key, value := range f.m {
		for _, parser := range customerParsers {
			if parser.key == key {
				results = append(results, parser.parse(key, value))
				continue BIGLOOP
			}
		}
		results = append(results, parse(key, value))
	}
	return
}

func parse(key, value string) (result []any) {
	result = append(result, key)
	for _, operater := range operaters {
		if v, found := strings.CutPrefix(value, operater); found {
			switch operater {
			case "!":
				result = append(result, "!=", v)
			case "~":
				result = append(result, "like", "%"+v+"%")
			case ">":
				result = append(result, ">", v)
			case "<":
				result = append(result, "<", v)
			case ">=":
				result = append(result, ">=", v)
			case "<=":
				result = append(result, "<=", v)
			}
			return
		}
	}
	result = append(result, value)
	return
}

func NewSortsFromCtx(ctx *gin.Context) *Sorts {
	tmp := ctx.QueryMap("sorts")
	for k, v := range tmp {
		delete(tmp, k)
		if v != "" && stringx.Contains([]string{"desc", "DESC", "asc", "ASC"}, v) {
			tmp[strcase.SnakeCase(k)] = v
		}
	}
	return &Sorts{
		m: tmp,
	}
}

type Sorts struct {
	m map[string]string
}

func (s Sorts) Apply(db *gorm.DB) (newDB *gorm.DB, err error) {
	newDB = db
	for key, value := range s.m {
		newDB = newDB.Order(fmt.Sprintf("%s %s", key, strings.ToUpper(value)))
	}

	return
}
