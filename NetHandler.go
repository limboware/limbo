package limbov1

type NetworkHandler [5]byte

func MakeNetworkHandler(nett NetConnType, msgt NetMessageType) NetworkHandler {
	return NetworkHandler{
		byte(msgt),
		byte(msgt >> 8),
		byte(msgt >> 16),
		byte(msgt >> 24),
		nett.Closer(),
	}
}

func (x NetworkHandler) ConnType() NetConnType {
	return (NetConnType)(x[4])
}

func (x NetworkHandler) MsgType() NetMessageType {
	return NetMessageType(uint32(x[3])<<24 | uint32(x[2])<<16 | uint32(x[1])<<8 | uint32(x[0]))
}

func (x NetworkHandler) Closer() [5]byte {
	return ([5]byte)(x)
}
