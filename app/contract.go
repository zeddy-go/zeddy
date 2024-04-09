package app

type HasSubModule interface {
	Register(subs ...Module) error
}

type Bootable interface {
	Boot() error
}

type HasName interface {
	Name() string
}

type Module interface {
	module()
}

type Initable interface {
	Init() error
}

type Startable interface {
	//Start 启动服务并阻塞, 框架一般会将这个方法作为协程调用, 报错应打日志记录
	Start()
}

type Stopable interface {
	//Stop 停止服务并阻塞, 报错应打日志记录
	Stop()
}

type Service interface {
	Startable
	Stopable
}
