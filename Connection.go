package limbov1

type ConnectionId uint64

func CreateConnectionId() ConnectionId {
	return ConnectionId(0)
}

func (x ConnectionId) Closer() uint64 {
	return uint64(x)
}

func (x ConnectionId) SetId(v uint64) ConnectionId {
	return ConnectionId(x.Closer() | uint64(v)<<16)
}

func (x ConnectionId) SetProtocol(v uint8) ConnectionId {
	return ConnectionId(x.Closer() | uint64(v)<<8)
}

func (x ConnectionId) SetState(v uint8) ConnectionId {
	return ConnectionId(x.Closer() | uint64(v))
}

func (x ConnectionId) Id() uint64 {
	return x.Closer() >> 16
}

func (x ConnectionId) Protocol() uint8 {
	return uint8(x.Closer() >> 8 & 0xFF)
}

func (x ConnectionId) State() uint8 {
	return uint8(x.Closer() & 0xFF)
}
