package ginx

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/reflectx"
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
		panic(errors.New("should not return results more than 3"))
	} else if fType.NumOut() == 3 {
		if !isNumber(fType.Out(0)) && fType.Out(0).Name() != "IMeta" {
			panic(errors.New("first one of results should be number(total of records) or IMeta"))
		}
	}

	return func(ctx *gin.Context) {
		params, err := buildParams(fType, ctx)
		if err != nil {
			response(reflect.ValueOf(err)).Do(ctx)
			return
		}

		results := reflect.ValueOf(f).Call(params)

		response(results...).Do(ctx)
	}
}

func buildParams(fType reflect.Type, ctx *gin.Context) (params []reflect.Value, err error) {
	params = make([]reflect.Value, fType.NumIn())
	valueCtx := reflect.ValueOf(ctx)
	for i := 0; i < fType.NumIn(); i++ {
		if fType.In(i) == valueCtx.Type() {
			params[i] = valueCtx
			continue
		}

		params[i], err = parseParam(ctx, fType.In(i))
		if err != nil {
			return
		}
	}

	return
}

func parseParam(ctx *gin.Context, t reflect.Type) (p reflect.Value, err error) {
	p, err = container.Default().Resolve(t)

	if err != nil && !errors.Is(err, container.ErrNotFound) {
		return
	}

	if err != nil || !p.IsValid() {
		psrc := newFromType(t)

		if len(ctx.Request.URL.Query()) > 0 {
			_ = ctx.ShouldBindQuery(psrc.Interface())
		}
		if len(ctx.Params) > 0 {
			_ = ctx.ShouldBindUri(psrc.Interface())
		}
		_ = ctx.ShouldBind(psrc.Interface())

		err = binding.Validator.ValidateStruct(psrc.Interface())
		if err != nil {
			return
		}

		pp := reflect.New(t)
		err = reflectx.SetValue(psrc, pp)
		if err != nil {
			return
		}
		p = pp.Elem()
	}

	return
}

func newFromType(t reflect.Type) (r reflect.Value) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	r = reflect.New(t)
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

func response(results ...reflect.Value) (resp IResponse[*gin.Context]) {
	switch len(results) {
	case 0:
		resp = NewJsonResponse(nil, nil, nil)
	case 1:
		r := results[0].Interface()
		if r == nil {
			resp = NewJsonResponse(nil, nil, nil)
			break
		}
		switch x := r.(type) {
		case error:
			resp = NewJsonResponse(nil, nil, x)
		default:
			resp = NewJsonResponse(r, nil, nil)
		}
	case 2:
		if results[1].IsValid() && !results[1].IsNil() {
			resp = NewJsonResponse(nil, nil, results[1].Interface().(error))
			break
		}
		resp = NewJsonResponse(results[0].Interface(), nil, nil)
	case 3:
		if results[2].IsValid() && !results[2].IsNil() {
			resp = NewJsonResponse(nil, nil, results[2].Interface().(error))
			break
		}

		if isNumber(results[0]) {
			tmp, err := reflectx.ConvertToKind(results[0], reflect.Int)
			if err != nil {
				panic(err)
			}
			resp = NewJsonResponse(results[1].Interface(), &Meta{Total: uint(tmp.Interface().(int))}, nil)
		} else if m, ok := results[0].Interface().(IMeta); ok {
			resp = NewJsonResponse(results[1].Interface(), m, nil)
		} else {
			panic(errors.New("three results only for pagination"))
		}
	}

	return
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
