package limbov1

import (
	"time"

	errnov1 "github.com/rejchev/errno"
)

var _ ISystem = (*World)(nil)

type World struct {
	tickCounter uint64
}

// OnDeActivate implements [ISystem].
func (x *World) Deactivate() {
	Systems().Deactivate()
}

// OnActivate implements [ISystem].
func (x *World) Activate() bool {
	return Systems().Activate()
}

// OnAllLoaded implements [ISystem].
func (x *World) OnAllLoaded() {
	Systems().OnAllLoaded()
}

// OnUnLoad implements [ISystem].
func (x *World) Unload() {
	Systems().Unload()
}

var world = World{tickCounter: 0}

func GetWorld() *World {
	return &world
}

func (x *World) CreateSystem(name string, v ISystem) int {
	return Systems().Create(name, v)
}

func (x *World) CreateEntity() Entity {
	return Entities().Create()
}

func (x *World) CreateCompotype(allocFn CompotypeAllocator, buff *Compotype) bool {
	return Compotypes().Register(allocFn, buff)
}

func (x *World) Load() errnov1.Code {
	x.tickCounter = 0

	if err := Entities().Init(); err != errnov1.OK {
		return err
	}

	if err := Compotypes().Init(); err != errnov1.OK {
		return err
	}

	if err := Components().Init(); err != errnov1.OK {
		return err
	}

	if err := Systems().Load(); err != errnov1.OK {
		return err
	}

	return errnov1.OK
}

func (x *World) Ticks() uint64 {
	return x.tickCounter
}

func (x *World) Update(dt time.Duration) {

	Systems().Update(dt)

	x.tickCounter++
}
