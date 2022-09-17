package internal

/*
import (
	"reflect"
	"unsafe"
)

func StringAsSlice(s string) (b []byte) {
	hs := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hb := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	hb.Data = hs.Data
	hb.Len = hs.Len
	hb.Cap = hs.Len
	return
}

func SliceAsString(b []byte) (s string) {
	hb := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hs := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hs.Data = hb.Data
	hs.Len = hb.Len
	return
}
*/
