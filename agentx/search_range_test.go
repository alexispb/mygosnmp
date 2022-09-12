package agentx

import (
	"fmt"
	"testing"

	"github.com/alexispb/mygosnmp/hex"
)

var searchRangeTestData = struct {
	order [2]byteOrder
	test  []struct {
		value SearchRange
		data  [2][]byte
	}
}{
	order: [2]byteOrder{
		FlagNetworkByteOrder.byteOrder(),
		FlagsNone.byteOrder(),
	},
	test: []struct {
		value SearchRange
		data  [2][]byte
	}{
		{
			value: SearchRange{},
			data: [2][]byte{
				{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
				{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: SearchRange{
				StartOid:      []uint32{1, 3, 6, 1, 4},
				EndOid:        []uint32{},
				StartIncluded: 1},
			data: [2][]byte{
				{
					0x00, 0x04, 0x01, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
				{
					0x00, 0x04, 0x01, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
			},
		},
		{
			value: SearchRange{
				StartOid:      []uint32{1, 3, 6, 1, 4, 999, 1},
				EndOid:        []uint32{1, 3, 6, 1, 4, 999, 5},
				StartIncluded: 1},
			data: [2][]byte{
				{
					0x02, 0x04, 0x01, 0x00,
					0x00, 0x00, 0x03, 0xE7,
					0x00, 0x00, 0x00, 0x01,
					0x02, 0x04, 0x00, 0x00,
					0x00, 0x00, 0x03, 0xE7,
					0x00, 0x00, 0x00, 0x05,
				},
				{
					0x02, 0x04, 0x01, 0x00,
					0xE7, 0x03, 0x00, 0x00,
					0x01, 0x00, 0x00, 0x00,
					0x02, 0x04, 0x00, 0x00,
					0xE7, 0x03, 0x00, 0x00,
					0x05, 0x00, 0x00, 0x00,
				},
			},
		},
	},
}

func TestSearchRangeEncoding(t *testing.T) {
	for i, order := range searchRangeTestData.order {
		for _, test := range searchRangeTestData.test {
			testid := fmt.Sprintf("%s %s\n", test.value.String(), order.String())

			size := searchRangeEncodingSize(test.value)
			if size != len(test.data[i]) {
				t.Errorf("%sinvalid evaluated encoding size: %d != %d",
					testid, size, len(test.data[i]))
			}

			data := encoder{byteOrder: order}.
				appendSearchRange(nil, test.value)

			if diff := hex.DumpDiff(test.data[i], data); len(diff) != 0 {
				t.Errorf("%sinvalid encoded data:\n%s",
					testid, diff)
			}

			value, tail := decoder{byteOrder: order}.
				parseSearchRange(test.data[i])

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
				appendSearchRange(nil, test.value)

			if diff := hex.DumpDiff(test.data[i], data); len(diff) != 0 {
				t.Errorf("%sinvalid debug-encoded data:\n%s",
					testid, diff)
			}

			value, nextpos := decoderDbg{byteOrder: order, data: test.data[i], log: lognone}.
				parseSearchRange(0)

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
