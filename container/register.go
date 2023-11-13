package container

func Register(providerOrInstance any, sets ...func(*Stuff)) {
	Default.Register(NewStuff(providerOrInstance, sets...))
}
