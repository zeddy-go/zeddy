package module

import "github.com/zeddy-go/core/contract"

func Init(list ...contract.IModule) {
	for _, item := range list {
		item.Init()
	}
}
