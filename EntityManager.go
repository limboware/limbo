package limbov1

import (
	errnov1 "github.com/rejchev/errno"
)

type EntityManager struct {
	ents []Entity
	free []uint32
	gens []uint16
	mask []ComponentMask

	entmap map[Entity]int
}

var entityManager = EntityManager{}

func Entities() *EntityManager {
	return &entityManager
}

func (x *EntityManager) Init() errnov1.Code {
	x.ents = make([]Entity, 0, 8)
	x.gens = make([]uint16, 0, 8)
	x.free = make([]uint32, 0, 8)
	x.mask = make([]ComponentMask, 0, 8)
	x.entmap = map[Entity]int{}

	return errnov1.OK
}

func (x *EntityManager) Cleanup() {
	x.entmap = map[Entity]int{}
	x.ents = x.ents[:]
	x.free = x.free[:]
	x.gens = x.gens[:]
	x.mask = x.mask[:]
}

func (x *EntityManager) Create() Entity {
	idx := uint32(0)

	if len(x.free) > 0 {
		idx = x.free[len(x.free)-1]
		x.free = x.free[:len(x.free)-1]
	} else {
		idx = uint32(len(x.gens))
		x.gens = append(x.gens, 0)
	}

	if len(x.mask) <= int(idx) {
		x.mask = append(x.mask, CreateComponentMask())
	}

	ent := CreateEntity(idx, x.gens[idx])

	x.ents = append(x.ents, ent)
	x.entmap[ent] = len(x.ents) - 1

	Events().Publish("entity.created", ent)

	return ent
}

// DestroyEntitySystem -> ent.Destroy -> sync Publish("entity.destroy")
func (x *EntityManager) Destroy(v Entity) {
	innerIdx := x.index(v)

	if innerIdx == -1 {
		return
	}

	Events().Publish("entity.destroy", v)

	idx := v.Id()
	x.gens[idx]++
	x.mask[idx].Reset()
	x.free = append(x.free, idx)

	entsLen := len(x.ents)

	if entsLen > (innerIdx + 1) {
		x.ents[innerIdx] = x.ents[entsLen-1]
		x.mask[innerIdx] = x.mask[entsLen-1]
		x.entmap[x.ents[entsLen-1]] = innerIdx
	}

	delete(x.entmap, v)
	x.ents = x.ents[:entsLen-1]
	x.mask = x.mask[:entsLen-1]

	Events().Publish("entity.destroyed", v)
}

func (x *EntityManager) Mask(v Entity) ComponentMask {
	return x.mask[v.Id()]
}

func (x *EntityManager) Iterator() *Iterator[Entity] {
	buff := make([]Entity, len(x.ents))

	copy(buff, x.ents)

	return NewIterator(buff)
}

func (x *EntityManager) WithComponent(e Entity, v Compotype) {
	x.Mask(e).With(v)
}

func (x *EntityManager) WithComponents(e Entity, v ...Compotype) {
	x.Mask(e).WithB(v...)
}

func (x *EntityManager) WithoutComponent(e Entity, v Compotype) {
	x.Mask(e).Without(v)
}

func (x *EntityManager) SetMask(e Entity, v ComponentMask) {
	x.mask[e.Id()] = v
}

func (x *EntityManager) Masks() []ComponentMask {
	return x.mask
}

func (x *EntityManager) Count() int {
	return len(x.entmap)
}

func (x *EntityManager) index(v Entity) int {
	if v, ok := x.entmap[v]; ok {
		return v
	}

	return -1
}

func (x *EntityManager) HasComponentB(e Entity, t Compotype) bool {
	return x.HasComponent(e, t)
}

func (x *EntityManager) HasComponent(e Entity, t Compotype) bool {
	return x.Mask(e).Has(t)
}

func (x *EntityManager) Has(e Entity, t Compotype) bool {
	return x.HasComponentB(e, t)
}

func (x *EntityManager) IsAlive(v Entity) bool {
	if uint32(len(x.gens)) <= v.Id() {
		return false
	}

	return x.gens[v.Id()] == v.Gen() && x.index(v) != -1
}
