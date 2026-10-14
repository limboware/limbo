package limbov1

import (
	"time"

	"github.com/limboware/limbo/debug"
	errnov1 "github.com/rejchev/errno"
)

type World struct {
	tickCounter uint64
}

var world = World{}

func (x *World) Deactivate() {
	Systems().Deactivate()
}

func (x *World) Activate() bool {
	return Systems().Activate()
}

func (x *World) Destroy() {
	debug.Get().Printf("<World::Destroy> start")
	Systems().Destroy()
	debug.Get().Printf("<World::Destroy> done")
}

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
	debug.Get().Printf("<World::Init> start")
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

	debug.Get().Printf("<World::Init> done")

	return errnov1.OK
}

func (x *World) Ticks() uint64 {
	return x.tickCounter
}

func (x *World) Update(dt time.Duration) {

	Systems().Update(dt)

	x.tickCounter++
}
