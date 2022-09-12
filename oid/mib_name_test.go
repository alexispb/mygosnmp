package oid

import "testing"

func TestMibName(t *testing.T) {
	MibName["myenterprise"] = []uint32{1, 3, 6, 1, 4, 1, 9999}

	if !Eq(MibName["myenterprise"], []uint32{1, 3, 6, 1, 4, 1, 9999}) {
		t.Error("TestMibName")
	}
}
