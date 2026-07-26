package limbov1

type ConnectionProtocol uint8

const (
	UnknownConnProto ConnectionProtocol = 0
	TCPConnProto     ConnectionProtocol = 1
	WsConnProto      ConnectionProtocol = 1 << 1
	UDPConnProto     ConnectionProtocol = 1 << 2
)

func (x ConnectionProtocol) Closer() uint8 {
	return (uint8)(x)
}

func (x ConnectionProtocol) Int() int {
	return (int)(x)
}

func (x ConnectionProtocol) String() string {
	switch x {
	case TCPConnProto:
		return "tcp"

	case UDPConnProto:
		return "udp"

	case WsConnProto:
		return "ws"

	default:
		return ""
	}
}
