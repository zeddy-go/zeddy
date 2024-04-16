package app

// HasSubModule 接口表示模块可以包含子模块
type HasSubModule interface {
	Register(subs ...Module) error
}

type HasName interface {
	Name() string
}

// Module 表示实现了这个接口的结构体是一个框架模块。
type Module interface {
	module()
}

// Initable 表示模块需要被初始化
type Initable interface {
	Init() error
}

// Bootable 表示模块需要被启动
type Bootable interface {
	Boot() error
}

type Service interface {
	//Start 启动服务并阻塞, 框架一般会将这个方法作为协程调用, 报错应打日志记录
	Start()
	//Stop 停止服务并阻塞, 报错应打日志记录
	Stop()
}
