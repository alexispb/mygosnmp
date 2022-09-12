package ipa

import (
	"strconv"
	"strings"
)

// Fprint appends ip address string representation to the builder.
// E.g. string representation of [4]byte{192, 158, 1, 38} is 192.158.1.38.
func Fprint(sb *strings.Builder, a [4]byte) {
	sb.WriteString(strconv.FormatUint(uint64(a[0]), 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(uint64(a[1]), 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(uint64(a[2]), 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(uint64(a[3]), 10))
}

// String return ip address string representation.
// E.g. string representation of [4]byte{192, 158, 1, 38} is 192.158.1.38.
func String(a [4]byte) string {
	var sb strings.Builder
	sb.Grow(15) // 15 = 4*(3 digits) + 3 dots
	Fprint(&sb, a)
	return sb.String()
}
