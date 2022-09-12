package agentx

import (
	"github.com/alexispb/mygosnmp/asn"
	"github.com/alexispb/mygosnmp/hex"
	"github.com/alexispb/mygosnmp/internal"
	"github.com/alexispb/mygosnmp/logger"
)

// PduHeaderSize is number of bytes required for encoding pdu header.
const PduHeaderSize = 20

type Pdu struct {
	// Header
	Tag           PduTag
	Flags         Flags
	SessionId     uint32
	TransactionId uint32
	PacketId      uint32
	PayloadSize   int32
	// Payload
	Context  string
	Params   PayloadParams
	Ranges   []SearchRange
	Varbinds []asn.Varbind
}

// String returns pdu string representation
func (pdu Pdu) String(depth int) string {
	return internal.StructString(pdu, depth)
}

// countPayloadSize returns number of bytes required for
// encoding pdu payload. This function panics if pdu.Tag
// is unknown or pdu.Params is nil.
func (pdu Pdu) countPayloadSize() (size int) {
	if pdu.Flags&FlagNonDefaultContext != 0 {
		size = octetStringEncodingSize(pdu.Context)
	}
	size += pdu.Params.encodingSize()
	if pduTable[pdu.Tag].includesRanges {
		for _, r := range pdu.Ranges {
			size += searchRangeEncodingSize(r)
		}
	}
	if pduTable[pdu.Tag].includesVarbinds {
		for _, vb := range pdu.Varbinds {
			size += varbindEncodingSize(vb)
		}
	}
	return
}

// EncodePdu encodes pdu. The returned data is ready for
// sending across the wire.
// If ok = false, call EncodePduDbg to get a detailed log
// of encoding process.
func EncodePdu(pdu Pdu) (data []byte, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	// This panics if pdu.Tag is unknown or pdu.Params = nil
	if !pduTable[pdu.Tag].isApplicableParams(pdu.Params) {
		ok = false
		return
	}

	// This panics if pdu.Tag is unknown or pdu.Params = nil
	payloadSize := pdu.countPayloadSize()
	pdu.PayloadSize = int32(payloadSize)

	data = make([]byte, 0, PduHeaderSize+payloadSize)
	e := encoder{byteOrder: pdu.Flags.byteOrder()}

	data = append(data, 1, byte(pdu.Tag), byte(pdu.Flags), 0)
	data = e.appendUint32(data, pdu.SessionId)
	data = e.appendUint32(data, pdu.TransactionId)
	data = e.appendUint32(data, pdu.PacketId)
	data = e.appendInt32(data, pdu.PayloadSize)

	if pdu.Flags&FlagNonDefaultContext != 0 {
		data = e.appendOctetString(data, pdu.Context)
	}

	// This panics if pdu.Params = nil
	data = pdu.Params.append(e, data)

	// This panics if pdu.Tag is unknown
	if pduTable[pdu.Tag].includesRanges {
		for _, r := range pdu.Ranges {
			data = e.appendSearchRange(data, r)
		}
	}

	// This panics if pdu.Tag is unknown
	if pduTable[pdu.Tag].includesVarbinds {
		for _, vb := range pdu.Varbinds {
			data = e.appendVarbind(data, vb)
		}
	}

	ok = true
	return
}

// DecodePduHeader returns pdu with header fields set to
// the result of decoding data.
// If ok = false, call DecodePduHeaderDbg to get a detailed log
// of decoding process.
func DecodePduHeader(data []byte) (pdu Pdu, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	pdu.Tag, pdu.Flags = PduTag(data[1]), Flags(data[2])
	if ok = data[0] == 1 && pdu.Tag.IsKnown() && data[3] == 0; !ok {
		return
	}
	if ok = pdu.Tag.NotAllowedFlags(pdu.Flags) == 0; !ok {
		return
	}

	d := decoder{byteOrder: pdu.Flags.byteOrder()}
	pdu.SessionId, data = d.parseUint32(data[4:])
	pdu.TransactionId, data = d.parseUint32(data)
	pdu.PacketId, data = d.parseUint32(data)
	pdu.PayloadSize, data = d.parseInt32(data)

	ok = pdu.PayloadSize >= 0 && pdu.PayloadSize&0x03 == 0 && len(data) == 0
	return
}

// DecodePduPayload sets pdu payload fields to the result
// of decoding data. If ok = false, call DecodePduPayloadDbg
// to get a detailed log of decoding process.
func DecodePduPayload(pdu *Pdu, data []byte) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	d := decoder{byteOrder: pdu.Flags.byteOrder()}

	if pdu.Flags&FlagNonDefaultContext != 0 {
		pdu.Context, data = d.parseOctetString(data)
	}

	pdu.Params, data = pduTable[pdu.Tag].parseParams(d, data)

	if pduTable[pdu.Tag].includesRanges {
		for len(data) > 0 {
			var r SearchRange
			r, data = d.parseSearchRange(data)
			pdu.Ranges = append(pdu.Ranges, r)
		}
	}
	if pduTable[pdu.Tag].includesVarbinds {
		for len(data) > 0 {
			var vb asn.Varbind
			vb, data = d.parseVarbind(data)
			pdu.Varbinds = append(pdu.Varbinds, vb)
		}
	}

	ok = len(data) == 0
	return
}

// EncodePduDbg is functionally equivalent to EncodePdu
// and additionally logs the process of encoding.
func EncodePduDbg(pdu Pdu, log logger.Log) (data []byte, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
		if !ok {
			log.Write("!!! Failed to encode pdu.")
		}
	}()

	log.Writef("encoding pdu\n%s", pdu.String(1))
	if ok = pdu.Tag.IsKnown(); !ok {
		log.Write("!!! Error. Unknown pdu tag")
		return
	}
	if ok = pdu.Params != nil; !ok {
		log.Write("!!! Error. Payload params are not set")
		return
	}
	if ok = pduTable[pdu.Tag].isApplicableParams(pdu.Params); !ok {
		log.Write("!!! Error. Invalid type of payload params")
		return
	}

	payloadSize := pdu.countPayloadSize()
	pdu.PayloadSize = int32(payloadSize)
	log.Writef("counted payload size: %d", payloadSize)

	log.Writef("creating data slice with cap = %d", PduHeaderSize+payloadSize)
	data = make([]byte, 0, PduHeaderSize+payloadSize)

	e := encoderDbg{
		byteOrder: pdu.Flags.byteOrder(),
		log:       log,
	}
	log.Writef("using %s", e.byteOrder.String())

	log.Writef("appending Version: 1, Tag: %s, Flags: %s, reserved byte: 0",
		pdu.Tag.String(), pdu.Flags.String())
	data = append(data, 1, byte(pdu.Tag), byte(pdu.Flags), 0)
	log.Write(hex.Dump("    ", data))

	log.Write("appending SessionId")
	data = e.appendUint32(data, pdu.SessionId)
	log.Write("appending TransactionId")
	data = e.appendUint32(data, pdu.TransactionId)
	log.Write("appending PacketId")
	data = e.appendUint32(data, pdu.PacketId)
	log.Write("appending PayloadSize")
	data = e.appendInt32(data, pdu.PayloadSize)

	if pdu.Flags&FlagNonDefaultContext != 0 {
		log.Write("appending Context")
		data = e.appendOctetString(data, pdu.Context)
	}

	data = pdu.Params.appendDbg(e, data)

	if pduTable[pdu.Tag].includesRanges {
		for _, r := range pdu.Ranges {
			data = e.appendSearchRange(data, r)
		}
	}
	if pduTable[pdu.Tag].includesVarbinds {
		for _, vb := range pdu.Varbinds {
			data = e.appendVarbind(data, vb)
		}
	}

	log.Writef("encoded pdu data:\n%s",
		hex.Dump("    ", data))

	ok = true
	return
}

// DecodePduHeaderDbg is functionally equivalent to
// DecodePduHeader and additionally logs the process
// of decoding.
func DecodePduHeaderDbg(data []byte, log logger.Log) (pdu Pdu, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
		if !ok {
			log.Write("!!! Failed to decode pdu header.")
		}
	}()

	log.Writef("decoding pdu header data:\n%s",
		hex.Dump("    ", data))
	if ok = len(data) == PduHeaderSize; !ok {
		log.Writef("!!! Error. Data size %d != %d",
			len(data), PduHeaderSize)
		return
	}

	d := decoderDbg{data: data, log: log}

	log.Write("parsing first 4 bytes")
	data, nextpos := d.nextChunk(0, 4)
	pdu.Tag, pdu.Flags = PduTag(data[1]), Flags(data[2])
	log.Writef("parsed agentx version: %d, pdu.Tag: %s, pdu.Flags: %s",
		data[0], pdu.Tag.String(), pdu.Flags.String())
	if ok = data[0] == 1; !ok {
		log.Writef("!!! Error. Invalid agentx version", data[0])
		return
	}
	if ok = pdu.Tag.IsKnown(); !ok {
		log.Write("!!! Error. Unknown pdu tag")
		return
	}
	fooFlags := pdu.Tag.NotAllowedFlags(pdu.Flags)
	if ok = fooFlags == 0; !ok {
		log.Writef("!!! Error. Flags %s are not allowed for this pdu", fooFlags.String())
		return
	}
	if ok = data[3] == 0; !ok {
		log.Write("!!! Error. Non-zero reserved byte")
		return
	}

	d.byteOrder = pdu.Flags.byteOrder()
	log.Writef("using %s", d.byteOrder.String())

	log.Write("parsing pdu.SessionId")
	pdu.SessionId, nextpos = d.parseUint32(nextpos)
	log.Write("parsing pdu.TransactionId")
	pdu.TransactionId, nextpos = d.parseUint32(nextpos)
	log.Write("parsing pdu.PacketId")
	pdu.PacketId, nextpos = d.parseUint32(nextpos)

	log.Write("parsing pdu.PayloadSize")
	pdu.PayloadSize, nextpos = d.parseInt32(nextpos)
	if ok = pdu.PayloadSize >= 0; !ok {
		log.Write("!!! Error. Negative payload size")
		return
	}
	if ok = pdu.PayloadSize&0x03 == 0; !ok {
		log.Write("!!! Error. Payload size is not multiple of 4")
		return
	}

	if ok = nextpos == len(d.data); !ok {
		log.Writef("!!! Error. Not decoded extra bytes:\n%s",
			hex.DumpSub("    ", d.data, nextpos, len(data)))
	}

	return
}

// DecodePduPayloadDbg functionally is equivalent to
// DecodePduPayload and additionally logs the process
// of decoding.
func DecodePduPayloadDbg(pdu *Pdu, data []byte, log logger.Log) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
		if !ok {
			log.Write("!!! Failed to decode pdu payload")
		}
	}()

	log.Writef("decoding pdu payload data:\n%s",
		hex.Dump("    ", data))

	if ok = len(data) == int(pdu.PayloadSize); !ok {
		log.Writef("!!! Error. Data size %d != pdu.PayloadSize %d",
			len(data), pdu.PayloadSize)
		return
	}

	d := decoderDbg{
		byteOrder: pdu.Flags.byteOrder(),
		data:      data,
		log:       log,
	}
	log.Writef("using %s", d.byteOrder.String())

	nextpos := 0
	if pdu.Flags&FlagNonDefaultContext != 0 {
		log.Write("parsing Context")
		pdu.Context, nextpos = d.parseOctetString(nextpos)
	}

	pdu.Params, nextpos = pduTable[pdu.Tag].parseParamsDbg(d, nextpos)

	if pduTable[pdu.Tag].includesRanges {
		for nextpos < len(d.data) {
			var r SearchRange
			r, nextpos = d.parseSearchRange(nextpos)
			pdu.Ranges = append(pdu.Ranges, r)
		}
	}
	if pduTable[pdu.Tag].includesVarbinds {
		for nextpos < len(d.data) {
			var vb asn.Varbind
			vb, nextpos = d.parseVarbind(nextpos)
			pdu.Varbinds = append(pdu.Varbinds, vb)
		}
	}

	if ok = nextpos == len(data); !ok {
		log.Writef("!!! Error. Not decoded extra bytes:\n%s",
			hex.DumpSub("    ", data, nextpos, len(data)))
	}

	return
}
