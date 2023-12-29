package contract

type IStartable interface {
	//Start 启动服务并阻塞, 框架一般会将这个方法作为协程调用, 报错应打日志记录
	Start()
}

type IStopable interface {
	//Stop 停止服务并阻塞, 报错应打日志记录
	Stop()
}

type Service interface {
	IStartable
	IStopable
}
