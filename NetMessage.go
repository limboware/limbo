package limbov1

import (
	"encoding/binary"
	"hash/crc32"
	"unsafe"
)

const (
	NETMSG_MAGIC_VALUE           = uint32(0xBEEEAAA3)
	NETMSG_VERSION               = 1
	NETMSG_HEADER_SIZE           = 64
	NETMSG_PAYLOAD_MAX_LENGTH    = 1 << 20
	NETMSG_GCM_SIZE              = 16
	NETMSG_RESERVE_SIZE          = 16
	NETMSG_PAYLOAD_WITH_GCM_SIZE = NETMSG_PAYLOAD_MAX_LENGTH + NETMSG_GCM_SIZE
	NETMSG_SIZE                  = NETMSG_HEADER_SIZE + NETMSG_PAYLOAD_MAX_LENGTH + NETMSG_GCM_SIZE // + gcm tag (16)
)

// NetMessageFlag
type NetMessageFlag byte

const (
	FlagNone       NetMessageFlag = 0x00
	FlagCompressed NetMessageFlag = 0x01
	FlagEncrypted  NetMessageFlag = 0x02
	FlagSigned     NetMessageFlag = 0x04
)

func (x NetMessageFlag) Is(v NetMessageFlag) bool {
	return x == v
}

func (x NetMessageFlag) Contains(v NetMessageFlag) bool {
	return (x & v) != 0
}

func (x NetMessageFlag) With(v NetMessageFlag) NetMessageFlag {
	return x | v
}

func (x NetMessageFlag) Without(v NetMessageFlag) NetMessageFlag {
	return x & ^v
}

func (x NetMessageFlag) Closer() byte {
	return (byte)(x)
}

// NetMessageType
type NetMessageType uint32

func (x NetMessageType) Closer() uint32 {
	return (uint32)(x)
}

// NetMessage
type NetMessage byte

func NetMessageOf(v []byte) *NetMessage {
	return (*NetMessage)(&v[0])
}

func (x *NetMessage) Magic() uint32 {
	return binary.BigEndian.Uint32(unsafe.Slice((*byte)(x), 4))
}

func (x *NetMessage) Version() uint16 {
	return binary.BigEndian.Uint16(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 4)), 2))
}

func (x *NetMessage) TotalLen() uint32 {
	return binary.BigEndian.Uint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 6)), 4))
}

func (x *NetMessage) Sequence() uint32 {
	return binary.BigEndian.Uint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 10)), 4))
}

func (x *NetMessage) Timestamp() uint64 {
	return binary.BigEndian.Uint64(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 14)), 8))
}

func (x *NetMessage) SessionID() uint64 {
	return binary.BigEndian.Uint64(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 22)), 8))
}

func (x *NetMessage) Type() NetMessageType {
	return (NetMessageType)(binary.BigEndian.Uint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 30)), 4)))
}

func (x *NetMessage) PayloadLen() uint32 {
	return binary.BigEndian.Uint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 34)), 4))
}

func (x *NetMessage) PayloadChecksum() uint32 {
	return binary.BigEndian.Uint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 38)), 4))
}

func (x *NetMessage) PayloadRealChecksum() uint32 {
	data := x.BytesB(NETMSG_HEADER_SIZE, int(x.PayloadLen()))
	return crc32.ChecksumIEEE(data)
}

func (x *NetMessage) HeaderChecksum() uint32 {
	return binary.BigEndian.Uint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 42)), 4))
}

func (x *NetMessage) HeaderRealChecksum() uint32 {
	data := x.Bytes(NETMSG_HEADER_SIZE)
	return crc32.ChecksumIEEE(append(data[0:42], data[46:]...))
}

func (x *NetMessage) Flags() NetMessageFlag {
	return NetMessageFlag(*(*byte)(unsafe.Add(unsafe.Pointer(x), 46)))
}

func (x *NetMessage) KeyID() uint8 {
	return *(*byte)(unsafe.Add(unsafe.Pointer(x), 47))
}

func (x *NetMessage) Reserved() *byte {
	return (*byte)(unsafe.Add(unsafe.Pointer(x), 48))
}

func (x *NetMessage) Header() *byte {
	return (*byte)(x)
}

func (x *NetMessage) Closer() *byte {
	return x.Header()
}

func (x *NetMessage) Payload() *byte {
	return (*byte)(unsafe.Add(unsafe.Pointer(x), 64))
}

func (x *NetMessage) Bytes(len int) []byte {
	return unsafe.Slice((*byte)(x), len)
}

func (x *NetMessage) BytesB(offset, len int) []byte {
	return unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), offset)), len)
}

func (x *NetMessage) Builder() *NetMessageBuilder {
	return (*NetMessageBuilder)(x)
}

// NetMessageBuilder
type NetMessageBuilder byte

func (x *NetMessageBuilder) putMagic() *NetMessageBuilder {
	binary.BigEndian.PutUint32(unsafe.Slice((*byte)(x), 4), NETMSG_MAGIC_VALUE)
	return x
}

func (x *NetMessageBuilder) putVersion() *NetMessageBuilder {
	binary.BigEndian.PutUint16(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 4)), 2), NETMSG_VERSION)
	return x
}

func (x *NetMessageBuilder) totalLen() *NetMessageBuilder {
	total := NETMSG_HEADER_SIZE + (*NetMessage)(x).PayloadLen()

	if NetMessageFlag((*NetMessage)(x).Flags()).Contains(FlagEncrypted) {
		total += NETMSG_GCM_SIZE
	}

	binary.BigEndian.PutUint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 6)), 4), total)
	return x
}

func (x *NetMessageBuilder) PutSequence(v uint32) *NetMessageBuilder {
	binary.BigEndian.PutUint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 10)), 4), v)
	return x
}

func (x *NetMessageBuilder) PutTimestamp(v uint64) *NetMessageBuilder {
	binary.BigEndian.PutUint64(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 14)), 8), v)
	return x
}

func (x *NetMessageBuilder) PutSessionID(v uint64) *NetMessageBuilder {
	binary.BigEndian.PutUint64(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 22)), 8), v)
	return x
}

func (x *NetMessageBuilder) PutType(v uint32) *NetMessageBuilder {
	binary.BigEndian.PutUint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 30)), 4), v)
	return x
}

func (x *NetMessageBuilder) PutPayloadLen(v uint32) *NetMessageBuilder {
	binary.BigEndian.PutUint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 34)), 4), v)
	return x
}

func (x *NetMessageBuilder) PutPayloadChecksum(v uint32) *NetMessageBuilder {
	binary.BigEndian.PutUint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 38)), 4), v)
	return x
}

func (x *NetMessageBuilder) headerChecksum() *NetMessageBuilder {
	binary.BigEndian.PutUint32(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 42)), 4), (*NetMessage)(x).HeaderRealChecksum())
	return x
}

func (x *NetMessageBuilder) Build() {
	x.putMagic().putVersion().totalLen().headerChecksum()
}

func (x *NetMessageBuilder) PutFlags(v NetMessageFlag) *NetMessageBuilder {
	*(*byte)(unsafe.Add(unsafe.Pointer(x), 46)) = v.Closer()
	return x
}

func (x *NetMessageBuilder) PutKeyID(v uint8) *NetMessageBuilder {
	*(*byte)(unsafe.Add(unsafe.Pointer(x), 47)) = v
	return x
}

func (x *NetMessageBuilder) ReserveBytes(offset int, v *byte, len int) *NetMessageBuilder {
	if len > NETMSG_RESERVE_SIZE {
		len = NETMSG_RESERVE_SIZE
	}

	if (offset + len) > NETMSG_RESERVE_SIZE {
		return x
	}

	copy(unsafe.Slice((*byte)(unsafe.Add(unsafe.Pointer(x), 48+offset)), NETMSG_RESERVE_SIZE-offset), unsafe.Slice(v, len))

	return x
}
