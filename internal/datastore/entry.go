package datastore

import (
	"hash/fnv"
	"unsafe"
)

type AVLEntry struct {
	node  AVLNode
	value int32
}

func NewAVLEntry(value int32) *AVLEntry {
	return &AVLEntry{
		value: value,
	}
}

func avlEntryEq(l, r *AVLNode) int {
	le := avlEntryContainerOf(l)
	re := avlEntryContainerOf(r)
	if le.value > re.value {
		return 1
	} else if le.value < re.value {
		return -1
	}
	return 0
}

// avlEntryContainerOf to have an intrusive data structure
func avlEntryContainerOf(node *AVLNode) *AVLEntry {
	return (*AVLEntry)(unsafe.Pointer(uintptr(unsafe.Pointer(node)) - unsafe.Offsetof(AVLEntry{}.node)))
}

type MapEntry struct {
	node  HNode
	key   string
	value string
}

func NewMapEntry(key string) *MapEntry {
	return &MapEntry{
		node: HNode{hcode: hash(key)},
		key:  key,
	}
}

// We can use the unsafe package to perform pointer arithmetic,
// and have an intrusive data structure
func mapEntryContainerOf(lhs *HNode) *MapEntry {
	return (*MapEntry)(unsafe.Pointer(uintptr(unsafe.Pointer(lhs)) - unsafe.Offsetof(MapEntry{}.node)))
}

func hash(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}
