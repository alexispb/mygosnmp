package agentx

import (
	"strconv"
	"strings"

	"github.com/alexispb/mygosnmp/hex"
	"github.com/alexispb/mygosnmp/oid"
)

type SearchRange struct {
	StartOid      []uint32
	EndOid        []uint32
	StartIncluded byte
}

func searchRangeEncodingSize(r SearchRange) int {
	return objectIdEncodingSize(r.StartOid) + objectIdEncodingSize(r.EndOid)
}

func (e encoder) appendSearchRange(data []byte, r SearchRange) []byte {
	data = e.appendObjectId(data, r.StartOid, r.StartIncluded)
	return e.appendObjectId(data, r.EndOid, 0)
}

func (e encoderDbg) appendSearchRange(data []byte, r SearchRange) []byte {
	e.log.Writef("appending SearchRange: %s", r.String())
	startindex := len(data)

	data = e.appendObjectId(data, r.StartOid, r.StartIncluded)
	data = e.appendObjectId(data, r.EndOid, 0)

	e.log.Writef("appended SearchRange data\n%s",
		hex.DumpSub("    ", data, startindex, len(data)))
	return data
}

func (d decoder) parseSearchRange(data []byte) (r SearchRange, next []byte) {
	r.StartOid, r.StartIncluded, next = d.parseObjectId(data)
	r.EndOid, _, next = d.parseObjectId(next)
	return
}

func (d decoderDbg) parseSearchRange(startpos int) (r SearchRange, nextpos int) {
	d.log.Write("parsing SearchRange")
	d.log.Write("parsing StartOid with include")
	r.StartOid, r.StartIncluded, nextpos = d.parseObjectId(startpos)
	d.log.Write("parsing EndOid")
	r.EndOid, _, nextpos = d.parseObjectId(nextpos)
	d.log.Writef("parsed SearchRange: %s", r.String())
	return
}

// Fprint appends search range string representation to the builder.
func (r SearchRange) Fprint(sb *strings.Builder) {
	sb.WriteString("{StartOid: ")
	oid.Fprint(sb, r.StartOid)
	sb.WriteString(" EndOid: ")
	oid.Fprint(sb, r.EndOid)
	sb.WriteString(" StartIncluded: ")
	sb.WriteString(strconv.FormatUint(uint64(r.StartIncluded), 10))
	sb.WriteByte('}')
}

// String return search range string representation.
func (r SearchRange) String() string {
	var sb strings.Builder
	r.Fprint(&sb)
	return sb.String()
}
