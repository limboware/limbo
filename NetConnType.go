package limbov1

import "math"

const MaxNetConnTypes = math.MaxUint8

type NetConnType uint8

func (x NetConnType) Closer() uint8 {
	return (uint8)(x)
}

func (x NetConnType) Int() int {
	return (int)(x)
}
