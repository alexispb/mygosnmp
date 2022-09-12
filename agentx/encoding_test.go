package agentx

import (
	"fmt"
	"testing"

	"github.com/alexispb/mygosnmp/hex"
)

var encodingTestData = struct {
	order [2]byteOrder
	test  []struct {
		value interface{}
		data  [2][]byte
	}
}{
	order: [2]byteOrder{
		FlagNetworkByteOrder.byteOrder(),
		FlagsNone.byteOrder(),
	},
	test: []struct {
		value interface{}
		data  [2][]byte
	}{
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ int values
		{
			value: int16(0x1F2F),
			data: [2][]byte{
				{0x1F, 0x2F},
				{0x2F, 0x1F},
			},
		},
		{
			value: int32(0x1F2F3F4F),
			data: [2][]byte{
				{0x1F, 0x2F, 0x3F, 0x4F},
				{0x4F, 0x3F, 0x2F, 0x1F},
			},
		},
		{
			value: uint32(0x1F2F3F4F),
			data: [2][]byte{
				{0x1F, 0x2F, 0x3F, 0x4F},
				{0x4F, 0x3F, 0x2F, 0x1F},
			},
		},
		{
			value: uint64(0x1F2F3F4F5F6F7F8F),
			data: [2][]byte{
				{0x1F, 0x2F, 0x3F, 0x4F, 0x5F, 0x6F, 0x7F, 0x8F},
				{0x8F, 0x7F, 0x6F, 0x5F, 0x4F, 0x3F, 0x2F, 0x1F},
			},
		},
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ OctetString
		{
			value: "",
			data: [2][]byte{
				{0x00, 0x00, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
		},
		{
			value: "a",
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x01,
					0x61, 0x00, 0x00, 0x00,
				},
				{
					0x01, 0x00, 0x00, 0x00,
					0x61, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: "ab",
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x02,
					0x61, 0x62, 0x00, 0x00,
				},
				{
					0x02, 0x00, 0x00, 0x00,
					0x61, 0x62, 0x00, 0x00,
				},
			},
		},
		{
			value: "abc",
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x03,
					0x61, 0x62, 0x63, 0x00,
				},
				{
					0x03, 0x00, 0x00, 0x00,
					0x61, 0x62, 0x63, 0x00,
				},
			},
		},
		{
			value: "abcd",
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x04,
					0x61, 0x62, 0x63, 0x64,
				},
				{
					0x04, 0x00, 0x00, 0x00,
					0x61, 0x62, 0x63, 0x64,
				},
			},
		},
		{
			value: "abcde",
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x05,
					0x61, 0x62, 0x63, 0x64,
					0x65, 0x00, 0x00, 0x00,
				},
				{
					0x05, 0x00, 0x00, 0x00,
					0x61, 0x62, 0x63, 0x64,
					0x65, 0x00, 0x00, 0x00,
				},
			},
		},
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ ObjectId
		{
			value: []uint32{},
			data: [2][]byte{
				{0x00, 0x00, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
		},
		{
			value: []uint32{1, 3, 6, 1},
			data: [2][]byte{
				{
					0x04, 0x00, 0x01, 0x00,
					0x00, 0x00, 0x00, 0x01,
					0x00, 0x00, 0x00, 0x03,
					0x00, 0x00, 0x00, 0x06,
					0x00, 0x00, 0x00, 0x01,
				},
				{
					0x04, 0x00, 0x01, 0x00,
					0x01, 0x00, 0x00, 0x00,
					0x03, 0x00, 0x00, 0x00,
					0x06, 0x00, 0x00, 0x00,
					0x01, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: []uint32{1, 3, 6, 1, 4},
			data: [2][]byte{
				{0x00, 0x04, 0x01, 0x00},
				{0x00, 0x04, 0x01, 0x00},
			},
		},
		{
			value: []uint32{1, 3, 6, 1, 4, 1, 9999},
			data: [2][]byte{
				{
					0x02, 0x04, 0x01, 0x00,
					0x00, 0x00, 0x00, 0x01,
					0x00, 0x00, 0x27, 0x0F,
				},
				{
					0x02, 0x04, 0x01, 0x00,
					0x01, 0x00, 0x00, 0x00,
					0x0F, 0x27, 0x00, 0x00,
				},
			},
		},
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ IpAddress
		{
			value: [4]byte{192, 158, 1, 38},
			data: [2][]byte{
				{0x00, 0x00, 0x00, 0x04, 0xC0, 0x9E, 0x01, 0x26},
				{0x04, 0x00, 0x00, 0x00, 0xC0, 0x9E, 0x01, 0x26},
			},
		},
		//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ Opaque
		{
			value: []byte{},
			data: [2][]byte{
				{0x00, 0x00, 0x00, 0x00},
				{0x00, 0x00, 0x00, 0x00},
			},
		},
		{value: []byte{0x01},
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x01,
					0x01, 0x00, 0x00, 0x00,
				},
				{
					0x01, 0x00, 0x00, 0x00,
					0x01, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: []byte{0x01, 0x02},
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x02,
					0x01, 0x02, 0x00, 0x00,
				},
				{
					0x02, 0x00, 0x00, 0x00,
					0x01, 0x02, 0x00, 0x00,
				},
			},
		},
		{
			value: []byte{0x01, 0x02, 0x03},
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x03,
					0x01, 0x02, 0x03, 0x00,
				},
				{
					0x03, 0x00, 0x00, 0x00,
					0x01, 0x02, 0x03, 0x00,
				},
			},
		},
		{
			value: []byte{0x01, 0x02, 0x03, 0x04},
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x04,
					0x01, 0x02, 0x03, 0x04,
				},
				{
					0x04, 0x00, 0x00, 0x00,
					0x01, 0x02, 0x03, 0x04,
				},
			},
		},
		{
			value: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x05,
					0x01, 0x02, 0x03, 0x04,
					0x05, 0x00, 0x00, 0x00,
				},
				{
					0x05, 0x00, 0x00, 0x00,
					0x01, 0x02, 0x03, 0x04,
					0x05, 0x00, 0x00, 0x00,
				},
			},
		},
	},
}

func TestEncoding(t *testing.T) {
	var (
		size     int         // evaluated encoded data size
		data     []byte      // encoded data
		value    interface{} // decoded value
		tail     []byte
		datadbg  []byte      // debug-encoded data
		valuedbg interface{} // debug-decoded value
		nextpos  int
	)

	for i, order := range encodingTestData.order {
		e := encoder{byteOrder: order}
		d := decoder{byteOrder: order}
		edbg := encoderDbg{byteOrder: order, log: lognone}
		ddbg := decoderDbg{byteOrder: order, log: lognone}

		for _, test := range encodingTestData.test {
			testid := fmt.Sprintf("%T %s\n", test.value, order.String())
			ddbg.data = test.data[i]

			switch testval := test.value.(type) {
			case int16:
				size = 2
				data = e.appendInt16(nil, testval)
				value, tail = d.parseInt16(test.data[i])
				datadbg = edbg.appendInt16(nil, testval)
				valuedbg, nextpos = ddbg.parseInt16(0)
			case int32:
				size = 4
				data = e.appendInt32(nil, testval)
				value, tail = d.parseInt32(test.data[i])
				datadbg = edbg.appendInt32(nil, testval)
				valuedbg, nextpos = ddbg.parseInt32(0)
			case uint32:
				size = 4
				data = e.appendUint32(nil, testval)
				value, tail = d.parseUint32(test.data[i])
				datadbg = edbg.appendUint32(nil, testval)
				valuedbg, nextpos = ddbg.parseUint32(0)
			case uint64:
				size = 8
				data = e.appendUint64(nil, testval)
				value, tail = d.parseUint64(test.data[i])
				datadbg = edbg.appendUint64(nil, testval)
				valuedbg, nextpos = ddbg.parseUint64(0)
			case string:
				size = octetStringEncodingSize(testval)
				data = e.appendOctetString(nil, testval)
				value, tail = d.parseOctetString(test.data[i])
				datadbg = edbg.appendOctetString(nil, testval)
				valuedbg, nextpos = ddbg.parseOctetString(0)
			case []uint32:
				size = objectIdEncodingSize(testval)
				data = e.appendObjectId(nil, testval, 1)
				value, _, tail = d.parseObjectId(test.data[i])
				datadbg = edbg.appendObjectId(nil, testval, 1)
				valuedbg, _, nextpos = ddbg.parseObjectId(0)
			case [4]byte:
				size = ipAddressEncodingSize(testval)
				data = e.appendIpAddress(nil, testval)
				value, tail = d.parseIpAddress(test.data[i])
				datadbg = edbg.appendIpAddress(nil, testval)
				valuedbg, nextpos = ddbg.parseIpAddress(0)
			case []byte:
				size = opaqueEncodingSize(testval)
				data = e.appendOpaque(nil, testval)
				value, tail = d.parseOpaque(test.data[i])
				datadbg = edbg.appendOpaque(nil, testval)
				valuedbg, nextpos = ddbg.parseOpaque(0)
			default:
				continue
			}

			if size != len(test.data[i]) {
				t.Errorf("%sinvalid evaluated encoding size %d != %d",
					testid, size, len(test.data[i]))
			}

			if diff := hex.DumpDiff(test.data[i], data); len(diff) != 0 {
				t.Errorf("%sinvalid encoded data\n%s",
					testid, diff)
			}

			str1 := fmt.Sprintf("%v", test.value)
			str2 := fmt.Sprintf("%v", value)

			if str2 != str1 {
				t.Errorf("%sninvalid decoded value %s != %s",
					testid, str2, str1)
			}

			if len(tail) != 0 {
				t.Errorf("%sinvalid tail length: %d != 0",
					testid, len(tail))
			}

			if diff := hex.DumpDiff(test.data[i], datadbg); len(diff) != 0 {
				t.Errorf("%sinvalid debug-encoded data\n%s",
					testid, diff)
			}

			str1 = fmt.Sprintf("%v", test.value)
			str2 = fmt.Sprintf("%v", valuedbg)

			if str2 != str1 {
				t.Errorf("%sninvalid debug-decoded value %s != %s",
					testid, str2, str1)
			}

			if nextpos != len(test.data[i]) {
				t.Errorf("%sinvalid next position: %d != %d",
					testid, nextpos, len(test.data[i]))
			}
		}
	}

}
