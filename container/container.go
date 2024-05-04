package container

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeddy-go/zeddy/errx"
	"reflect"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
)

type provider struct {
	Value     reflect.Value
	Singleton bool
}

type bindOpts struct {
	Singleton bool
	Key       string
}

type resolveOpts struct {
	Key string
}

func NoSingleton() func(*bindOpts) {
	return func(opts *bindOpts) {
		opts.Singleton = false
	}
}

func WithKey(key string) func(*bindOpts) {
	return func(opts *bindOpts) {
		opts.Key = key
	}
}

func WithResolveKey(key string) func(*resolveOpts) {
	return func(opts *resolveOpts) {
		opts.Key = key
	}
}

func NewContainer() *Container {
	return &Container{
		providers: make(map[reflect.Type]map[string]*provider),
		instances: make(map[reflect.Type]map[string]reflect.Value),
	}
}

type Container struct {
	providers map[reflect.Type]map[string]*provider
	instances map[reflect.Type]map[string]reflect.Value
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
		err = c.bindInstance(t, value, options)
	} else if canBindProvider(t, value) {
		err = c.bindProvider(t, value, options)
	} else {
		err = errx.New(fmt.Sprintf("can not bind <%s> to <%s>", value.Type(), t))
	}

	return
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

// bindConsistent 目标类型与给定值类型一致的情况
func (c *Container) bindInstance(t reflect.Type, value reflect.Value, options *bindOpts) (err error) {
	if t != value.Type() {
		value, err = c.convert(t, value)
		if err != nil {
			return
		}
	}

	group, ok := c.instances[t]
	if !ok {
		group = make(map[string]reflect.Value)
		c.instances[t] = group
	}
	group[options.Key] = value

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

func (c *Container) bindProvider(t reflect.Type, value reflect.Value, options *bindOpts) (err error) {
	group, ok := c.providers[t]
	if !ok {
		group = make(map[string]*provider)
		c.providers[t] = group
	}
	group[options.Key] = &provider{
		Value:     value,
		Singleton: options.Singleton,
	}

	return
}

func (c *Container) Resolve(t reflect.Type, opts ...func(*resolveOpts)) (result reflect.Value, err error) {
	options := &resolveOpts{}
	for _, opt := range opts {
		opt(options)
	}
	return c.resolve(context.Background(), t, options)
}

func (c *Container) resolve(ctx context.Context, t reflect.Type, opts *resolveOpts) (result reflect.Value, err error) {
	group, ok := c.instances[t]
	if ok {
		result, ok = group[opts.Key]
		if ok {
			return
		}
	}

	var chain []reflect.Type
	tmp := ctx.Value("chain")
	if tmp == nil {
		chain = []reflect.Type{}
		ctx = context.WithValue(ctx, "chain", chain)
	} else {
		chain = tmp.([]reflect.Type)
	}

	for _, item := range chain {
		if item == t {
			result = reflect.New(t).Elem()
			return
		}
	}

	group2, ok := c.providers[t]
	if ok {
		var f *provider
		f, ok = group2[opts.Key]
		if ok {
			chain = append(chain, t)
			ctx = context.WithValue(ctx, "chain", chain)
			result, err = c.invokeAndGetType(ctx, f.Value, t)
			if err != nil {
				return
			}

			if f.Singleton {
				group, ok = c.instances[t]
				if !ok {
					group = make(map[string]reflect.Value)
					c.instances[t] = group
				}
				group[opts.Key] = result
			}

			return
		}
	}

	result, err = c.resolveSlow(ctx, t, opts)
	if err != nil {
		return
	}
	if !result.IsValid() {
		err = errx.Wrap(ErrNotFound, fmt.Sprintf("type <%s>", t.String()))
		return
	}

	return
}

// resolveSlow 尝试对已有类型进行转换
func (c *Container) resolveSlow(ctx context.Context, t reflect.Type, opts *resolveOpts) (result reflect.Value, err error) {
	chain := ctx.Value("chain").([]reflect.Type)

	for typ, group := range c.instances {
		if typ.ConvertibleTo(t) {
			if item, ok := group[opts.Key]; ok {
				result = item.Convert(t)
				group, ok := c.instances[t]
				if !ok {
					group = make(map[string]reflect.Value)
					c.instances[t] = group
				}
				group[opts.Key] = result
				chain = append(chain, t)
				ctx = context.WithValue(ctx, "chain", chain)
				return
			}
		}
	}

	for typ, group := range c.providers {
		if typ.ConvertibleTo(t) {
			if provider, ok := group[opts.Key]; ok {
				chain = append(chain, t)
				ctx = context.WithValue(ctx, "chain", chain)

				result, err = c.invokeAndGetType(ctx, provider.Value, t)
				if err != nil {
					return
				}

				if provider.Singleton {
					group, ok := c.instances[t]
					if !ok {
						group = make(map[string]reflect.Value)
						c.instances[t] = group
					}
					group[opts.Key] = result
				}

				return
			}
		}
	}

	return
}

func (c *Container) invokeAndGetType(ctx context.Context, f reflect.Value, resultType reflect.Type) (result reflect.Value, err error) {
	results, err := c.invoke(ctx, f)
	if err != nil {
		return
	}

	//if results[0].IsNil() {
	//	err = errx.New("get nil result")
	//	return
	//}

	result = results[0]

	if result.Kind() == reflect.Interface {
		result = result.Elem()
	}

	result, err = c.convert(resultType, result)
	return
}

func WithParams(params map[int]any) func(*invokeOpts) {
	return func(opts *invokeOpts) {
		opts.params = params
	}
}

type invokeOpts struct {
	params map[int]any
	keys   map[int]string
}

var waits = make(map[reflect.Type][]reflect.Value)
var waitsLock sync.Mutex

func (c *Container) invoke(ctx context.Context, f reflect.Value, opts ...func(*invokeOpts)) (results []reflect.Value, err error) {
	options := &invokeOpts{}
	for _, opt := range opts {
		opt(options)
	}

	p := make([]reflect.Value, 0, f.Type().NumIn())
	ts := make([]reflect.Type, 0, f.Type().NumIn())
	for i := 0; i < f.Type().NumIn(); i++ {
		var param reflect.Value
		if len(options.params) <= 0 {
			param, err = c.resolve(ctx, f.Type().In(i), &resolveOpts{Key: options.keys[i]})
			if err != nil {
				return
			}
			if param.IsNil() {
				ts = append(ts, f.Type().In(i))
			}
		} else {
			if p, ok := options.params[i]; ok {
				param = reflect.ValueOf(p)
			}
		}

		p = append(p, param)
	}

	results = f.Call(p)

	if len(results) > 0 && !results[len(results)-1].IsNil() {
		var ok bool
		if err, ok = results[len(results)-1].Interface().(error); ok {
			return
		}
	}

	if len(waits) > 0 {
		waitsLock.Lock()
		for _, result := range results {
			if targets, ok := waits[result.Type()]; ok {
				for _, target := range targets {
					target.Elem().Set(result)
				}
				delete(waits, result.Type())
			}
		}
		waitsLock.Unlock()
	}

	if len(ts) > 0 {
		waitsLock.Lock()
		for _, result := range results {
			if result.IsNil() {
				continue
			}
			v := reflect.Indirect(result)
			if v.Kind() != reflect.Struct {
				continue
			}

			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				for _, t := range ts {
					if t == field.Type() {
						target := reflect.NewAt(t, field.Addr().UnsafePointer())
						targets, ok := waits[t]
						if !ok {
							targets = make([]reflect.Value, 0, 1)
						}
						targets = append(targets, target)
						waits[t] = targets
						break
					}
				}
			}
		}
		waitsLock.Unlock()
	}

	return
}

func (c *Container) Invoke(f reflect.Value, opts ...func(*invokeOpts)) (results []reflect.Value, err error) {
	return c.invoke(context.Background(), f, opts...)
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
