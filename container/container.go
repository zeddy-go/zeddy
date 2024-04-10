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

	if t == value.Type() || (t.Kind() == reflect.Interface && value.Type().Implements(t)) {
		c.bindConsistent(t, value, options)
	} else {
		delete(c.instances, t)
		err = c.bindProvider(t, value, options)
	}

	return
}

func (c *Container) bindProvider(t reflect.Type, value reflect.Value, options *bindOpts) (err error) {
	if value.Kind() != reflect.Func {
		err = errx.Wrap(ErrCanNotBind, fmt.Sprintf("value is not a func, type <%s>", t.String()))
		return
	}

	bindable, _ := c.isConsistent(t, value.Type().Out(0))
	if !bindable {
		err = errx.Wrap(ErrCanNotBind, fmt.Sprintf("type <%s>", t.String()))
		return
	}

	c.providers[t] = &provider{
		Value:     value,
		Singleton: options.Singleton,
	}

	return
}

func (c *Container) isConsistent(dest reflect.Type, src reflect.Type) (consistent bool, shouldConvert bool) {
	//一致性检查:
	// 1. 如果t是接口, 检查给定的**provider返回值或者实例**的类型是否实现了这个接口
	// 2. 如果t是普通类型, 检查给定的**provider返回值或者实例**的类型是否等于或者可以转换成t
	// 3. 判断按执行速度排序
	if dest == src {
		consistent = true
		return
	}

	if dest.Kind() == reflect.Interface && src.Implements(dest) {
		consistent = true
		return
	}

	if src.ConvertibleTo(dest) {
		consistent = true
		shouldConvert = true
		return
	}

	return
}

// bindConsistent 目标类型与给定值类型一致的情况
func (c *Container) bindConsistent(t reflect.Type, value reflect.Value, options *bindOpts) {
	//如果要绑定的类型是函数, 那就无视options中的singleton选项, 直接绑定到instances去
	if t.Kind() == reflect.Func {
		c.instances[t] = value
		return
	}

	//如果要绑定的类型不是函数, 就要判断是否为singleton, 来做合适的绑定
	if options.Singleton && value.Kind() == reflect.Pointer {
		c.instances[t] = value
		return
	}

	if options.Singleton && value.Kind() != reflect.Pointer {
		c.instances[t] = value.Addr()
		return
	}

	if value.Kind() == reflect.Pointer {
		c.providers[t] = &provider{
			Value: reflect.ValueOf(func() any {
				return reflect.New(value.Type().Elem()).Interface()
			}),
		}
	} else {
		c.providers[t] = &provider{
			Value: reflect.ValueOf(func() any {
				return reflect.New(value.Type()).Elem().Interface()
			}),
		}
	}
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
			_, shouldConvert := c.isConsistent(t, f.Value.Type().Out(0))
			if shouldConvert {
				c.instances[t] = result.Convert(t)
			} else {
				c.instances[t] = result
			}
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

	if len(results) == 2 {
		if !results[1].IsNil() {
			err = results[1].Interface().(error)
			return
		}
	}

	if results[0].IsNil() {
		err = ErrNotFound
		return
	}

	result = results[0]

	if result.Kind() == reflect.Interface {
		result = result.Elem()
	}

	consistent, shouldConvert := c.isConsistent(resultType, result.Type())
	if !consistent {
		err = ErrTypeNotMatch
		return
	}

	if shouldConvert {
		result = result.Convert(resultType)
	}

	return
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
