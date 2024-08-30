package i64adder

import (
	"unsafe"
)

//go:linkname getm runtime.getm
func getm() uintptr

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

//go:noescape
//go:linkname memhash32 runtime.memhash32
func memhash32(p unsafe.Pointer, h uintptr) uintptr

func hash() uint64 {
	m := getm()
	return uint64(memhash(unsafe.Pointer(&m), 0, unsafe.Sizeof(m)))
}
