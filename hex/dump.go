package hex

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alexispb/mygosnmp/generics"
)

const (
	diffmark1 = '\u250C'
	diffmark2 = '\u2514'
	table     = "0123456789ABCDEF"
)

// must be: len(dst) >= 10
func encodeOffset(dst []byte, offset int) {
	b := byte(offset >> 24)
	dst[0] = table[b>>4]
	dst[1] = table[b&0x0F]
	b = byte(offset >> 16)
	dst[2] = table[b>>4]
	dst[3] = table[b&0x0F]
	b = byte(offset >> 8)
	dst[4] = table[b>>4]
	dst[5] = table[b&0x0F]
	b = byte(offset)
	dst[6] = table[b>>4]
	dst[7] = table[b&0x0F]
	dst[8] = ' '
	dst[9] = ' '
}

// DumpSub returns a hex dump of data[i1:i2].
// Example:
//     data := []byte{0, 1, 3, 4, ..., 15}
//     fmt.Println(DumpSub(data, 4, len(data)))
// results with
//     00000004  04 05 06 07
//     00000008  08 09 0A 0B
//     0000000C  0C 0D 0E 0F
// If i1 = i2 HexDump returns only offset.
// Example:
//     data := []byte{0, 1, 3, 4, ..., 15}
//     fmt.Println(DumpSub(data, 4, 4))
// results with
//     00000004
// If either i1 or i2 is out of ranges DumpSub returns
// an error message.
// Example:
//     data := []byte{0, 1, 3, 4, ..., 15}
//     fmt.Println(DumpSub(data, 4, 3))
// results with
//     out of range: i1 = 4, i2 = 3, len(data) = 16
func DumpSub(prefix string, data []byte, i1, i2 int) string {
	if i1 < 0 || i1 > i2 || i2 > len(data) {
		return fmt.Sprintf(
			"%sout of range: i1 = %d, i2 = %d, len(data) = %d",
			prefix, i1, i2, len(data))
	}

	var buf [13]byte

	if i1 == i2 {
		var sb strings.Builder
		sb.Grow(10)
		encodeOffset(buf[:10], i1)
		sb.WriteString(prefix)
		sb.Write(buf[:10])
		return sb.String()
	}

	var sb strings.Builder
	// 23 = len("00000000  00 00 00 00 \n")
	sb.Grow((1 + ((i2 - i1 - 1) / 4)) * (len(prefix) + 23))

	used := 0
	for i := i1; i < i2; i++ {
		if used == 0 {
			encodeOffset(buf[:], i)
			sb.WriteString(prefix)
			sb.Write(buf[:10])
		}
		buf[10] = table[data[i]>>4]
		buf[11] = table[data[i]&0x0F]
		buf[12] = ' '
		sb.Write(buf[10:13])
		used++
		if used == 4 && i != i2-1 {
			sb.WriteByte('\n')
			used = 0
		}
	}

	return sb.String()
}

func DumpTail(prefix string, data []byte, size int) string {
	return DumpSub(prefix, data, len(data)-size, len(data))
}

func Dump(prefix string, data []byte) string {
	return DumpSub(prefix, data, 0, len(data))
}

// must be: len(dst) >= 3*len(data)
func encodeData(dst, data []byte) (size int) {
	for i := range data {
		dst[size] = table[data[i]>>4]
		size++
		dst[size] = table[data[i]&0x0F]
		size++
		dst[size] = ' '
		size++
	}
	return
}

// DumpDiff compares two slices. If slices are equal, it
// returns an empty string. Otherwise it returns a string
// which displays the discrepancy between slices.
// Example:
//     data1 := []byte{
//			0x00, 0x01, 0x02, 0x03,
//			0x04, 0x05, 0x06, 0x07
//			0x08, 0x09, 0x0A, 0x0B
//			0x0C, 0x0D, 0x0E, 0x0F
//			0x10, 0x11}
//     data2 := []byte{
//			0x00, 0x63, 0x02, 0x03,
//			0x04, 0x05, 0x06, 0x07
//			0x08, 0x09, 0x0A, 0x0B
//			0x0C, 0x0D, 0x0E, 0x0F}
//     fmt.Println(DumpDiff(data1, data2))
// results with
//     ┌00000000  00 01 02 03
//     └00000000  00 63 02 03
//      00000004  04 05 06 07
//      00000008  08 09 0A 0B
//      0000000C  0C 0D 0E 0F
//     ┌00000010  10 11
//     └00000010
func DumpDiff(data1, data2 []byte) string {
	if bytes.Equal(data1, data2) {
		return ""
	}
	var (
		sb     strings.Builder
		chunk1 []byte
		chunk2 []byte
		maxlen = generics.Max(len(data1), len(data2))
		buf    [22]byte
	)
	// 52 = 2*(diffmark:3 + len("00000000  00 00 00 00 \n)":23)
	sb.Grow((1 + ((maxlen - 1) / 4)) * 52)

	for i := 0; i < maxlen; i += 4 {
		chunk1 = data1[generics.Min(i, len(data1)):generics.Min(i+4, len(data1))]
		chunk2 = data2[generics.Min(i, len(data2)):generics.Min(i+4, len(data2))]

		encodeOffset(buf[:10], i)
		if bytes.Equal(chunk1, chunk2) {
			sb.WriteByte(' ')
			sb.Write(buf[:10+encodeData(buf[10:], chunk1)])
		} else {
			sb.WriteRune(diffmark1)
			sb.Write(buf[:10+encodeData(buf[10:], chunk1)])
			sb.WriteByte('\n')
			sb.WriteRune(diffmark2)
			sb.Write(buf[:10+encodeData(buf[10:], chunk2)])
		}
		sb.WriteByte('\n')
	}

	return strings.TrimSuffix(sb.String(), "\n")
}
