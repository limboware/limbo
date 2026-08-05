package limbov1

import errnov1 "github.com/rejchev/errno"

var netconntypes = NetConnTypeManager{}

type NetConnTypeManager struct {
	count     uint8
}

func NetConnTypes() *NetConnTypeManager {
	return &netconntypes
}

func (x *NetConnTypeManager) Init() errnov1.Code {
	x.count = 1
	return errnov1.OK
}

func (x *NetConnTypeManager) Register(buff *NetConnType) bool {
	if *buff = NetConnType(x.count); x.count >= MaxNetConnTypes {
		return false
	}

	x.count++

	return true
}

func (x *NetConnTypeManager) Count() int {
	return int(x.count) - 1
}

func (x *NetConnTypeManager) Reset() {
	x.count = 1
}
