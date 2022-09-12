package agentx

import (
	"github.com/alexispb/mygosnmp/asn"
	"github.com/alexispb/mygosnmp/hex"
)

func varbindEncodingSize(vb asn.Varbind) int {
	return 4 + objectIdEncodingSize(vb.Oid) +
		vbtable[vb.Tag].encodingSize(vb.Value)
}

func (e encoder) appendVarbind(data []byte, vb asn.Varbind) []byte {
	data = e.appendInt16(data, int16(vb.Tag))
	data = append(data, 0, 0)
	data = e.appendObjectId(data, vb.Oid, 0)
	return vbtable[vb.Tag].appendValue(e, data, vb.Value)
}

func (d decoder) parseVarbind(data []byte) (vb asn.Varbind, next []byte) {
	var tag int16
	tag, next = d.parseInt16(data)
	if vb.Tag = asn.Tag(tag); !vb.Tag.IsValueTag() {
		panic(ErrEncoding)
	}
	if next[0] != 0 || next[1] != 0 {
		panic(ErrEncoding)
	}
	vb.Oid, _, next = d.parseObjectId(next[2:])
	vb.Value, next = vbtable[vb.Tag].parseValue(d, next)
	return
}

func (e encoderDbg) appendVarbind(data []byte, vb asn.Varbind) []byte {
	e.log.Writef("appending Varbind: %s", vb.String())
	startindex := len(data)

	e.log.Write("appending Varbind Tag")
	data = e.appendInt16(data, int16(vb.Tag))

	e.log.Write("appending 2 reserved bytes")
	data = append(data, 0, 0)
	e.log.Write(hex.DumpTail("    ", data, 2))

	e.log.Write("appending Varbind Oid")
	data = e.appendObjectId(data, vb.Oid, 0)

	e.log.Write("appending Varbind Value")
	data = vbtable[vb.Tag].appendValueDbg(e, data, vb.Value)

	e.log.Writef("appended Varbind data:\n%s")
	hex.DumpSub("    ", data, startindex, len(data))
	return data
}

func (d decoderDbg) parseVarbind(startpos int) (vb asn.Varbind, nextpos int) {
	d.log.Write("parsing Varbind")
	d.log.Write("parsing Varbind.Tag")
	var tag int16
	tag, nextpos = d.parseInt16(startpos)
	vb.Tag = asn.Tag(tag)
	d.log.Writef("parsed Tag: %s", vb.Tag.String())
	if vb.Tag = asn.Tag(tag); !vb.Tag.IsValueTag() {
		panic("Unknown value tag")
	}
	d.log.Write("parsing 2 reserved bytes")
	data, nextpos := d.nextChunk(nextpos, 2)
	if data[0] != 0 || data[1] != 0 {
		panic("non-zero reserved bytes")
	}
	d.log.Write("parsing Varbind.Oid")
	vb.Oid, _, nextpos = d.parseObjectId(nextpos)
	d.log.Write("parsing Varbind.Value")
	vb.Value, nextpos = vbtable[vb.Tag].parseValueDbg(d, nextpos)
	d.log.Writef("parsed Varbind: %s", vb.String())
	return
}

type vbentry struct {
	encodingSize   func(v interface{}) int
	appendValue    func(e encoder, data []byte, v interface{}) []byte
	parseValue     func(d decoder, data []byte) (interface{}, []byte)
	appendValueDbg func(e encoderDbg, data []byte, v interface{}) []byte
	parseValueDbg  func(d decoderDbg, startpos int) (v interface{}, nextpos int)
}

var vbtable = map[asn.Tag]vbentry{
	asn.TagInteger32: {
		encodingSize: func(v interface{}) int {
			return 4
		},
		appendValue: func(e encoder, data []byte, v interface{}) []byte {
			return e.appendInt32(data, v.(int32))
		},
		parseValue: func(d decoder, data []byte) (interface{}, []byte) {
			return d.parseInt32(data)
		},
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte {
			return e.appendInt32(data, v.(int32))
		},
		parseValueDbg: func(d decoderDbg, startpos int) (v interface{}, nextpos int) {
			return d.parseInt32(startpos)
		},
	},
	asn.TagOctetString: {
		encodingSize: func(v interface{}) int {
			return octetStringEncodingSize(v.(string))
		},
		appendValue: func(e encoder, data []byte, v interface{}) []byte {
			return e.appendOctetString(data, v.(string))
		},
		parseValue: func(d decoder, data []byte) (interface{}, []byte) {
			return d.parseOctetString(data)
		},
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte {
			return e.appendOctetString(data, v.(string))
		},
		parseValueDbg: func(d decoderDbg, startpos int) (v interface{}, nextpos int) {
			return d.parseOctetString(startpos)
		},
	},
	asn.TagObjectId: {
		encodingSize: func(v interface{}) int {
			return objectIdEncodingSize(v.([]uint32))
		},
		appendValue: func(e encoder, data []byte, v interface{}) []byte {
			return e.appendObjectId(data, v.([]uint32), 0)
		},
		parseValue: func(d decoder, data []byte) (v interface{}, next []byte) {
			v, _, next = d.parseObjectId(data)
			return
		},
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte {
			return e.appendObjectId(data, v.([]uint32), 0)
		},
		parseValueDbg: func(d decoderDbg, startpos int) (v interface{}, nextpos int) {
			v, _, nextpos = d.parseObjectId(startpos)
			return
		},
	},
	asn.TagIpAddress: {
		encodingSize: func(v interface{}) int {
			return ipAddressEncodingSize(v.([4]byte))
		},
		appendValue: func(e encoder, data []byte, v interface{}) []byte {
			return e.appendIpAddress(data, v.([4]byte))
		},
		parseValue: func(d decoder, data []byte) (interface{}, []byte) {
			return d.parseIpAddress(data)
		},
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte {
			return e.appendIpAddress(data, v.([4]byte))
		},
		parseValueDbg: func(d decoderDbg, startpos int) (v interface{}, nextpos int) {
			return d.parseIpAddress(startpos)
		},
	},
	asn.TagCounter32: {
		encodingSize: func(v interface{}) int {
			return 4
		},
		appendValue: func(e encoder, data []byte, v interface{}) []byte {
			return e.appendUint32(data, v.(uint32))
		},
		parseValue: func(d decoder, data []byte) (interface{}, []byte) {
			return d.parseUint32(data)
		},
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte {
			return e.appendUint32(data, v.(uint32))
		},
		parseValueDbg: func(d decoderDbg, startpos int) (v interface{}, nextpos int) {
			return d.parseUint32(startpos)
		},
	},
	asn.TagGauge32: {
		encodingSize: func(v interface{}) int {
			return 4
		},
		appendValue: func(e encoder, data []byte, v interface{}) []byte {
			return e.appendUint32(data, v.(uint32))
		},
		parseValue: func(d decoder, data []byte) (interface{}, []byte) {
			return d.parseUint32(data)
		},
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte {
			return e.appendUint32(data, v.(uint32))
		},
		parseValueDbg: func(d decoderDbg, startpos int) (v interface{}, nextpos int) {
			return d.parseUint32(startpos)
		},
	},
	asn.TagTimeTicks: {
		encodingSize: func(v interface{}) int {
			return 4
		},
		appendValue: func(e encoder, data []byte, v interface{}) []byte {
			return e.appendUint32(data, v.(uint32))
		},
		parseValue: func(d decoder, data []byte) (interface{}, []byte) {
			return d.parseUint32(data)
		},
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte {
			return e.appendUint32(data, v.(uint32))
		},
		parseValueDbg: func(d decoderDbg, startpos int) (v interface{}, nextpos int) {
			return d.parseUint32(startpos)
		},
	},
	asn.TagOpaque: {
		encodingSize: func(v interface{}) int {
			return opaqueEncodingSize(v.([]byte))
		},
		appendValue: func(e encoder, data []byte, v interface{}) []byte {
			return e.appendOpaque(data, v.([]byte))
		},
		parseValue: func(d decoder, data []byte) (interface{}, []byte) {
			return d.parseOpaque(data)
		},
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte {
			return e.appendOpaque(data, v.([]byte))
		},
		parseValueDbg: func(d decoderDbg, startpos int) (v interface{}, nextpos int) {
			return d.parseOpaque(startpos)
		},
	},
	asn.TagCounter64: {
		encodingSize: func(v interface{}) int {
			return 8
		},
		appendValue: func(e encoder, data []byte, v interface{}) []byte {
			return e.appendUint64(data, v.(uint64))
		},
		parseValue: func(d decoder, data []byte) (interface{}, []byte) {
			return d.parseUint64(data)
		},
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte {
			return e.appendUint64(data, v.(uint64))
		},
		parseValueDbg: func(d decoderDbg, startpos int) (v interface{}, nextpos int) {
			return d.parseUint64(startpos)
		},
	},
	asn.TagNull: {
		encodingSize:   func(v interface{}) int { return 0 },
		appendValue:    func(e encoder, data []byte, v interface{}) []byte { return data },
		parseValue:     func(d decoder, data []byte) (interface{}, []byte) { return nil, data },
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte { return data },
		parseValueDbg:  func(d decoderDbg, startpos int) (v interface{}, nextpos int) { return nil, startpos },
	},
	asn.TagNoSuchObject: {
		encodingSize:   func(v interface{}) int { return 0 },
		appendValue:    func(e encoder, data []byte, v interface{}) []byte { return data },
		parseValue:     func(d decoder, data []byte) (interface{}, []byte) { return nil, data },
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte { return data },
		parseValueDbg:  func(d decoderDbg, startpos int) (v interface{}, nextpos int) { return nil, startpos },
	},
	asn.TagNoSuchInstance: {
		encodingSize:   func(v interface{}) int { return 0 },
		appendValue:    func(e encoder, data []byte, v interface{}) []byte { return data },
		parseValue:     func(d decoder, data []byte) (interface{}, []byte) { return nil, data },
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte { return data },
		parseValueDbg:  func(d decoderDbg, startpos int) (v interface{}, nextpos int) { return nil, startpos },
	},
	asn.TagEndOfMibView: {
		encodingSize:   func(v interface{}) int { return 0 },
		appendValue:    func(e encoder, data []byte, v interface{}) []byte { return data },
		parseValue:     func(d decoder, data []byte) (interface{}, []byte) { return nil, data },
		appendValueDbg: func(e encoderDbg, data []byte, v interface{}) []byte { return data },
		parseValueDbg:  func(d decoderDbg, startpos int) (v interface{}, nextpos int) { return nil, startpos },
	},
}
