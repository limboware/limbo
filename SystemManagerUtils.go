package limbov1

import (
	"reflect"
	"unsafe"
)

func SystemPtr[T any]() *T {
	if v := Systems().System(systemType[T]()); v != nil {
		return (*T)(v.Instance)
	}

	return nil
}

func SystemExists[T any]() bool {
	return Systems().Contains(systemType[T]())
}

func SystemRegister[T any](allocFn func(*System)) bool {
	return Systems().Register(func(x *System) uintptr {
		allocFn(x)
		return systemType[T]()
	})
}

func systemType[T any]() uintptr {
	t := reflect.TypeFor[*T]()
	return (*[2]uintptr)(unsafe.Pointer(&t))[1]
}
