package asn

import (
	"strconv"
	"strings"

	"github.com/alexispb/mygosnmp/ipa"
	"github.com/alexispb/mygosnmp/oid"
)

func (tag Tag) IsKnown() (ok bool) {
	_, ok = table[tag]
	return
}

func (tag Tag) IsValueTag() (ok bool) {
	_, ok = table[tag]
	return
}

func (tag Tag) String() string {
	if !tag.IsKnown() {
		return "?" + strconv.FormatInt(int64(tag), 10)
	}
	return table[tag].tagstr
}

func (tag Tag) IsValidValue(v interface{}) bool {
	return tag.IsKnown() && table[tag].isvalid(v)
}

func (tag Tag) Fprint(sb *strings.Builder, v interface{}) {
	if tag.IsKnown() && table[tag].isvalid(v) {
		table[tag].fprint(sb, v)
	} else {
		sb.WriteByte('?')
	}
}

type entry struct {
	// tagstr is string representation.
	tagstr string
	// isvalid tests whether value has valid type.
	isvalid func(v interface{}) bool
	// fprint appends value string representation to the builder.
	fprint func(sb *strings.Builder, v interface{})
}

var table = map[Tag]entry{
	TagInteger32: {
		tagstr: "Integer32",
		isvalid: func(v interface{}) bool {
			_, ok := v.(int32)
			return ok
		},
		fprint: func(sb *strings.Builder, v interface{}) {
			sb.WriteString(strconv.FormatInt(int64(v.(int32)), 10))
		},
	},
	TagOctetString: {
		tagstr: "OctetString",
		isvalid: func(v interface{}) bool {
			_, ok := v.(string)
			return ok
		},
		fprint: func(sb *strings.Builder, v interface{}) {
			sb.WriteString(v.(string))
		},
	},
	TagNull: {
		tagstr:  "Null",
		isvalid: func(v interface{}) bool { return true },
		fprint:  func(sb *strings.Builder, v interface{}) {},
	},
	TagObjectId: {
		tagstr: "ObjectId",
		isvalid: func(v interface{}) bool {
			_, ok := v.([]uint32)
			return ok
		},
		fprint: func(sb *strings.Builder, v interface{}) {
			oid.Fprint(sb, v.([]uint32))
		},
	},
	TagSequence: {
		tagstr:  "Sequence",
		isvalid: func(v interface{}) bool { return false },
		fprint:  func(sb *strings.Builder, v interface{}) {},
	},
	TagIpAddress: {
		tagstr: "IpAddress",
		isvalid: func(v interface{}) bool {
			_, ok := v.([4]byte)
			return ok
		},
		fprint: func(sb *strings.Builder, v interface{}) {
			ipa.Fprint(sb, v.([4]byte))
		},
	},
	TagCounter32: {
		tagstr: "Counter32",
		isvalid: func(v interface{}) bool {
			_, ok := v.(uint32)
			return ok
		},
		fprint: func(sb *strings.Builder, v interface{}) {
			sb.WriteString(strconv.FormatUint(uint64(v.(uint32)), 10))
		},
	},
	TagGauge32: {
		tagstr: "Unsigned32",
		isvalid: func(v interface{}) bool {
			_, ok := v.(uint32)
			return ok
		},
		fprint: func(sb *strings.Builder, v interface{}) {
			sb.WriteString(strconv.FormatUint(uint64(v.(uint32)), 10))
		},
	},
	TagTimeTicks: {
		tagstr: "TimeTicks",
		isvalid: func(v interface{}) bool {
			_, ok := v.(uint32)
			return ok
		},
		fprint: func(sb *strings.Builder, v interface{}) {
			sb.WriteString(strconv.FormatUint(uint64(v.(uint32)), 10))
		},
	},
	TagOpaque: {
		tagstr: "Opaque",
		isvalid: func(v interface{}) bool {
			_, ok := v.([]byte)
			return ok
		},
		fprint: func(sb *strings.Builder, v interface{}) {
			b := v.([]byte)
			sb.WriteByte('{')
			for i, n := 0, len(b); i < n; i++ {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.FormatUint(uint64(b[i]), 10))
			}
			sb.WriteByte('}')
		},
	},
	TagCounter64: {
		tagstr: "Counter64",
		isvalid: func(v interface{}) bool {
			_, ok := v.(uint64)
			return ok
		},
		fprint: func(sb *strings.Builder, v interface{}) {
			sb.WriteString(strconv.FormatUint(v.(uint64), 10))
		},
	},
	TagNoSuchObject: {
		tagstr:  "NoSuchObject",
		isvalid: func(v interface{}) bool { return true },
		fprint:  func(sb *strings.Builder, v interface{}) {},
	},
	TagNoSuchInstance: {
		tagstr:  "NoSuchInstance",
		isvalid: func(v interface{}) bool { return true },
		fprint:  func(sb *strings.Builder, v interface{}) {},
	},
	TagEndOfMibView: {
		tagstr:  "EndOfMibView",
		isvalid: func(v interface{}) bool { return true },
		fprint:  func(sb *strings.Builder, v interface{}) {},
	},
}
