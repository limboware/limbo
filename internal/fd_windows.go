//go:build windows
package internal

import (
	"runtime"
	"sync/atomic"
	"syscall"
)

// poll.fd_windows
type FD struct {
	// Lock sysfd and serialize access to Read and Write methods.
	fdmu [16]byte

	// System file descriptor. Immutable until Close.
	Sysfd syscall.Handle

	// I/O poller.
	pd byte

	// The file offset for the next read or write.
	// Overlapped IO operations don't use the real file pointer,
	// so we need to keep track of the offset ourselves.
	offset int64

	// For console I/O.
	lastbits       []byte   // first few bytes of the last incomplete rune in last write
	readuint16     []uint16 // buffer to hold uint16s obtained with ReadConsole
	readbyte       []byte   // buffer to hold decoding of readuint16 from utf16 to utf8
	readbyteOffset int      // readbyte[readOffset:] is yet to be consumed with file.Read

	// Semaphore signaled when file is closed.
	csema uint32

	skipSyncNotif bool

	// Whether this is a streaming descriptor, as opposed to a
	// packet-based descriptor like a UDP socket.
	IsStream bool

	// Whether a zero byte read indicates EOF. This is false for a
	// message based socket connection.
	ZeroReadIsEOF bool

	// Whether the handle is owned by os.File.
	isFile bool

	// The kind of this file.
	kind byte

	// Whether FILE_FLAG_OVERLAPPED was not set when opening the file.
	isBlocking bool

	disassociated atomic.Bool

	// readPinner and writePinner are automatically unpinned
	// before execIO returns.
	readPinner  runtime.Pinner
	writePinner runtime.Pinner
}
