package module

func NewBaseModule(name string) *BaseModule {
	return &BaseModule{
		name: name,
	}
}

type BaseModule struct {
	name string
}

func (b BaseModule) Name() string {
	return b.name
}
