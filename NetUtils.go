package limbov1

import (
	"net"
	"unsafe"

	"github.com/limboware/limbo/sys"
)

func GetDescriptor(v net.Conn) uintptr {
	if ptr := (*sys.NetFD)(*(*unsafe.Pointer)(*(*unsafe.Pointer)(unsafe.Pointer(&v)))); ptr != nil {
		return uintptr(ptr.Pfd.Sysfd)
	}

	return 0
}
