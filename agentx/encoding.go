package agentx

import (
	"errors"

	"github.com/alexispb/mygosnmp/internal"
)

var ErrEncoding = errors.New("encoding error")

// byteOrder specifies how to encode/decode integer types int16,
// int32, uint32, and uint64.
// Note. The bitwise operators are endian neutral (e.g. the >>
// operator always shifts the bits towards the least significant
// digit, so 0x1234 >> 8 always results with 0x12).
type byteOrder interface {
	String() string

	appendInt16([]byte, int16) []byte
	appendInt32([]byte, int32) []byte
	appendUint32([]byte, uint32) []byte
	appendUint64([]byte, uint64) []byte

	parseInt16([]byte) (int16, []byte)
	parseInt32([]byte) (int32, []byte)
	parseUint32([]byte) (uint32, []byte)
	parseUint64([]byte) (uint64, []byte)
}

// networkByteOrder is the big-endian implementation of byteOrder
// (used if the FlagNetworkByteOrder is set).
type networkByteOrder struct{}

func (networkByteOrder) String() string {
	return "direct byte order"
}

func (networkByteOrder) appendInt16(data []byte, val int16) []byte {
	return append(data, byte(val>>8), byte(val))
}
func (networkByteOrder) parseInt16(data []byte) (int16, []byte) {
	return int16(data[0])<<8 | int16(data[1]), data[2:]
}
func (networkByteOrder) appendInt32(data []byte, val int32) []byte {
	return append(data, byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}
func (networkByteOrder) parseInt32(data []byte) (int32, []byte) {
	return int32(data[0])<<24 | int32(data[1])<<16 | int32(data[2])<<8 | int32(data[3]), data[4:]
}
func (networkByteOrder) appendUint32(data []byte, val uint32) []byte {
	return append(data, byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}
func (networkByteOrder) parseUint32(data []byte) (uint32, []byte) {
	return uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3]), data[4:]
}
func (networkByteOrder) appendUint64(data []byte, val uint64) []byte {
	return append(data, byte(val>>56), byte(val>>48), byte(val>>40), byte(val>>32),
		byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}
func (networkByteOrder) parseUint64(data []byte) (uint64, []byte) {
	return uint64(data[0])<<56 | uint64(data[1])<<48 | uint64(data[2])<<40 | uint64(data[3])<<32 |
		uint64(data[4])<<24 | uint64(data[5])<<16 | uint64(data[6])<<8 | uint64(data[7]), data[8:]
}

// inverseByteOrder is the little-endian implementation of byteOrder
// (used if the FlagNetworkByteOrder is not set).
type inverseByteOrder struct{}

func (inverseByteOrder) String() string {
	return "inverse byte order"
}

func (inverseByteOrder) appendInt16(data []byte, val int16) []byte {
	return append(data, byte(val), byte(val>>8))
}
func (inverseByteOrder) parseInt16(data []byte) (int16, []byte) {
	return int16(data[0]) | int16(data[1])<<8, data[2:]
}
func (inverseByteOrder) appendInt32(data []byte, val int32) []byte {
	return append(data, byte(val), byte(val>>8), byte(val>>16), byte(val>>24))
}
func (inverseByteOrder) parseInt32(data []byte) (int32, []byte) {
	return int32(data[0]) | int32(data[1])<<8 | int32(data[2])<<16 | int32(data[3])<<24, data[4:]
}
func (inverseByteOrder) appendUint32(data []byte, val uint32) []byte {
	return append(data, byte(val), byte(val>>8), byte(val>>16), byte(val>>24))
}
func (inverseByteOrder) parseUint32(data []byte) (uint32, []byte) {
	return uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24, data[4:]
}
func (inverseByteOrder) appendUint64(data []byte, val uint64) []byte {
	return append(data, byte(val), byte(val>>8), byte(val>>16), byte(val>>24),
		byte(val>>32), byte(val>>40), byte(val>>48), byte(val>>56))
}
func (inverseByteOrder) parseUint64(data []byte) (uint64, []byte) {
	return uint64(data[0]) | uint64(data[1])<<8 | uint64(data[2])<<16 | uint64(data[3])<<24 |
		uint64(data[4])<<32 | uint64(data[5])<<40 | uint64(data[6])<<48 | uint64(data[7])<<56, data[8:]
}

// encoder implements encoding of composite types such as
// OctetString, ObjectId, IpAddress, Opaque, etc.
type encoder struct {
	byteOrder
}

// decoder implements decoding of composite types such as
// OctetString, ObjectId, IpAddress, Opaque, etc.
type decoder struct {
	byteOrder
}

var mod2pad [4]int = [4]int{0, 3, 2, 1}

func paddingSize(n int) int {
	return mod2pad[n&0x03]
}

func (encoder) appendPaddingBytes(data []byte) []byte {
	switch paddingSize(len(data)) {
	case 1:
		return append(data, 0)
	case 2:
		return append(data, 0, 0)
	case 3:
		return append(data, 0, 0, 0)
	default:
		return data
	}
}

func (d decoder) parsePaddingBytes(data []byte, size int) []byte {
	size = paddingSize(size)
	for i := 0; i < size; i++ {
		if data[i] != 0 {
			panic(ErrEncoding)
		}
	}
	return data[size:]
}

func octetStringEncodingSize(val string) int {
	return 4 + len(val) + paddingSize(len(val))
}

func (e encoder) appendOctetString(data []byte, val string) []byte {
	data = e.appendInt32(data, int32(len(val)))
	data = append(data, internal.StringAsSlice(val)...)
	return e.appendPaddingBytes(data)
}

func (d decoder) parseOctetString(data []byte) (val string, next []byte) {
	size, next := d.parseInt32(data)
	if size < 0 {
		panic(ErrEncoding)
	}
	val, next = string(next[:int(size)]), next[int(size):]
	next = d.parsePaddingBytes(next, int(size))
	return
}

func objectIdEncodingSize(id []uint32) (size int) {
	size = 4
	if len(id) > 0 {
		nsubids := len(id)
		if len(id) >= 5 && id[0] == 1 && id[1] == 3 && id[2] == 6 && id[3] == 1 && id[4] <= 256 {
			nsubids -= 5
		}
		size += 4 * nsubids
	}
	return
}

func (e encoder) appendObjectId(data []byte, val []uint32, include byte) []byte {
	if len(val) == 0 {
		return append(data, 0, 0, 0, 0)
	}
	nsubids := len(val)
	prefix := byte(0)
	if nsubids >= 5 && val[0] == 1 && val[1] == 3 && val[2] == 6 && val[3] == 1 && val[4] <= 256 {
		nsubids -= 5
		prefix = byte(val[4])
		val = val[5:]
	}
	data = append(data, byte(nsubids), prefix, include, 0)
	for _, subid := range val {
		data = e.appendUint32(data, subid)
	}
	return data
}

func (d decoder) parseObjectId(data []byte) (val []uint32, include byte, next []byte) {
	if data[3] != 0 {
		panic(ErrEncoding)
	}
	nsubids, prefix, include := int(data[0]), uint32(data[1]), data[2]
	next = data[4:]
	if nsubids == 0 && prefix == 0 {
		return
	}
	if prefix == 0 {
		val = make([]uint32, 0, nsubids)
	} else {
		val = make([]uint32, 0, nsubids+5)
		val = append(val, 1, 3, 6, 1, prefix)
	}
	for i := 0; i < nsubids; i++ {
		var subid uint32
		subid, next = d.parseUint32(next)
		val = append(val, subid)
	}
	return
}

func ipAddressEncodingSize(val [4]byte) int {
	return 8
}

func (e encoder) appendIpAddress(data []byte, val [4]byte) []byte {
	data = e.appendInt32(data, int32(4))
	return append(data, val[:]...)
}

func (d decoder) parseIpAddress(data []byte) (val [4]byte, next []byte) {
	size, next := d.parseInt32(data)
	if size != 4 {
		panic(ErrEncoding)
	}
	copy(val[:], next[:4])
	next = next[4:]
	return
}

func opaqueEncodingSize(val []byte) int {
	return 4 + len(val) + paddingSize(len(val))
}

func (e encoder) appendOpaque(data []byte, val []byte) []byte {
	data = e.appendInt32(data, int32(len(val)))
	data = append(data, val...)
	return e.appendPaddingBytes(data)
}

func (d decoder) parseOpaque(data []byte) (val []byte, next []byte) {
	size, next := d.parseInt32(data)
	if size < 0 {
		panic(ErrEncoding)
	}
	val = make([]byte, int(size))
	copy(val, next[:int(size)])
	next = next[int(size):]
	next = d.parsePaddingBytes(next, int(size))
	return
}
