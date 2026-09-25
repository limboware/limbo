package sys

import "syscall"

// poll.netFD
type FD struct {
	// Lock sysfd and serialize access to Read and Write methods.
	fdmu [16]byte

	// System file descriptor. Immutable until Close.
	Sysfd int

	// Platform dependent state of the file descriptor.
	SysFile [8]byte

	// I/O poller.
	pd byte

	// Semaphore signaled when file is closed.
	csema uint32

	// Non-zero if this file has been set to blocking mode.
	isBlocking uint32

	// Whether this is a streaming descriptor, as opposed to a
	// packet-based descriptor like a UDP socket. Immutable.
	IsStream bool

	// Whether a zero byte read indicates EOF. This is false for a
	// message based socket connection.
	ZeroReadIsEOF bool

	// Whether this is a file rather than a network socket.
	isFile bool
}

func SetNonblock(fd uintptr, nonb bool) error {
	return syscall.SetNonblock(int(fd), nonb)
}

func Read(fd uintptr, buff []byte) (int, error) {
	return syscall.Read(int(fd), buff)
}

func Write(fd uintptr, buff []byte) (int, error) {
	return syscall.Write(int(fd), buff)
}
