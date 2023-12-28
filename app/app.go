package app

import "github.com/zeddy-go/zeddy/contract"

var modules = []contract.IModule{}
var stopers = []func(){}
var starters = []func(){}

func Use(module contract.IModule) {
	modules = append(modules, module)
}

func Boot() (err error) {
	for _, module := range modules {
		if m, ok := module.(contract.IShouldBoot); ok {
			err = m.Boot()
			if err != nil {
				return
			}
		}
	}

	return
}

func Start() (err error) {
	for _, module := range modules {
		if m, ok := module.(contract.Service); ok {
			starters = append(starters, m.Start)
			stopers = append(stopers, m.Stop)
		}
	}

	//TODO: 启动
	return
}
