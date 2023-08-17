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
	le := (*AVLEntry)(containerOf(unsafe.Pointer(l), unsafe.Offsetof(AVLEntry{}.node)))
	re := (*AVLEntry)(containerOf(unsafe.Pointer(r), unsafe.Offsetof(AVLEntry{}.node)))
	if le.value > re.value {
		return 1
	} else if le.value < re.value {
		return -1
	}
	return 0
}
