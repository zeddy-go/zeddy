package contract

import "github.com/gin-gonic/gin"

type IShouldRegisterRoute interface {
	RegisterRoute(r gin.IRouter)
}

type IModule interface {
	Register(sub IModule)

	Init()
}
