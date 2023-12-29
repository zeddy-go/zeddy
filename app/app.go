package app

import (
	"github.com/zeddy-go/zeddy/contract"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var moduleList = make([]contract.IModule, 0)

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
		if module, ok := m.(contract.Startable); ok {
			go module.Start()
			n++
		}
	}

	return
}

func Stop() {
	for _, m := range moduleList {
		if module, ok := m.(contract.Stopable); ok {
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

	if n == 0 {
		slog.Info("nothing started")
	} else {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

		for range signals {
			signal.Stop(signals)
			close(signals)

			Stop()
		}

		println("bye bye~")
	}

	return
}
