package oid

import "testing"

var testDataOidClone = []struct {
	id []uint32
}{
	{id: nil},
	{id: []uint32{}},
	{id: []uint32{1, 2, 3}},
}

func TestOidClone(t *testing.T) {
	for i, test := range testDataOidClone {
		if res := Clone(test.id); !Eq(res, test.id) {
			t.Errorf("TestOidClone[%d]", i)
		}
	}
}

var testDataOidCat = []struct {
	id1 []uint32
	id2 []uint32
	res []uint32
}{
	{id1: nil, id2: nil, res: []uint32{}},
	{id1: []uint32{}, id2: nil, res: []uint32{}},
	{id1: []uint32{}, id2: []uint32{}, res: []uint32{}},
	{id1: []uint32{}, id2: []uint32{1, 2, 3}, res: []uint32{1, 2, 3}},
	{id1: []uint32{1, 2, 3}, id2: []uint32{}, res: []uint32{1, 2, 3}},
	{id1: []uint32{1, 2, 3}, id2: []uint32{4, 5, 6}, res: []uint32{1, 2, 3, 4, 5, 6}},
}

func TestOidCat(t *testing.T) {
	for i, test := range testDataOidCat {
		res := Cat(test.id1, test.id2...)
		if !Eq(res, test.res) {
			t.Errorf("TestOidCat[%d]", i)
		}
	}
}

var testDataOidHasPrefix = []struct {
	id1 []uint32
	id2 []uint32
	res bool
}{
	{id1: nil, id2: nil, res: true},
	{id1: nil, id2: []uint32{}, res: true},
	{id1: []uint32{}, id2: nil, res: true},
	{id1: []uint32{}, id2: []uint32{}, res: true},
	{id1: []uint32{1}, id2: nil, res: true},
	{id1: []uint32{1}, id2: []uint32{}, res: true},
	{id1: []uint32{1, 2, 3}, id2: []uint32{1, 2}, res: true},
	{id1: []uint32{1, 2, 3}, id2: []uint32{1, 2, 3}, res: true},
	{id1: []uint32{1, 2, 3}, id2: []uint32{1, 2, 3, 4}, res: false},
	{id1: []uint32{1, 2, 3}, id2: []uint32{1, 3}, res: false},
}

func TestOidHasPrefix(t *testing.T) {
	for i, test := range testDataOidHasPrefix {
		if res := HasPrefix(test.id1, test.id2...); res != test.res {
			t.Errorf("TestOidHasPrefix[%d]", i)
		}
	}
}

var testDataOidCompare = []struct {
	id1 []uint32
	id2 []uint32
	res int
}{
	// id1 == id2
	{id1: nil, id2: nil, res: 0},
	{id1: []uint32{}, id2: nil, res: 0},
	{id1: []uint32{}, id2: []uint32{}, res: 0},
	{id1: []uint32{1, 2, 3}, id2: []uint32{1, 2, 3}, res: 0},
	// id1 < id2
	{id1: nil, id2: []uint32{1}, res: -1},
	{id1: []uint32{1, 2}, id2: []uint32{1, 2, 3}, res: -1},
	{id1: []uint32{1, 2, 3}, id2: []uint32{1, 3, 3}, res: -1},
	{id1: []uint32{1, 2, 3, 4}, id2: []uint32{1, 3, 3}, res: -1},
	// id1 > id2
	{id1: []uint32{1}, id2: nil, res: +1},
	{id1: []uint32{1, 2, 3}, id2: []uint32{1, 2}, res: +1},
	{id1: []uint32{1, 3, 3}, id2: []uint32{1, 2, 3}, res: +1},
	{id1: []uint32{1, 3, 3}, id2: []uint32{1, 2, 3, 4}, res: +1},
}

func TestOidCompare(t *testing.T) {
	for i, test := range testDataOidCompare {
		if res := Compare(test.id1, test.id2); res != test.res {
			t.Errorf("TestOidCompare[%d]", i)
		}
	}
}

var testDataOidString = []struct {
	id  []uint32
	res string
}{
	{id: nil, res: ""},
	{id: []uint32{}, res: ""},
	{id: []uint32{1, 2, 3}, res: "1.2.3"},
}

func TestOidString(t *testing.T) {
	for i, test := range testDataOidString {
		if res := String(test.id); res != test.res {
			t.Errorf("TestOidString[%d]", i)
		}
	}
}

var testDataOidParse = []struct {
	str string
	res []uint32
}{
	{str: "", res: []uint32{}},
	{str: ".", res: []uint32{}},
	{str: "1.2.3", res: []uint32{1, 2, 3}},
	{str: ".1.2.3", res: []uint32{1, 2, 3}},
}

func TestOidParse(t *testing.T) {
	for i, test := range testDataOidParse {
		if res := Parse(test.str); !Eq(res, test.res) {
			t.Errorf("TestOidParse[%d]", i)
		}
	}
}
