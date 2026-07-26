package limbov1

import (
	errnov1 "github.com/rejchev/errno"
)

type NetworkHandlerFn func(Entity, *byte, uint32) errnov1.Code

type NetworkHandlerManager struct {
	container map[NetworkHandler]NetworkHandlerFn
}

var msgH = NetworkHandlerManager{container: map[NetworkHandler]NetworkHandlerFn{}}

func NetworkHandlers() *NetworkHandlerManager {
	return &msgH
}

func (x *NetworkHandlerManager) Init() errnov1.Code {
	return errnov1.OK
}

func (x *NetworkHandlerManager) Exist(v NetworkHandler) bool {
	_, ok := x.container[v]
	return ok
}

func (x *NetworkHandlerManager) Get(v NetworkHandler) NetworkHandlerFn {
	y, _ := x.container[v]
	return y
}

func (x *NetworkHandlerManager) Register(proto ConnectionProtocol, typ NetMessageType, fn NetworkHandlerFn) bool {
	netHandler := NetworkHandlerFor(proto, typ)

	if !x.Exist(netHandler) {
		x.container[netHandler] = fn
		return true
	}

	return false
}

func (x *NetworkHandlerManager) Unregister(v NetworkHandler) {
	delete(x.container, v)
}
