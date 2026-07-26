package limbov1

import (
	"unsafe"

	errnov1 "github.com/rejchev/errno"
)

type location uint64

func locationOf(i int, j int) location {
	return location(uint64(uint8(i))<<32 | uint64(uint32(j)))
}

func (x location) isSame(another location) bool {
	return x.getI() == another.getI() && x.getJ() == another.getJ()
}

func (x location) getI() int {
	return int(uint8(x >> 32))
}

func (x location) getJ() int {
	return int(uint32(x))
}

type ComponentManager struct {
	// comp is vertical (Compotype idexes)
	// ent is hori (safe set + map)
	container [][]unsafe.Pointer

	// entidx:compidx idx
	navigator  map[Component]location
	navigator2 map[location]Component
}

// OnActivate implements [ISystemComponent].
func (x *ComponentManager) OnActivate() {
}

// OnPreLoad implements [ISystemComponent].
func (x *ComponentManager) OnPreLoad() {
}

// OnUnLoad implements [ISystemComponent].
func (x *ComponentManager) OnUnLoad() {
}

var componentManager = ComponentManager{
	container:  make([][]unsafe.Pointer, 0, 8),
	navigator:  map[Component]location{},
	navigator2: map[location]Component{},
}

func Components() *ComponentManager {
	return &componentManager
}

func (x *ComponentManager) Init() errnov1.Code {
	return errnov1.OK
}

func (x *ComponentManager) OnAllLoaded() {}

func (x *ComponentManager) New(e Entity, t Compotype, buff *Component) errnov1.Code {
	if buff == nil {
		return errnov1.EINVAL
	}

	y := ComponentFor(e, t)

	if x.Contains(y) {
		return errnov1.ECALL
	}

	index := t.Int()
	if len(x.container) <= index {
		for i := len(x.container); i < (index + 1); i++ {
			x.container = append(x.container, []unsafe.Pointer{})
		}
	}

	loc := locationOf(t.Int(), len(x.container[index]))

	// allocate new comp
	x.container[index] = append(x.container[index], Compotypes().Allocate(t))

	// store location
	x.navigator[y] = loc
	x.navigator2[loc] = y

	Entities().WithComponent(e, t)

	*buff = y

	return errnov1.OK
}

func (x *ComponentManager) Get(v Component) unsafe.Pointer {
	loc := locationOf(0, 0)

	if errnov1.SUCCESS(x.navigate(v, &loc)) {
		return x.Components()[loc.getI()][loc.getJ()]
	}

	return nil
}

func (x *ComponentManager) Destroy(v Component) {
	loc := locationOf(0, 0)

	if errnov1.FAIL(x.navigate(v, &loc)) {
		return
	}

	delete(x.navigator, v)

	if Entities().IsAlive(v.Entity()) {
		Entities().WithoutComponent(v.Entity(), v.Compotype())
	}

	locLast := locationOf(loc.getI(), len(x.container[loc.getI()])-1)

	if last, ok := x.navigator2[locLast]; ok && !loc.isSame(locLast) {
		x.container[loc.getI()][loc.getJ()] = x.container[locLast.getI()][locLast.getJ()]
		x.navigator[last] = loc
		x.navigator2[loc] = last
	}

	x.container[locLast.getI()] = x.container[locLast.getI()][:locLast.getJ()]
	delete(x.navigator2, locLast)
}

func (x *ComponentManager) Components() [][]unsafe.Pointer {
	return x.container
}

func (x *ComponentManager) Iterator(v Compotype) *Iterator[unsafe.Pointer] {
	if len(x.container) <= v.Int() {
		return nil
	}

	buff := make([]unsafe.Pointer, len(x.container[v.Int()]))

	copy(buff, x.container[v.Int()])

	return NewIterator(buff)
}

func (x *ComponentManager) IterateB(v Compotype, fn func(Entity, unsafe.Pointer)) {
	if len(x.container) <= v.Int() {
		return
	}

	i := v.Int()
	loc := locationOf(0, 0)
	component := Component(0)
	entities := Entities()

	for j, y := range x.container[v.Int()] {
		if loc = locationOf(i, j); errnov1.SUCCESS(x.navigate2(loc, &component)) {
			if entities.IsAlive(component.Entity()) {
				fn(component.Entity(), y)
			}
		}
	}
}

func (x *ComponentManager) IterateC(e Entity, handlerFn func(t Compotype, ptr unsafe.Pointer)) {
	entities := Entities()

	for i := range x.container {
		if entities.Has(e, Compotype(i)) {
			handlerFn(Compotype(i), x.Get(ComponentFor(e, Compotype(i))))
		}
	}
}

func (x *ComponentManager) FindFirstEntity(t Compotype, buff *Entity, condFn func(Entity, unsafe.Pointer) bool) bool {
	if len(x.container) <= t.Int() {
		return false
	}

	entities := Entities()

	return entities.Iterator().First(buff, func(ent Entity) bool {
		if entities.IsAlive(ent) {
			if component := ComponentFor(ent, t); x.Contains(component) {
				return condFn(ent, x.Get(component))
			}
		}

		return false
	})
}

func (x *ComponentManager) navigate(v Component, buff *location) errnov1.Code {
	if v, ok := x.navigator[v]; ok {
		*buff = v
		return errnov1.OK
	}

	return errnov1.EINVAL
}

func (x *ComponentManager) navigate2(i location, buff *Component) errnov1.Code {
	if v, ok := x.navigator2[i]; ok {
		*buff = v
		return errnov1.OK
	}

	return errnov1.EINVAL
}

func (x *ComponentManager) navigateA(e Entity, t Compotype, buff *location) errnov1.Code {
	return x.navigate(ComponentFor(e, t), buff)
}

func (x *ComponentManager) Contains(v Component) bool {
	_, ok := x.navigator[v]
	return ok
}

func (x *ComponentManager) ContainsA(e Entity, t Compotype) bool {
	return x.Contains(ComponentFor(e, t))
}
