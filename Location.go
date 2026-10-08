package limbov1

type location uint64

func LocationOf(i int, j int) location {
	return location(uint64(uint32(i))<<32 | uint64(uint32(j)))
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