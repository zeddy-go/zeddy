package ginx

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/zeddy-go/zeddy/errx"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/reflectx"
)

type NewResponseFunc func(data any, meta IMeta, err error) IResponse[*gin.Context]

var DefaultNewResponseFunc NewResponseFunc = NewJsonResponse

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
			parseAndResponse(reflect.ValueOf(errx.Wrap(err, "build params failed", errx.WithAbort())))
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
			parseAndResponse(reflect.ValueOf(err)).Do(ctx)
			return
		}

		results := reflect.ValueOf(f).Call(params)

		parseAndResponse(results...).Do(ctx)
	}
}

func buildParams(fType reflect.Type, ctx *gin.Context) (params []reflect.Value, err error) {
	params = make([]reflect.Value, fType.NumIn())
	valueCtx := reflect.ValueOf(ctx)
	for i := 0; i < fType.NumIn(); i++ {
		switch fType.In(i) {
		case valueCtx.Type():
			params[i] = valueCtx
		case reflect.TypeOf((*Page)(nil)):
			params[i] = reflect.ValueOf(NewPageFromCtx(ctx, 15))
		case reflect.TypeOf((*Filters)(nil)):
			params[i] = reflect.ValueOf(NewFiltersFromCtx(ctx))
		case reflect.TypeOf((*Sorts)(nil)):
			params[i] = reflect.ValueOf(NewSortsFromCtx(ctx))
		default:
			params[i], err = parseParam(ctx, fType.In(i))
			if err != nil {
				return
			}
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
			err := errx.Wrap(x, "middleware error", errx.WithAbort())
			DefaultNewResponseFunc(nil, nil, err).Do(ctx)
		}
		return
	default:
		panic(errors.New("middleware should return results not more than 1"))
	}
}

func parseAndResponse(results ...reflect.Value) (resp IResponse[*gin.Context]) {
	switch len(results) {
	case 0:
		resp = DefaultNewResponseFunc(nil, nil, nil)
	case 1:
		r := results[0].Interface()
		if r == nil {
			resp = DefaultNewResponseFunc(nil, nil, nil)
			break
		}
		switch x := r.(type) {
		case error:
			resp = DefaultNewResponseFunc(nil, nil, x)
		default:
			resp = DefaultNewResponseFunc(r, nil, nil)
		}
	case 2:
		if results[1].IsValid() && !results[1].IsNil() {
			resp = DefaultNewResponseFunc(nil, nil, results[1].Interface().(error))
			break
		}
		resp = DefaultNewResponseFunc(results[0].Interface(), nil, nil)
	case 3:
		if results[2].IsValid() && !results[2].IsNil() {
			resp = DefaultNewResponseFunc(nil, nil, results[2].Interface().(error))
			break
		}

		if isNumber(results[0]) {
			tmp, err := reflectx.ConvertToKind(results[0], reflect.Int)
			if err != nil {
				panic(err)
			}
			resp = DefaultNewResponseFunc(results[1].Interface(), &Meta{Total: uint(tmp.Interface().(int))}, nil)
		} else if m, ok := results[0].Interface().(IMeta); ok {
			resp = DefaultNewResponseFunc(results[1].Interface(), m, nil)
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
