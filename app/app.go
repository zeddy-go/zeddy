package app

import (
	"github.com/zeddy-go/zeddy/container"
	"github.com/zeddy-go/zeddy/contract"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var moduleList = make([]contract.IModule, 0)

var beforeWaits = make([]any, 0)

// BeforeWaits 等待钩子
func BeforeWaits(funcs ...any) {
	beforeWaits = append(beforeWaits, funcs...)
}

func Use(modules ...contract.IModule) {
	moduleList = append(moduleList, modules...)
}

func Boot() (err error) {
	for _, module := range moduleList {
		if m, ok := module.(contract.IShouldInit); ok {
			err = m.Init()
			if err != nil {
				return
			}
		}
	}

	for _, module := range moduleList {
		if m, ok := module.(contract.IShouldBoot); ok {
			err = m.Boot()
			if err != nil {
				return
			}
		}
	}

	return
}

func Start() (n int) {
	for _, m := range moduleList {
		if module, ok := m.(contract.IStartable); ok {
			go module.Start()
			n++
		}
	}

	return
}

func Stop() {
	for _, m := range moduleList {
		if module, ok := m.(contract.IStopable); ok {
			go module.Stop()
		}
	}
}

func StartAndWait() (err error) {
	err = Boot()
	if err != nil {
		return
	}

	n := Start()
	defer Stop()

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
		}

		println("bye bye~")
	}

	return
}
