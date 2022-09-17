package oid

import "testing"

func TestMibName(t *testing.T) {
	Name["myenterprise"] = []uint32{1, 3, 6, 1, 4, 1, 9999}

	if !Eq(Name["myenterprise"], []uint32{1, 3, 6, 1, 4, 1, 9999}) {
		t.Error("TestMibName")
	}
}
