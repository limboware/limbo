package limbov1

import (
	errnov1 "github.com/rejchev/errno"
)

type NetHandlerFn func(Entity, *byte, uint32) errnov1.Code

type NetHandlerManager struct {
	container map[NetworkHandler]NetHandlerFn
}

var msgH = NetHandlerManager{}

func NetworkHandlers() *NetHandlerManager {
	return &msgH
}

func (x *NetHandlerManager) Init() errnov1.Code {
	x.container = map[NetworkHandler]NetHandlerFn{}
	return errnov1.OK
}

func (x *NetHandlerManager) Exist(v NetworkHandler) bool {
	_, ok := x.container[v]
	return ok
}

func (x *NetHandlerManager) Get(v NetworkHandler) NetHandlerFn {
	y, _ := x.container[v]
	return y
}

func (x *NetHandlerManager) Register(nett NetConnType, msgt NetMessageType, fn NetHandlerFn) bool {
	netHandler := MakeNetworkHandler(nett, msgt)

	if !x.Exist(netHandler) {
		x.container[netHandler] = fn
		return true
	}

	return false
}

func (x *NetHandlerManager) Unregister(v NetworkHandler) {
	delete(x.container, v)
}
