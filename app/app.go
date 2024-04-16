package app

import (
	"github.com/zeddy-go/zeddy/container"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var moduleList = make([]Module, 0)

var beforeWaits = make([]any, 0)

// BeforeWaits 等待钩子
func BeforeWaits(funcs ...any) {
	beforeWaits = append(beforeWaits, funcs...)
}

func Use(modules ...Module) {
	moduleList = append(moduleList, modules...)
}

func Boot() (err error) {
	for _, module := range moduleList {
		if m, ok := module.(Initable); ok {
			err = m.Init()
			if err != nil {
				return
			}
		}
	}

	for _, module := range moduleList {
		if m, ok := module.(Bootable); ok {
			err = m.Boot()
			if err != nil {
				return
			}
		}
	}

	return
}

func Start(wg *sync.WaitGroup) (n int) {
	for _, m := range moduleList {
		if module, ok := m.(Service); ok {
			wg.Add(1)
			go func(start func()) {
				defer wg.Done()
				start()
			}(module.Start)
			n++
		}
	}

	return
}

func Stop() {
	for _, m := range moduleList {
		if module, ok := m.(Service); ok {
			go module.Stop()
		}
	}
}

func StartAndWait() (err error) {
	err = Boot()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	n := Start(&wg)

	if n == 0 {
		slog.Info("nothing started, shutdown.")
	} else {
		if len(beforeWaits) > 0 {
			for _, f := range beforeWaits {
				err = container.Invoke(f)
				if err != nil {
					return
				}
			}
		}
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

		for range signals {
			signal.Stop(signals)
			close(signals)
			Stop()
		}

		println("wait module stop...")
		wg.Wait()
		println("bye bye~")
	}

	return
}
