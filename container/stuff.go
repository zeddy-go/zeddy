package container

import (
	"errors"
	"reflect"
)

func NewStuff(instanceOrProvider any, sets ...func(*Stuff)) *Stuff {
	v := reflect.ValueOf(instanceOrProvider)
	switch v.Kind() {
	case reflect.Func:
		return NewStuffUseProvider(v, sets...)
	default:
		return NewStuffUseInstance(v, sets...)
	}
}

func NewStuffUseProvider(provider reflect.Value, sets ...func(*Stuff)) *Stuff {
	s := &Stuff{
		provider: provider,
	}

	for _, set := range sets {
		set(s)
	}

	return s
}

func NewStuffUseInstance(instance reflect.Value, sets ...func(*Stuff)) *Stuff {
	s := &Stuff{
		instance: instance,
	}

	for _, set := range sets {
		set(s)
	}

	return s
}

type Stuff struct {
	key       string        //键
	provider  reflect.Value //实例化函数
	instance  reflect.Value //实例
	singleton bool          //是否单例
	container *Container    //容器
}

func (s *Stuff) SetContainer(container *Container) {
	s.container = container
}

func (s *Stuff) GetType() reflect.Type {
	if s.instance.IsValid() {
		return s.instance.Type()
	}

	if s.provider.IsValid() {
		return s.provider.Type().Out(0)
	}

	panic(errors.New("stuff is invalid for lack of both instance and provider"))
}

func (s *Stuff) create() (instance reflect.Value, err error) {
	results, err := s.container.Invoke(s.provider)
	if err != nil {
		return
	}
	if len(results) > 1 {
		if results[1].Interface() != nil {
			err = results[1].Interface().(error)
			return
		}
	}

	if results[0].IsValid() {
		err = errors.New("no valid result")
		return
	}

	return results[0], nil
}

func (s *Stuff) GetInstance() (v reflect.Value, err error) {
	if s.singleton {
		if s.instance.IsValid() {
			s.instance, err = s.create()
			if err != nil {
				return
			}
		}

		return s.instance, nil
	} else {
		return s.create()
	}
}
