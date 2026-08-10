package limbov1

import (
	"time"

	errnov1 "github.com/rejchev/errno"
)

type World struct {
	tickCounter uint64
}

func (x *World) Deactivate() {
	Systems().Deactivate()
}

func (x *World) Activate() bool {
	return Systems().Activate()
}

func (x *World) Destroy() {
	Systems().Destroy()
}

var world = World{}

func GetWorld() *World {
	return &world
}

func (x *World) CreateEntity() Entity {
	return Entities().Create()
}

func (x *World) CreateCompotype(allocFn CompotypeAllocator, buff *Compotype) bool {
	return Compotypes().Register(allocFn, buff)
}

func (x *World) Init() errnov1.Code {
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

	Events().Publish("world.loaded", nil)

	return errnov1.OK
}

func (x *World) Ticks() uint64 {
	return x.tickCounter
}

func (x *World) Update(dt time.Duration) {

	Systems().Update(dt)

	x.tickCounter++
}
