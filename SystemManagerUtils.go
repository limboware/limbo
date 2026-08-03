package limbov1

import (
	"reflect"
	"unsafe"
)

func SystemPtr[T any]() *T {
	if v := Systems().System(reflect.TypeFor[T]().Name()); v != nil {
		return (*T)(unsafe.Pointer(reflect.ValueOf(v).Pointer()))
	}

	return nil
}

func SystemExists[T any]() bool {
	return Systems().Contains(reflect.TypeFor[T]().Name())
}

func SystemRegister[T ISystem](v T) bool {
	return Systems().Create(reflect.TypeFor[T]().Name(), v) != -1
}