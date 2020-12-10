package iunsafe

import (
	"reflect"
	"unsafe"
)

type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type interfaceHeader struct {
	Table unsafe.Pointer
	Data unsafe.Pointer
}

func ByteSliceAsInt64Slice(i []byte) []int64 {
	var o []int64

	hdrI := (*sliceHeader)(unsafe.Pointer(&i))
	hdrO := (*sliceHeader)(unsafe.Pointer(&o))
	hdrO.Data = hdrI.Data
	hdrO.Cap = hdrI.Cap / 8
	hdrO.Len = hdrI.Len / 8
	return o
}

func Int64SliceAsByteSlice(i []int64) []byte {
	var o []byte

	hdrI := (*sliceHeader)(unsafe.Pointer(&i))
	hdrO := (*sliceHeader)(unsafe.Pointer(&o))
	hdrO.Data = hdrI.Data
	hdrO.Cap = hdrI.Cap * 8
	hdrO.Len = hdrI.Len * 8
	return o
}

func Uint64SliceAsByteSlice(i []uint64) []byte {
	var o []byte

	hdrI := (*sliceHeader)(unsafe.Pointer(&i))
	hdrO := (*sliceHeader)(unsafe.Pointer(&o))
	hdrO.Data = hdrI.Data
	hdrO.Cap = hdrI.Cap * 8
	hdrO.Len = hdrI.Len * 8
	return o
}

func ByteSliceAsArbSlice(i []byte, interfaceO interface{}) {
	t := reflect.TypeOf(interfaceO)
	if t.Kind() != reflect.Ptr {
		panic("not a ptr")
	}
	t = t.Elem()
	if t.Kind() != reflect.Slice {
		panic("not a ptr slice")
	}
	t = t.Elem()
	// This might be dumb.
	//if t.Kind() != reflect.Struct {
	//	panic("not a ptr slice struct")
	//}
	actualS := (*sliceHeader)((*interfaceHeader)(unsafe.Pointer(&interfaceO)).Data)
	hdrI := (*sliceHeader)(unsafe.Pointer(&i))
	actualS.Data = hdrI.Data
	// This is deliberate so that attempting to append to a slice will force a reallocation.
	// That may cause other problems with memory leaks.  Caveat emptor.
	inputLen := hdrI.Len / int(t.Size())
	actualS.Len = inputLen
	actualS.Cap = inputLen
}
