package utils

import (
	"hash/fnv"
	"unsafe"
)

// ContainerOf We can use the unsafe package to perform pointer arithmetic,
// and have an intrusive data structure
func ContainerOf(ptr unsafe.Pointer, memberOffset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) - memberOffset)
}

func Hash(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

func Max(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
