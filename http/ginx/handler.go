package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/reflectx"
	"io"
	"net/http"
	"reflect"
)

func GinMiddleware(f any) gin.HandlerFunc {
	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		panic(errors.New("func only"))
	}

	if fType.NumOut() > 1 {
		panic(errors.New("middleware should return results not more than 1"))
	}

	return func(ctx *gin.Context) {
		params, err := buildParams(fType, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &Response{
				Message: err.Error(),
			})
			return
		}

		results := reflect.ValueOf(f).Call(params)

		checkResult(ctx, results)
	}
}

func GinHandler(f any) gin.HandlerFunc {
	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		panic(errors.New("func only"))
	}

	if fType.NumOut() > 3 {
		panic(errors.New("middleware should return results not more than 3"))
	} else if fType.NumOut() == 3 {
		if !isNumber(fType.Out(0)) && fType.Out(0).Name() != "IMeta" {
			panic(errors.New("first result should be number(total of records) or IMeta"))
		}
	}

	return func(ctx *gin.Context) {
		params, err := buildParams(fType, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &Response{
				Message: err.Error(),
			})
			return
		}

		results := reflect.ValueOf(f).Call(params)

		response(ctx, results)
	}
}

func buildParams(fType reflect.Type, ctx *gin.Context) (params []reflect.Value, err error) {
	params = make([]reflect.Value, fType.NumIn())
	for i := 0; i < fType.NumIn(); i++ {
		if fType.In(i) == reflect.TypeOf(&gin.Context{}) {
			params[i] = reflect.ValueOf(ctx)
			continue
		}

		//TODO: 优化
		var p reflect.Value
		p, err = container.ResolveType(fType.In(i))
		if err != nil || !p.IsValid() {
			p := reflect.New(fType.In(i))
			if len(ctx.Request.URL.Query()) > 0 {
				err = ctx.ShouldBindQuery(p.Interface())
				if err != nil {
					return
				}
			}
			err = ctx.ShouldBindJSON(p.Interface())
			if err != nil {
				if err == io.EOF {
					err = nil
				} else {
					return
				}
			}
			params[i] = p.Elem()
		} else {
			params[i] = p
		}
	}

	return
}

func checkResult(ctx *gin.Context, results []reflect.Value) {
	switch len(results) {
	case 0:
		fallthrough
	case 1:
		r := results[0].Interface()
		if r == nil {
			return
		}
		switch x := r.(type) {
		case error:
			ErrorAbort(ctx, x, 0)
		}
		return
	default:
		panic(errors.New("middleware should return results not more than 1"))
	}
}

func response(ctx *gin.Context, results []reflect.Value) {
	switch len(results) {
	case 0:
		Success(ctx, nil, 0)
		return
	case 1:
		r := results[0].Interface()
		if r == nil {
			ctx.JSON(http.StatusNoContent, nil)
			return
		}
		switch x := r.(type) {
		case error:
			Error(ctx, x, 0)
		default:
			Success(ctx, r, 0)
		}
		return
	case 2:
		if results[1].IsValid() && !results[1].IsNil() {
			Error(ctx, results[1].Interface().(error), 0)
			return
		}

		Success(ctx, results[0].Interface(), 0)
		return
	case 3:
		if results[2].IsValid() && !results[2].IsNil() {
			Error(ctx, results[2].Interface().(error), 0)
			return
		}

		if isNumber(results[0]) {
			tmp, err := reflectx.ConvertToKind(results[0], reflect.Int)
			if err != nil {
				panic(err)
			}
			Pagination(ctx, results[1].Interface(), tmp.Interface().(int))
		} else if m, ok := results[0].Interface().(IMeta); ok {
			resp := &Response{
				Data: results[1].Interface(),
				Meta: m.GetMeta(),
			}

			Json(ctx, resp, http.StatusOK, false)
		} else {
			panic(errors.New("three results only for pagination"))
		}
	}
}

func isNumber(v interface{ Kind() reflect.Kind }) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}
