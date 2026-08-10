package limbov1

import (
	"time"
	"unsafe"

	errnov1 "github.com/rejchev/errno"
)

type InitFn = func() errnov1.Code
type ActivateFn = func() bool
type UpdateFn = func(time.Duration)
type DeactivateFn = func()
type DestroyFn = func()

type System struct {
	Instance   unsafe.Pointer
	Init       InitFn
	Activate   ActivateFn
	Update     UpdateFn
	Deactivate DeactivateFn
	Destroy    DestroyFn
}
