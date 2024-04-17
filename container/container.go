package container

import (
	"errors"
	"fmt"
	"github.com/zeddy-go/zeddy/errx"
	"reflect"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrCanNotBind   = errors.New("can not bind")
	ErrTypeNotMatch = errors.New("type not match")
)

type provider struct {
	Value     reflect.Value
	Singleton bool
}

func NewContainer() *Container {
	return &Container{
		providers: make(map[reflect.Type]*provider),
		instances: make(map[reflect.Type]reflect.Value),
	}
}

type Container struct {
	providers map[reflect.Type]*provider
	instances map[reflect.Type]reflect.Value
}

type bindOpts struct {
	Singleton bool
}

// AsSingleton enable singleton mode
// Deprecated: singleton mode is now default enabled
func AsSingleton() func(*bindOpts) {
	return func(opts *bindOpts) {
		opts.Singleton = true
	}
}

func NoSingleton() func(*bindOpts) {
	return func(opts *bindOpts) {
		opts.Singleton = false
	}
}

func (c *Container) Bind(t reflect.Type, value reflect.Value, opts ...func(*bindOpts)) (err error) {
	options := &bindOpts{
		Singleton: true,
	}
	for _, opt := range opts {
		opt(options)
	}

	delete(c.instances, t)
	if canBindConsistent(t, value) {
		c.bindConsistent(t, value, options)
	} else if canBindProvider(t, value) {
		err = c.bindProvider(t, value, options)
	} else {
		err = errx.New(fmt.Sprintf("can not bind <%s> to <%s>", value.Type(), t))
	}

	return
}

// canBindProvider
func canBindProvider(t reflect.Type, value reflect.Value) bool {
	if value.Kind() != reflect.Func {
		return false
	}

	if value.Type().NumOut() <= 0 {
		return false
	}

	resultType := value.Type().Out(0)

	if t == resultType {
		return true
	}

	if resultType.ConvertibleTo(t) {
		return true
	}

	return false
}

func canBindConsistent(t reflect.Type, value reflect.Value) bool {
	if t == value.Type() {
		return true
	}

	if value.Type().ConvertibleTo(t) {
		return true
	}

	return false
}

func (c *Container) bindProvider(t reflect.Type, value reflect.Value, options *bindOpts) (err error) {
	c.providers[t] = &provider{
		Value:     value,
		Singleton: options.Singleton,
	}

	return
}

// bindConsistent 目标类型与给定值类型一致的情况
func (c *Container) bindConsistent(t reflect.Type, value reflect.Value, options *bindOpts) (err error) {
	//如果要绑定的类型是函数, 那就无视options中的singleton选项, 直接绑定到instances去
	if t.Kind() == reflect.Func {
		c.instances[t] = value
		return
	}

	//如果要绑定的类型不是函数, 就要判断是否为singleton, 来做合适的绑定
	if options.Singleton {
		var tmp reflect.Value
		tmp, err = c.convert(t, value)
		if err != nil {
			return
		}
		c.instances[t] = tmp
		return
	}

	c.providers[t] = &provider{
		Value: reflect.ValueOf(func() any {
			return reflect.New(value.Type().Elem()).Interface()
		}),
	}

	return
}

func (c *Container) convert(t reflect.Type, value reflect.Value) (result reflect.Value, err error) {
	if t == value.Type() {
		result = value
	} else if value.Type().ConvertibleTo(t) {
		result = value.Convert(t)
	} else {
		err = errx.New(fmt.Sprintf("can not convert <%s> to <%s>", value.Type().String(), t.String()))
	}
	return
}

func (c *Container) Resolve(t reflect.Type) (result reflect.Value, err error) {
	result, ok := c.instances[t]
	if ok {
		return
	}

	f, ok := c.providers[t]
	if ok {
		result, err = c.invokeAndGetType(f.Value, t)
		if err != nil {
			return
		}
		if f.Singleton {
			c.instances[t] = result
		}
		return
	}

	err = errx.Wrap(ErrNotFound, fmt.Sprintf("type <%s>", t.String()))
	return
}

func (c *Container) invokeAndGetType(f reflect.Value, resultType reflect.Type) (result reflect.Value, err error) {
	results, err := c.Invoke(f)
	if err != nil {
		return
	}

	if len(results) == 0 {
		err = errors.New("no result returned")
		return
	}

	if len(results) == 2 && !results[1].IsNil() {
		err = results[1].Interface().(error)
		return
	}

	if results[0].IsNil() {
		err = errx.New("get nil result")
		return
	}

	result = results[0]

	if result.Kind() == reflect.Interface {
		result = result.Elem()
	}

	return c.convert(resultType, result)
}

func WithParams(params map[int]any) func(*invokeOpts) {
	return func(opts *invokeOpts) {
		opts.params = params
	}
}

type invokeOpts struct {
	params map[int]any
}

func (c *Container) Invoke(f reflect.Value, opts ...func(*invokeOpts)) (results []reflect.Value, err error) {
	options := &invokeOpts{}
	for _, opt := range opts {
		opt(options)
	}

	p := make([]reflect.Value, 0, f.Type().NumIn())
	for i := 0; i < f.Type().NumIn(); i++ {
		var param reflect.Value
		if len(options.params) > 0 {
			if p, ok := options.params[i]; ok {
				param = reflect.ValueOf(p)
			} else {
				param, err = c.Resolve(f.Type().In(i))
				if err != nil {
					return
				}
			}
		} else {
			param, err = c.Resolve(f.Type().In(i))
			if err != nil {
				return
			}
		}

		p = append(p, param)
	}

	results = f.Call(p)

	return
}

func (c *Container) Has(t reflect.Type) bool {
	_, ok := c.instances[t]
	if ok {
		return true
	}

	_, ok = c.providers[t]
	if ok {
		return true
	}

	return false
}
