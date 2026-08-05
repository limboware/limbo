package limbov1

type NetConn uint64

func MakeNetConn() NetConn {
	return NetConn(0)
}

func (x NetConn) Closer() uint64 {
	return uint64(x)
}

func (x NetConn) SetId(v uint64) NetConn {
	return NetConn(x.Closer() | uint64(v)<<16)
}

func (x NetConn) SetType(v NetConnType) NetConn {
	return NetConn(x.Closer() | uint64(v)<<8)
}

func (x NetConn) SetGen(v uint8) NetConn {
	return NetConn(x.Closer() | uint64(v))
}

func (x NetConn) Id() uint64 {
	return x.Closer() >> 16
}

func (x NetConn) Type() NetConnType {
	return NetConnType(x.Closer() >> 8)
}

func (x NetConn) Gen() uint8 {
	return uint8(x.Closer())
}
