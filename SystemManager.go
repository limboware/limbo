package limbov1

import (
	"reflect"
	"time"

	errnov1 "github.com/rejchev/errno"
)

var _ ISystem = (*SystemManager)(nil)

type SystemManager struct {
	systems []ISystem

	alias map[string]int
}

// OnDeActivate implements [ISystem].
func (x *SystemManager) Deactivate() {
	for _, y := range x.systems {
		y.Deactivate()
	}
}

var systemManager = SystemManager{
	systems: make([]ISystem, 0, 8),
	alias:   map[string]int{},
}

func Systems() *SystemManager {
	return &systemManager
}

// OnActivate implements [ISystemComponent].
func (x *SystemManager) Activate() bool {
	for _, y := range x.systems {
		if !y.Activate() {
			return false
		}
	}

	return true
}

// OnUnLoad implements [ISystemComponent].
func (x *SystemManager) Unload() {
	for _, y := range x.systems {
		y.Unload()
	}
}

func (x *SystemManager) Load() errnov1.Code {
	for _, y := range x.systems {
		if res := y.Load(); res != errnov1.OK {
			return res
		}
	}

	return errnov1.OK
}

func (x *SystemManager) OnAllLoaded() {
	for _, y := range x.systems {
		y.OnAllLoaded()
	}
}

func (x *SystemManager) Count() int {
	return len(x.systems)
}

func (x *SystemManager) SystemByIndex(v int) ISystem {
	return x.systems[v]
}

func (x *SystemManager) SystemTypeByIndex(v int) string {
	return reflect.TypeOf(x.systems[v]).String()
}

func (x *SystemManager) Create(k string, v ISystem) int {
	idx := len(x.systems)
	x.systems = append(x.systems, v)
	x.alias[k] = idx
	return idx
}

func (x *SystemManager) Destroy(v string) {
	if !x.Contains(v) {
		return
	}

	l := len(x.systems)
	if l == 1 {
		x.systems = make([]ISystem, 0, 8)
		delete(x.alias, v)
		return
	}

	idx := x.alias[v]
	delete(x.alias, v)

	last := l - 1
	if last == idx {
		x.systems = x.systems[:last]
		return
	}

	ls := x.systems[last]
	lt := reflect.TypeOf(ls).String()

	x.systems = x.systems[:last]
	x.systems[idx] = ls
	x.alias[lt] = idx
}

func (x *SystemManager) System(v string) ISystem {
	if !x.Contains(v) {
		return nil
	}

	return x.SystemByIndex(x.alias[v])
}

func (x *SystemManager) Contains(v string) bool {
	_, ok := x.alias[v]
	return ok
}

func (x *SystemManager) Update(dt time.Duration) {
	for _, v := range x.systems {
		v.Update(dt)
	}
}
