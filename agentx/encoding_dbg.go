package agentx

import (
	"fmt"

	"github.com/alexispb/mygosnmp/generics"
	"github.com/alexispb/mygosnmp/hex"
	"github.com/alexispb/mygosnmp/internal"
	"github.com/alexispb/mygosnmp/ipa"
	"github.com/alexispb/mygosnmp/logger"
	"github.com/alexispb/mygosnmp/oid"
)

// encoderDbg implements methods which encode corresponding values
// with additional logging of encoding process.
type encoderDbg struct {
	byteOrder
	log logger.Log
}

// decoderDbg implements methods which decode corresponding values
// with additional logging of decoding process.
type decoderDbg struct {
	byteOrder
	data []byte
	log  logger.Log
}

func (d decoderDbg) nextChunk(startpos, size int) (chunk []byte, nextpos int) {
	n := generics.Min(len(d.data)-startpos, size)
	d.log.Write(hex.DumpSub("    ", d.data, startpos, startpos+n))
	if n < size {
		panic(fmt.Sprintf("not enough data: required %d bytes", size))
	}
	nextpos = startpos + size
	chunk = d.data[startpos:nextpos]
	return
}

func (e encoderDbg) appendInt16(data []byte, val int16) []byte {
	e.log.Writef("appending int16: %d 0x%04X", val, val)
	data = e.byteOrder.appendInt16(data, val)
	e.log.Write(hex.DumpTail("    ", data, 2))
	return data
}

func (d decoderDbg) parseInt16(startpos int) (val int16, nextpos int) {
	d.log.Write("parsing int16 data:")
	data, nextpos := d.nextChunk(startpos, 2)
	val, _ = d.byteOrder.parseInt16(data)
	d.log.Writef("parsed int16: %d", val)
	return
}

func (e encoderDbg) appendInt32(data []byte, val int32) []byte {
	e.log.Writef("appending int32: %d 0x%08X", val, val)
	data = e.byteOrder.appendInt32(data, val)
	e.log.Write(hex.DumpTail("    ", data, 4))
	return data
}

func (d decoderDbg) parseInt32(startpos int) (val int32, nextpos int) {
	d.log.Write("parsing int32 data:")
	data, nextpos := d.nextChunk(startpos, 4)
	val, _ = d.byteOrder.parseInt32(data)
	d.log.Writef("parsed int32: %d", val)
	return
}

func (e encoderDbg) appendUint32(data []byte, val uint32) []byte {
	e.log.Writef("appending uint32: %d 0x%08X", val, val)
	data = e.byteOrder.appendUint32(data, val)
	e.log.Write(hex.DumpTail("    ", data, 4))
	return data
}

func (d decoderDbg) parseUint32(startpos int) (val uint32, nextpos int) {
	d.log.Write("parsing uint32 data:")
	data, nextpos := d.nextChunk(startpos, 4)
	val, _ = d.byteOrder.parseUint32(data)
	d.log.Writef("parsed uint32: %d", val)
	return
}

func (e encoderDbg) appendUint64(data []byte, val uint64) []byte {
	e.log.Writef("appending uint64: %d 0x%016X", val, val)
	data = e.byteOrder.appendUint64(data, val)
	e.log.Write(hex.DumpTail("    ", data, 8))
	return data
}

func (d decoderDbg) parseUint64(startpos int) (val uint64, nextpos int) {
	d.log.Write("parsing uint64 data:")
	data, nextpos := d.nextChunk(startpos, 8)
	val, _ = d.byteOrder.parseUint64(data)
	d.log.Writef("parsed uint64: %d", val)
	return
}

func (e encoderDbg) appendPaddingBytes(data []byte) []byte {
	if size := paddingSize(len(data)); size > 0 {
		e.log.Writef("appending %d padding byte(s)", size)
		switch size {
		case 1:
			data = append(data, 0)
		case 2:
			data = append(data, 0, 0)
		case 3:
			data = append(data, 0, 0, 0)
		}
		e.log.Write(hex.DumpTail("    ", data, size))
	}
	return data
}

func (d decoderDbg) parsePaddingBytes(startpos, size int) (nextpos int) {
	if size = paddingSize(size); size == 0 {
		return startpos
	}
	d.log.Writef("parsing %d padding bytes:", size)
	data, nextpos := d.nextChunk(startpos, size)
	for _, b := range data {
		if b != 0 {
			panic("non-zero padding byte(s)")
		}
	}
	return nextpos
}

func (e encoderDbg) appendOctetString(data []byte, val string) []byte {
	e.log.Writef("appending OctetString: %s", val)
	startindex := len(data)

	e.log.Write("appending OctetString size")
	data = e.appendInt32(data, int32(len(val)))
	e.log.Write("appending OctetString bytes")
	data = append(data, internal.StringAsSlice(val)...)
	e.log.Write(hex.DumpTail("    ", data, len(val)))
	data = e.appendPaddingBytes(data)

	e.log.Writef("appended OctetString data\n%s",
		hex.DumpSub("    ", data, startindex, len(data)))
	return data
}

func (d decoderDbg) parseOctetString(startpos int) (val string, nextpos int) {
	d.log.Write("parsing OctetString data")
	d.log.Write("parsing OctetString size")
	size, nextpos := d.parseInt32(startpos)
	if size < 0 {
		panic("negative size")
	}
	d.log.Write("parsing OctetString bytes:")
	data, nextpos := d.nextChunk(nextpos, int(size))
	val = string(data)
	nextpos = d.parsePaddingBytes(nextpos, int(size))
	d.log.Writef("parsed OctetString: %s", val)
	return
}

func (e encoderDbg) appendObjectId(data []byte, val []uint32, include byte) []byte {
	e.log.Writef("appending ObjectId: %s, include: %d", oid.String(val), include)
	if len(val) == 0 {
		data = append(data, 0, 0, 0, 0)
		e.log.Write("appended null ObjectId")
		e.log.Write(hex.DumpTail("    ", data, 4))
		return data
	}

	startindex := len(data)
	nsubids := len(val)
	prefix := byte(0)

	if nsubids >= 5 && val[0] == 1 && val[1] == 3 && val[2] == 6 && val[3] == 1 && val[4] <= 256 {
		nsubids -= 5
		prefix = byte(val[4])
		val = val[5:]
	}

	e.log.Writef("appending ObjectId header {nsubids: %d, prefix: %d, include: %d, reserved byte: 0}:",
		nsubids, prefix, include)
	data = append(data, byte(nsubids), prefix, include, 0)
	e.log.Write(hex.DumpTail("    ", data, 4))

	e.log.Write("appending subids")
	for _, subid := range val {
		data = e.appendUint32(data, subid)
	}

	e.log.Writef("appended ObjectId data:\n%s",
		hex.DumpSub("    ", data, startindex, len(data)))
	return data
}

func (d *decoderDbg) parseObjectId(startpos int) (val []uint32, include byte, nextpos int) {
	d.log.Write("parsing ObjectId")
	d.log.Write("parsing ObjectId header:")
	data, nextpos := d.nextChunk(startpos, 4)
	nsubids, prefix, include := int(data[0]), uint32(data[1]), data[2]
	d.log.Writef("header: {nsubids: %d, prefix: %d, include: %d, reserved: %d}",
		nsubids, prefix, include, data[3])
	if data[3] != 0 {
		panic("non-zero reserved byte")
	}
	if nsubids == 0 && prefix == 0 {
		d.log.Write("parsed ObjectId: <null>")
		return
	}
	if prefix == 0 {
		val = make([]uint32, 0, nsubids)
	} else {
		val = make([]uint32, 0, nsubids+5)
		val = append(val, 1, 3, 6, 1, prefix)
		d.log.Writef("ObjectId starts with %s", oid.String(val))
	}
	for i := 0; i < nsubids; i++ {
		d.log.Write("parsing subid")
		var subid uint32
		subid, nextpos = d.parseUint32(nextpos)
		val = append(val, subid)
	}
	d.log.Writef("parsed ObjectId: %s", oid.String(val))
	return
}

func (e encoderDbg) appendIpAddress(data []byte, val [4]byte) []byte {
	e.log.Writef("appending IpAddress: %s", ipa.String(val))

	e.log.Write("appending IpAddress size")
	data = e.appendInt32(data, int32(4))
	e.log.Write("appending IpAddress bytes")
	data = append(data, val[:]...)
	e.log.Write(hex.DumpTail("    ", data, 4))

	e.log.Writef("appended IpAddress data\n%s",
		hex.DumpTail("    ", data, 8))
	return data
}

func (d decoderDbg) parseIpAddress(startpos int) (val [4]byte, nextpos int) {
	d.log.Write("parsing IpAddress")
	d.log.Write("parsing IpAddress size")
	size, nextpos := d.parseInt32(startpos)
	if size != 4 {
		panic("size != 4")
	}
	d.log.Write("parsing IpAddress bytes")
	data, nextpos := d.nextChunk(nextpos, 4)
	copy(val[:], data)
	d.log.Writef("parsed IpAddress: %s", ipa.String(val))
	return
}

func (e encoderDbg) appendOpaque(data []byte, val []byte) []byte {
	e.log.Writef("appending Opaque: %v", val)
	startindex := len(data)

	e.log.Write("appending Opaque size")
	data = e.appendInt32(data, int32(len(val)))
	e.log.Write("appending Opaque bytes")
	data = append(data, val...)
	e.log.Write(hex.DumpTail("    ", data, len(val)))
	data = e.appendPaddingBytes(data)

	e.log.Writef("appended Opaque data\n%s",
		hex.DumpSub("    ", data, startindex, len(data)))
	return data
}

func (d decoderDbg) parseOpaque(startpos int) (val []byte, nextpos int) {
	d.log.Write("parsing Opaque")
	d.log.Write("parsing Opaque size")
	size, nextpos := d.parseInt32(startpos)
	if size < 0 {
		panic("negative size")
	}
	d.log.Write("parsing Opaque bytes")
	data, nextpos := d.nextChunk(nextpos, int(size))
	val = make([]byte, int(size))
	copy(val, data)
	nextpos = d.parsePaddingBytes(nextpos, int(size))
	d.log.Writef("parsed Opaque: %v", val)
	return
}
