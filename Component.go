package limbov1

type Component uint64

// FFFF__FFFF_FFFF_FFFF
// type_ent
func ComponentFor(e Entity, t Compotype) Component {
	return Component(uint64(t)<<48 | e.Closer())
}

func (x Component) Closer() uint64 {
	return (uint64)(x)
}

func (x Component) Compotype() Compotype {
	return Compotype(x >> 48)
}

func (x Component) Entity() Entity {
	return Entity(x & 0xFFFF_FFFF_FFFF)
}
