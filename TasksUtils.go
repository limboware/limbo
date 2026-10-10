package limbov1

import (
	"reflect"
	"unsafe"
)

func NewTask[T any](allocFn func(*Task_s), buff *Task) bool {
	return Tasks().New(func(t *Task_s) uintptr {
		allocFn(t)
		return TypePtr[T]()
	}, buff)
}

func IsTypeEqual[T any](another uintptr) bool {
	return TypePtr[T]() == another
} 

func TypePtr[T any]() uintptr {
	t := reflect.TypeFor[*T]()
	return (*[2]uintptr)(unsafe.Pointer(&t))[1]
}
