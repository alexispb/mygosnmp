package internal

import "unsafe"

var IsLittleEndian bool

func defineEndian() {
	var i int32 = 0x01020304
	pu := unsafe.Pointer(&i)
	pb := (*byte)(pu)
	IsLittleEndian = (*pb == 0x04)
}

func init() {
	defineEndian()
}
