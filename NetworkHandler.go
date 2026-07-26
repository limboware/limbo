package limbov1

type NetworkHandler [5]byte

func NetworkHandlerFor(proto ConnectionProtocol, typ NetMessageType) NetworkHandler {
	return NetworkHandler{
		byte(typ),
		byte(typ >> 8),
		byte(typ >> 16),
		byte(typ >> 24),
		proto.Closer(),
	}
}

func (x NetworkHandler) Proto() ConnectionProtocol {
	return (ConnectionProtocol)(x[4])
}

func (x NetworkHandler) Type() NetMessageType {
	return NetMessageType(uint32(x[3])<<24 | uint32(x[2])<<16 | uint32(x[1])<<8 | uint32(x[0]))
}

func (x NetworkHandler) Closer() [5]byte {
	return ([5]byte)(x)
}
