package container

import (
	"errors"
	"reflect"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrCanNotBind   = errors.New("can not bind")
	ErrTypeNotMatch = errors.New("type not match")
)

func NewContainer() *Container {
	return &Container{
		providers: make(map[reflect.Type]reflect.Value),
		instances: make(map[reflect.Type]reflect.Value),
	}
}

type Container struct {
	providers map[reflect.Type]reflect.Value
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

func (c *Container) Bind(t reflect.Type, value reflect.Value, opts ...func(*bindOpts)) (err error) {
	options := &bindOpts{}
	for _, opt := range opts {
		opt(options)
	}

	if t == value.Type() {
		c.bindConsistent(t, value, options)
	} else {
		err = c.bindNotConsistent(t, value, options)
	}

	return
}

func (c *Container) bindNotConsistent(t reflect.Type, value reflect.Value, options *bindOpts) (err error) {
	if value.Kind() != reflect.Func {
		err = ErrCanNotBind
		return
	}

	bindable, shouldConvert := c.isConsistent(t, value.Type().Out(0))
	if !bindable {
		err = ErrCanNotBind
		return
	}

	if options.Singleton {
		var result reflect.Value
		result, err = c.invokeAndGetType(value, t)
		if err != nil {
			return
		}
		if shouldConvert {
			c.instances[t] = reflect.ValueOf(result).Convert(t)
		} else {
			c.instances[t] = result
		}
	} else {
		c.providers[t] = value
	}

	return
}

func (c *Container) isConsistent(dest reflect.Type, src reflect.Type) (canBind bool, shouldConvert bool) {
	//一致性检查:
	// 1. 如果t是接口, 检查给定的**provider返回值或者实例**的类型是否实现了这个接口
	// 2. 如果t是普通类型, 检查给定的**provider返回值或者实例**的类型是否等于或者可以转换成t
	// 3. 判断按执行速度排序
	if dest == src {
		canBind = true
		return
	}

	if dest.Kind() == reflect.Interface && src.Implements(dest) {
		canBind = true
		return
	}

	if src.ConvertibleTo(dest) {
		canBind = true
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
	if options.Singleton && t.Kind() == reflect.Pointer {
		c.instances[t] = value
		return
	}

	if options.Singleton && t.Kind() != reflect.Pointer {
		c.instances[t] = value.Addr()
		return
	}

	if t.Kind() == reflect.Pointer {
		c.providers[t] = reflect.ValueOf(func() any {
			return reflect.New(t.Elem()).Interface()
		})
	} else {
		c.providers[t] = reflect.ValueOf(func() any {
			return reflect.New(t).Elem().Interface()
		})
	}
}

func (c *Container) Resolve(t reflect.Type) (result reflect.Value, err error) {
	result, ok := c.instances[t]
	if ok {
		return
	}

	f, ok := c.providers[t]
	if ok {
		result, err = c.invokeAndGetType(f, t)
		if err != nil {
			return
		}
		return
	}

	err = ErrNotFound
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

	result = results[0]

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
