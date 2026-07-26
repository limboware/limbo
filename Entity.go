package limbov1

type Entity uint64

func CreateEntity(id uint32, gen uint16) Entity {
	return Entity(uint64(id)<<16 | uint64(gen))
}

func (x Entity) Closer() uint64 {
	return (uint64)(x)
}

func (x Entity) Id() uint32 {
	return uint32(x >> 16)
}

func (x Entity) Gen() uint16 {
	return uint16(x & 0xFFFF)
}
