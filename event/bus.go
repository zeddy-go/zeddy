package event

import (
	"errors"
	"github.com/zeddy-go/zeddy/container"
	"reflect"
	"sync"
)

type Bus struct {
	lock      sync.RWMutex
	subs      map[reflect.Type][]any
	container *container.Container
}

// Sub 订阅事件, f必须是函数, f的第一个参数必须是事件
func (h *Bus) Sub(f any) {
	h.lock.Lock()
	defer h.lock.Unlock()

	vFunc := reflect.ValueOf(f)
	if vFunc.Kind() != reflect.Func {
		panic(errors.New("func only"))
	}

	if h.container == nil && vFunc.Type().NumIn() != 1 {
		panic(errors.New("event handler require must one param in normal mode"))
	} else if h.container != nil && vFunc.Type().NumIn() < 1 {
		panic(errors.New("event handler require at least one param in container mode"))
	}

	eventType := vFunc.Type().In(0)

	if _, ok := h.subs[eventType]; !ok {
		h.subs[eventType] = make([]any, 0, 10)
	}

	h.subs[eventType] = append(h.subs[eventType], f)
}

func (h *Bus) Pub(event any) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	eventValue := reflect.ValueOf(event)

	if group, ok := h.subs[eventValue.Type()]; ok {
		for _, item := range group {
			go func(f any) {
				if h.container != nil {
					h.container.Invoke(f, container.WithInvokeParams(map[int]any{0: event}))
				} else {
					reflect.ValueOf(f).Call([]reflect.Value{eventValue})
				}
			}(item)
		}
	}
}
