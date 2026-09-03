package internal

import (
	"net"
)

// poll.netFD
type NetFD struct {
	Pfd FD

	// immutable until Close
	family      int
	sotype      int
	isConnected bool // handshake completed or use of association with peer
	net         string
	laddr       net.Addr
	raddr       net.Addr
}
