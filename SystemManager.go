package limbov1

import (
	"time"

	"github.com/limboware/limbo/debug"
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
	debug.Get().Printf("<SystemManager::Destroy> start (count: %d)", len(x.container))

	for i := range len(x.container) {
		if x.container[i].Destroy != nil {
			debug.Get().Printf("<SystemManager::Destroy> #%d system destroy", i)
			x.container[i].Destroy()
			debug.Get().Printf("<SystemManager::Destroy> #%d system destroyed", i)
		}
	}

	debug.Get().Printf("<SystemManager::Destroy> done")
}

func (x *SystemManager) Init() errnov1.Code {
	debug.Get().Printf("<SystemManager::Init> start")
	x.container = make([]*System, 0, 8)
	x.router = map[uintptr]int{}
	debug.Get().Printf("<SystemManager::Init> done")
	return errnov1.OK
}

func (x *SystemManager) Load() errnov1.Code {
	debug.Get().Printf("<SystemManager::Load> start (count: %d)", len(x.container))

	res := errnov1.OK
	for i := range len(x.container) {
		if x.container[i].Init != nil {
			debug.Get().Printf("<SystemManager::Load> #%d system init", i)
			res = x.container[i].Init()

			debug.Get().Printf("<SystemManager::Load> #%d system inited (code: %s)", i, res.String())

			if errnov1.FAIL(res) {
				return res
			}
		}
	}

	return res
}

func (x *SystemManager) Count() int {
	i := len(x.container)

	debug.Get().Printf("<SystemManager::Count> %d systems available", i)

	return i
}

func (x *SystemManager) Register(allocFn func(*System) uintptr) bool {
	debug.Get().Printf("<SystemManager::Register> #%d system %v", len(x.container), allocFn)

	sys := new(System)
	if typePtr := allocFn(sys); typePtr != 0 {

		debug.Get().Printf("<SystemManager::Register> #%d system %v is %v type", len(x.container), allocFn, typePtr)

		if !x.Contains(typePtr) {
			x.container = append(x.container, sys)
			x.router[typePtr] = len(x.container) - 1

			debug.Get().Printf("<SystemManager::Register> #%d system %v done", len(x.container)-1, allocFn)
			return true
		}
	}

	debug.Get().Printf("<SystemManager::Register> #%d system %v failed", len(x.container), allocFn)

	return false
}

func (x *SystemManager) System(v uintptr) *System {
	debug.Get().Printf("<SystemManager::System> type %v", v)

	if !x.Contains(v) {

		debug.Get().Printf("<SystemManager::System> type %v is unknown", v)
		return nil
	}

	debug.Get().Printf("<SystemManager::System> type %v is %v", v, x.container[x.route(v)])

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
