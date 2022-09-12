package agentx

import (
	"fmt"
	"testing"

	"github.com/alexispb/mygosnmp/asn"
	"github.com/alexispb/mygosnmp/hex"
)

var varbindTestData = struct {
	order [2]byteOrder
	test  []struct {
		value asn.Varbind
		data  [2][]byte
	}
}{
	order: [2]byteOrder{
		FlagNetworkByteOrder.byteOrder(),
		FlagsNone.byteOrder(),
	},
	test: []struct {
		value asn.Varbind
		data  [2][]byte
	}{
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagInteger32, Value: int32(123)},
			data: [2][]byte{
				{
					0x00, 0x02, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x7B,
				},
				{
					0x02, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x7B, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagOctetString, Value: "abc"},
			data: [2][]byte{
				{
					0x00, 0x04, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x03,
					0x61, 0x62, 0x63, 0x00,
				},
				{
					0x04, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x00,
					0x61, 0x62, 0x63, 0x00,
				},
			},
		},

		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagNull},
			data: [2][]byte{
				{
					0x00, 0x05, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
				{
					0x05, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagObjectId, Value: []uint32{}},
			data: [2][]byte{
				{
					0x00, 0x06, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
				{
					0x06, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagIpAddress, Value: [4]byte{192, 158, 1, 38}},
			data: [2][]byte{
				{
					0x00, 0x40, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x04,
					0xC0, 0x9E, 0x01, 0x26,
				},
				{
					0x40, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x04, 0x00, 0x00, 0x00,
					0xC0, 0x9E, 0x01, 0x26,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagCounter32, Value: uint32(123)},
			data: [2][]byte{
				{
					0x00, 0x41, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x7B,
				},
				{
					0x41, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x7B, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagGauge32, Value: uint32(123)},
			data: [2][]byte{
				{
					0x00, 0x42, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x7B,
				},
				{
					0x42, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x7B, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagTimeTicks, Value: uint32(123)},
			data: [2][]byte{
				{
					0x00, 0x43, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x7B,
				},
				{
					0x43, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x7B, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagOpaque, Value: []byte{}},
			data: [2][]byte{
				{
					0x00, 0x44, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
				{
					0x44, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagCounter64, Value: uint64(123)},
			data: [2][]byte{
				{
					0x00, 0x46, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7B,
				},
				{
					0x46, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x7B, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagNoSuchObject},
			data: [2][]byte{
				{
					0x00, 0x80, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
				{
					0x80, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagNoSuchInstance},
			data: [2][]byte{
				{
					0x00, 0x81, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
				{
					0x81, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: asn.Varbind{Oid: nil, Tag: asn.TagEndOfMibView},
			data: [2][]byte{
				{
					0x00, 0x82, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
				{
					0x82, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
	},
}

func TestVarbindEncoding(t *testing.T) {
	for i, order := range encodingTestData.order {
		for _, test := range varbindTestData.test {
			testid := fmt.Sprintf("%s %s\n", test.value.String(), order.String())

			size := varbindEncodingSize(test.value)
			if size != len(test.data[i]) {
				t.Errorf("%sinvalid evaluated encoding size %d != %d",
					testid, size, len(test.data[i]))
			}

			data := encoder{byteOrder: order}.
				appendVarbind(nil, test.value)

			if diff := hex.DumpDiff(test.data[i], data); len(diff) != 0 {
				t.Errorf("%sinvalid encoded data:\n%s",
					testid, diff)
			}

			value, tail := decoder{byteOrder: order}.
				parseVarbind(test.data[i])

			str1 := test.value.String()
			str2 := value.String()
			if str2 != str1 {
				t.Errorf("%sinvalid decoded value:\n%s\n%s",
					testid, str2, str1)
			}

			if len(tail) != 0 {
				t.Errorf("%sinvalid tail length: %d != 0",
					testid, len(tail))
			}

			data = encoderDbg{byteOrder: order, log: lognone}.
				appendVarbind(nil, test.value)

			if diff := hex.DumpDiff(test.data[i], data); len(diff) != 0 {
				t.Errorf("%sinvalid debug-encoded data:\n%s",
					testid, diff)
			}

			value, nextpos := decoderDbg{byteOrder: order, data: test.data[i], log: lognone}.
				parseVarbind(0)

			str1 = fmt.Sprintf("%v", test.value)
			str2 = fmt.Sprintf("%v", value)
			if str2 != str1 {
				t.Errorf("%sinvalid debug-decoded value %s != %s",
					testid, str2, str1)
			}

			if nextpos != len(test.data[i]) {
				t.Errorf("%sinvalid next position: %d != %d",
					testid, nextpos, len(test.data[i]))
			}
		}
	}
}
