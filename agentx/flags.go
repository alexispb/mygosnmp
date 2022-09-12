package agentx

import (
	"strings"
)

type Flags byte

const (
	FlagsNone                Flags = 0
	FlagInstanceRegistration Flags = 1 << 0
	FlagNewIndex             Flags = 1 << 1
	FlagAnyIndex             Flags = 1 << 2
	FlagNonDefaultContext    Flags = 1 << 3
	FlagNetworkByteOrder     Flags = 1 << 4
)

var flagStrings = [...]string{
	"InstanceRegistration",
	"NewIndex",
	"AnyIndex",
	"NonDefaultContext",
	"NetworkByteOrder",
}

const hextable = "0123456789ABCDEF"
const flagsStringMaxLen = 96

// Fprintf appends Flags string representation to sb.
func (f Flags) Fprint(sb *strings.Builder) {
	if f == FlagsNone {
		return
	}

	first := true
	for i := 0; i < 8; i++ {
		flag := Flags(1 << i)
		if f&flag == 0 {
			continue
		}
		if first {
			first = false
		} else {
			sb.WriteByte('|')
		}
		if i < len(flagStrings) {
			sb.WriteString(flagStrings[i])
		} else {
			sb.WriteString("?0x")
			foo := f & flag
			sb.WriteByte(hextable[foo>>4])
			sb.WriteByte(hextable[foo&0x0F])
		}
	}
}

// String returns Flags string representation.
func (f Flags) String() string {
	var sb strings.Builder
	sb.Grow(flagsStringMaxLen)
	f.Fprint(&sb)
	return sb.String()
}

// byteOrder returns byteOrder interface implementation
// which is required by the FlagNetworkByteOrder.
func (f Flags) byteOrder() byteOrder {
	if f&FlagNetworkByteOrder != 0 {
		return networkByteOrder{}
	}
	return inverseByteOrder{}
}
