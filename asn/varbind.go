package asn

import (
	"strings"

	"github.com/alexispb/mygosnmp/oid"
)

type Varbind struct {
	Oid   []uint32
	Tag   Tag
	Value interface{}
}

// Fprint appends varbind string representation to the builder.
func (vb Varbind) Fprint(sb *strings.Builder) {
	sb.WriteByte('{')
	oid.Fprint(sb, vb.Oid)
	sb.WriteByte(' ')
	sb.WriteString(vb.Tag.String())
	sb.WriteByte(':')
	vb.Tag.Fprint(sb, vb.Value)
	sb.WriteByte('}')
}

// String return varbind string representation.
func (vb Varbind) String() string {
	var sb strings.Builder
	vb.Fprint(&sb)
	return sb.String()
}
