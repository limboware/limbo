package limbov1

import (
	"net"

	errnov1 "github.com/rejchev/errno"
)

const ConnectionsMax = uint64(0xFFFF_FFFF_FFFF_0000)

var manager NetworkManager = NetworkManager{
	conns: make([]net.Conn, 0, 8),
	gens:  make([]uint8, 0, 8),
	free:  make([]uint64, 0, 8),

	count: 0,
}

type NetworkManager struct {
	conns []net.Conn
	gens  []uint8
	free  []uint64

	count int
}

func Networks() *NetworkManager {
	return &manager
}

// Connection implements [IConnectionManager].
func (x *NetworkManager) Connection(v ConnectionId) net.Conn {
	if x.gens[v.Id()] != v.State() {
		return nil
	}

	return x.conns[v.Id()]
}

func (x *NetworkManager) OnAllLoaded() {}

// Destroy implements [IConnectionManager].
func (x *NetworkManager) Destroy(v ConnectionId) {

	if !x.IsAlive(v) {
		return
	}

	Events().Publish("connection.close", v)

	idx := v.Id()

	x.gens[idx]++
	if x.gens[idx] >= 128 {
		x.gens[idx] = 0
	}

	x.free = append(x.free, idx)

	if con := x.conns[idx]; con != nil {
		_ = con.Close()
	}

	x.conns[idx] = nil
	x.count--

	Events().Publish("connection.closed", v)
}

func (x *NetworkManager) Protocol(v ConnectionId) uint8 {
	return v.Protocol()
}

func (x *NetworkManager) Addr(v ConnectionId) net.Addr {
	if conn := x.Connection(v); conn != nil {
		return conn.RemoteAddr()
	}

	return nil
}

func (x *NetworkManager) TCPAddr(v ConnectionId) *net.TCPAddr {
	if v.Protocol() == 1 {
		if addr := x.Addr(v); addr != nil {
			if tcp, ok := addr.(*net.TCPAddr); ok {
				return tcp
			}
		}
	}

	return nil
}

func (x NetworkManager) IsAlive(v ConnectionId) bool {
	return x.gens[v.Id()] == v.State()
}

func (x *NetworkManager) Init() errnov1.Code {
	return errnov1.OK
}

// Count implements [IConnectionManager].
func (x *NetworkManager) Count() int {
	return x.count
}

func (x *NetworkManager) Connect(con net.Conn, buff *ConnectionId) errnov1.Code {
	return x.ConnectB(con, 1, buff)
}

func (x *NetworkManager) ConnectB(conn net.Conn, proto uint8, buff *ConnectionId) errnov1.Code {
	idx := uint64(0)

	if len(x.free) > 0 {
		idx = x.free[len(x.free)-1]
		x.free = x.free[:len(x.free)-1]
	} else {
		if idx = uint64(len(x.gens)); idx >= ConnectionsMax {
			return errnov1.EINVAL
		}

		x.gens = append(x.gens, 0)
	}

	if idx >= uint64(len(x.conns)) {
		x.conns = append(x.conns, conn)
	} else {
		if x.conns[idx] != nil {
			_ = x.conns[idx].Close()
		}

		x.conns[idx] = conn
	}

	x.count++

	*buff = CreateConnectionId().SetId(idx).SetProtocol(proto).SetState(x.gens[idx])

	Events().PublishAsync("connection.new", *buff)

	return errnov1.OK
}
