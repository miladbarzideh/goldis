package datastore

import (
	"hash/fnv"
	"unsafe"
)

// We can use the unsafe package to perform pointer arithmetic,
// and have an intrusive data structure
func containerOf(ptr unsafe.Pointer, memberOffset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) - memberOffset)
}

func hash(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

func max(a uint32, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}
