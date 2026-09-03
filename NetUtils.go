package limbov1

import (
	"net"
	"unsafe"

	"github.com/limboware/limbo/internal"
)

func GetDescriptor(v net.Conn) uintptr {
	if ptr := (*internal.NetFD)(*(*unsafe.Pointer)(*(*unsafe.Pointer)(unsafe.Pointer(&v)))); ptr != nil {
		return uintptr(ptr.Pfd.Sysfd)
	}

	return 0
}
