package datastore

import (
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
