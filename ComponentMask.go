package limbov1

import (
	"fmt"
)

type ComponentMask [32]byte

func CreateComponentMask() ComponentMask {
	return ComponentMask([32]byte{}) // compotype is uint8 -> max 256
}

// in bytes
func (x ComponentMask) Len() int {
	return len(x)
}

func (x ComponentMask) With(i Compotype) ComponentMask {
	x[x.bucketId(i)] |= (1 << x.bucketBitId(i))
	return x
}

// one bucket eq 8bits -> each bucket is byte in []byte
// compotype is idx of bit in byte seq eq []byte
func (x ComponentMask) bucketId(i Compotype) int {
	return i.Int() / 8
}

func (x ComponentMask) bucketBitId(i Compotype) int {
	return i.Int() % 8
}

func (x ComponentMask) WithB(v ...Compotype) ComponentMask {
	for _, i := range v {
		x.With(i)
	}

	return x
}

func (x ComponentMask) Get(id int) byte {
	return x[id]
}

func (x ComponentMask) Get64(id int) uint64 {
	v := uint64(0)

	if (len(x) - id) > 7 {
		v = uint64(x[id])
		v |= uint64(x[id+1]) << 8
		v |= uint64(x[id+2]) << 16
		v |= uint64(x[id+3]) << 24
		v |= uint64(x[id+4]) << 32
		v |= uint64(x[id+5]) << 40
		v |= uint64(x[id+6]) << 48
		v |= uint64(x[id+7]) << 56
	}

	return v
}

func (x ComponentMask) Set64(id int, v uint64) ComponentMask {
	if (len(x) - id) > 7 {
		x[id] = byte(v)
		x[id+1] = byte(v >> 8)
		x[id+2] = byte(v >> 16)
		x[id+3] = byte(v >> 24)
		x[id+4] = byte(v >> 32)
		x[id+5] = byte(v >> 40)
		x[id+6] = byte(v >> 48)
		x[id+7] = byte(v >> 56)
	}

	return x
}

func (x ComponentMask) With8(id int, v uint8) ComponentMask {
	x[id] |= v
	return x
}

func (x ComponentMask) With16(id int, v uint16) ComponentMask {
	if (len(x) - id) > 1 {
		x[id] |= byte(v)
		x[id+1] |= byte(v >> 8)
	}

	return x
}

func (x ComponentMask) With32(id int, v uint32) ComponentMask {
	if (len(x) - id) > 3 {
		x.With16(id, uint16(v))
		x.With16(id+2, uint16(v>>16))
	}

	return x
}

func (x ComponentMask) With64(id int, v uint64) ComponentMask {
	if (len(x) - id) > 7 {
		x.With32(id, uint32(v))
		x.With32(id+4, uint32(v>>32))
	}

	return x
}

func (x ComponentMask) Without(i Compotype) ComponentMask {
	x[x.bucketId(i)] &= ^(1 << x.bucketBitId(i))
	return x
}

func (x ComponentMask) Has(v Compotype) bool {
	return (x[x.bucketId(v)] & (1 << x.bucketBitId(v))) != 0
}

func (x ComponentMask) Reset() {
	x.Set64(0, 0)
	x.Set64(8, 0)
	x.Set64(16, 0)
	x.Set64(24, 0)
}

func (x ComponentMask) IsEmpty() bool {
	return x.Get64(0) == 0 && x.Get64(8) == 0 && x.Get64(16) == 0 && x.Get64(24) == 0
}

func (x ComponentMask) String() string {
	return fmt.Sprintf("%d %d %d %d", x.Get64(0), x.Get64(8), x.Get64(16), x.Get64(24))
}
