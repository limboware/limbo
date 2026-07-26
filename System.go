package limbov1

import (
	"time"

	errnov1 "github.com/rejchev/errno"
)

type ISystem interface {
	Load() errnov1.Code

	OnAllLoaded()

	Activate() bool
	Update(dt time.Duration)
	Deactivate()

	Unload()
}
