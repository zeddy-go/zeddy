package contract

type Startable interface {
	Start() //阻塞
}

type Stopable interface {
	Stop()
}

type Service interface {
	Startable
	Stopable
}
