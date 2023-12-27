package container

import (
	"errors"
	"fmt"
	"github.com/zeddy-go/zeddy/slicex"
	"reflect"
	"sync"
)

func NewContainer() *Container {
	return &Container{
		stuffs: make(map[reflect.Type][]*Stuff),
	}
}

type Container struct {
	stuffs map[reflect.Type][]*Stuff
	lock   sync.RWMutex
}

func (c *Container) BindType(destType reflect.Type, srcType reflect.Type) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	stuffs, ok := c.stuffs[srcType]
	if !ok {
		return
	}

	c.stuffs[destType] = stuffs
}

// Bind 实列或者provider函数到指定的类型上,
// 它会检查**provider返回值/实例**的类型是否与指定的type一致.
func (c *Container) Bind(t reflect.Type, providerOrInstance any, sets ...func(*Stuff)) {
	typeV := reflect.TypeOf(providerOrInstance)
	if typeV.Kind() == reflect.Func {
		typeV = typeV.Out(0) //固定期望第一个就是provider返回的值
	}

	if !c.bindable(t, typeV) {
		panic(errors.New(fmt.Sprintf("type [%s] can not bind to type [%s]", typeV.String(), t.String())))
	}

	stuff := NewStuff(providerOrInstance, sets...)
	stuff.SetContainer(c)

	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.stuffs[t]; !ok {
		c.stuffs[t] = make([]*Stuff, 0, 5)
	}

	c.stuffs[t] = append(c.stuffs[t], stuff)
}

func (c *Container) bindable(dest reflect.Type, src reflect.Type) bool {
	//一致性检查:
	// 1. 如果t是接口, 检查给定的**provider返回值或者实例**的类型是否实现了这个接口
	// 2. 如果t是普通类型, 检查给定的**provider返回值或者实例**的类型是否等于或者可以转换成t
	// 3. 判断按执行速度排序
	if dest == src {
		return true
	}

	if dest.Kind() == reflect.Interface && src.Implements(dest) {
		return true
	}

	if src.ConvertibleTo(dest) {
		return true
	}

	return false
}

func (c *Container) Register(stuff *Stuff) {
	c.lock.Lock()
	defer c.lock.Unlock()

	stuff.SetContainer(c)
	tp := stuff.GetType()
	if _, ok := c.stuffs[tp]; !ok {
		c.stuffs[tp] = make([]*Stuff, 0, 5)
	}

	c.stuffs[tp] = append(c.stuffs[tp], stuff)
}

type InvokeOpts struct {
	params map[int]any //map[参数序号]参数
}

func WithInvokeParams(params map[int]any) func(*InvokeOpts) {
	return func(opts *InvokeOpts) {
		opts.params = params
	}
}

// Invoke 执行一个函数,函数参数通过容器注入,如果调用函数没有问题,那么err返回值尝试用被调用函数的最后一个返回值.
func (c *Container) Invoke(f any, sets ...func(*InvokeOpts)) (result []reflect.Value, err error) {
	var (
		x    reflect.Value
		ok   bool
		opts InvokeOpts
	)

	for _, set := range sets {
		set(&opts)
	}

	if x, ok = f.(reflect.Value); !ok {
		x = reflect.ValueOf(f)
	}

	params := make([]reflect.Value, 0, x.Type().NumIn())
	for i := 0; i < x.Type().NumIn(); i++ {
		var p reflect.Value
		if param, ok := opts.params[i]; ok {
			p = reflect.ValueOf(param)
		} else {
			p, err = c.resolveValue(x.Type().In(i))
			if err != nil {
				return
			}
		}
		params = append(params, p)
	}

	result = x.Call(params)

	if len(result) > 0 {
		if e, ok := result[len(result)-1].Interface().(error); ok {
			err = e
		}
	}

	return
}

func (c *Container) Resolve(tp reflect.Type, key ...string) (instance any, err error) {
	tmp, err := c.resolveValue(tp, key...)
	if err != nil {
		return
	}

	instance = tmp.Interface()

	return
}

func (c *Container) resolveValue(tp reflect.Type, key ...string) (instance reflect.Value, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	stuffs, ok := c.stuffs[tp]
	if !ok {
		return reflect.Value{}, errors.New("stuff not found")
	}

	for _, item := range stuffs {
		if item.Key == slicex.First(key...) {
			instance, err = item.GetInstance()
			return
		}
	}

	return reflect.Value{}, errors.New("stuff not found")
}
