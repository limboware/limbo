package limbov1

import (
	"time"

	errnov1 "github.com/rejchev/errno"
)

type SystemManager struct {
	container []*System

	router map[uintptr]int
}

// OnDeActivate implements [ISystem].
func (x *SystemManager) Deactivate() {
	for i := range len(x.container) {
		if x.container[i].Deactivate != nil {
			x.container[i].Deactivate()
		}
	}
}

var systemManager = SystemManager{}

func Systems() *SystemManager {
	return &systemManager
}

func (x *SystemManager) Activate() bool {
	res := true
	for i := range len(x.container) {
		if x.container[i].Activate != nil {
			if res = x.container[i].Activate(); !res {
				return res
			}
		}
	}

	return res
}

func (x *SystemManager) Destroy() {
	for i := range len(x.container) {
		if x.container[i].Destroy != nil {
			x.container[i].Destroy()
		}
	}
}

func (x *SystemManager) Init() errnov1.Code {
	x.container = make([]*System, 0, 8)
	x.router = map[uintptr]int{}
	return errnov1.OK
}

func (x *SystemManager) Load() errnov1.Code {
	res := errnov1.OK
	for i := range len(x.container) {
		if x.container[i].Init != nil {
			if res = x.container[i].Init(); errnov1.FAIL(res) {
				return res
			}
		}
	}

	return res
}

func (x *SystemManager) Count() int {
	return len(x.container)
}

func (x *SystemManager) Register(allocFn func(*System) uintptr) bool {
	sys := new(System)
	if typePtr := allocFn(sys); typePtr != 0 {
		if !x.Contains(typePtr) {
			x.container = append(x.container, sys)
			x.router[typePtr] = len(x.container) - 1
			return true
		}
	}

	return false
}

func (x *SystemManager) System(v uintptr) *System {
	if !x.Contains(v) {
		return nil
	}

	return x.container[x.route(v)]
}

func (x *SystemManager) route(v uintptr) int {
	return x.router[v]
}

func (x *SystemManager) Contains(v uintptr) bool {
	_, ok := x.router[v]
	return ok
}

func (x *SystemManager) Update(dt time.Duration) {
	for i := range len(x.container) {
		if x.container[i].Update != nil {
			x.container[i].Update(dt)
		}
	}
}
