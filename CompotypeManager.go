package limbov1

import (
	"unsafe"

	errnov1 "github.com/rejchev/errno"
)

const MaxCompotypes = 1<<8 - 1

type CompotypeAllocator func() unsafe.Pointer

type CompotypeManager struct {
	alloc []CompotypeAllocator
}

var componentTypes = CompotypeManager{
	alloc: make([]CompotypeAllocator, MaxCompotypes)[:0],
}

func Compotypes() *CompotypeManager {
	return &componentTypes
}

func (x *CompotypeManager) Init() errnov1.Code {
	return errnov1.OK
}

func (x *CompotypeManager) Register(allocFn CompotypeAllocator, buff *Compotype) bool {
	if *buff = Compotype(len(x.alloc)); (*buff + 1) == MaxCompotypes {
		return false
	}

	x.alloc = append(x.alloc, allocFn)

	return true
}

func (x *CompotypeManager) Allocate(v Compotype) unsafe.Pointer {
	return x.alloc[v.Int()]()
}
