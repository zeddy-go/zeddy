package container

import (
	"errors"
	"github.com/zeddy-go/core/slicex"
	"reflect"
	"sync"
)

type Container struct {
	stuffs map[reflect.Type][]*Stuff
	lock   sync.Mutex
}

func (c *Container) Register(stuff *Stuff) {
	c.lock.Lock()
	defer c.lock.Unlock()

	stuff.SetContainer(c)
	tp := stuff.GetType()
	if _, ok := c.stuffs[tp]; !ok {
		c.stuffs[stuff.GetType()] = make([]*Stuff, 5)
	}

	c.stuffs[tp] = append(c.stuffs[tp], stuff)
}

func (c *Container) Invoke(f any) (result []reflect.Value, err error) {
	var (
		x  reflect.Value
		ok bool
	)

	if x, ok = f.(reflect.Value); !ok {
		x = reflect.ValueOf(f)
	}

	params := make([]reflect.Value, 0, x.Type().NumIn())
	for i := 0; i < x.Type().NumIn(); i++ {
		var p reflect.Value
		p, err = c.resolveValue(x.Type().In(i))
		if err != nil {
			return
		}
		params = append(params, p)
	}

	result = x.Call(params)

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
	stuffs, ok := c.stuffs[tp]
	if !ok {
		return reflect.Value{}, errors.New("stuff not found")
	}

	for _, item := range stuffs {
		if item.key == slicex.First(key...) {
			instance, err = item.GetInstance()
			return
		}
	}

	return reflect.Value{}, errors.New("stuff not found")
}
