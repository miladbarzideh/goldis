package repository

import (
	"hash/fnv"
	"unsafe"
)

type DataStore struct {
	db *HMap
}

type Entry struct {
	node  HNode
	key   string
	value string
}

func NewDataStore() *DataStore {
	return &DataStore{
		db: NewHMap(),
	}
}

func (ds *DataStore) Get(key string) string {
	entry := Entry{
		node: HNode{hcode: hash(key)},
		key:  key,
	}

	node := ds.db.lookup(&entry.node, entryEq)
	if node == nil {
		return "(nil)"
	}

	return (*Entry)(unsafe.Pointer(uintptr(unsafe.Pointer(node)) - unsafe.Offsetof(Entry{}.node))).value
}

func (ds *DataStore) Set(key string, value string) string {
	entry := Entry{
		node: HNode{hcode: hash(key)},
		key:  key,
	}

	node := ds.db.lookup(&entry.node, entryEq)
	//update the value
	if node != nil {
		(*Entry)(unsafe.Pointer(uintptr(unsafe.Pointer(node)) - unsafe.Offsetof(Entry{}.node))).value = value
	} else {
		entry.value = value
		ds.db.insert(&entry.node)
	}
	return "OK"
}

func entryEq(lhs, rhs *HNode) bool {
	// We can use the unsafe package to perform pointer arithmetic and access members within a structure without knowing its type
	le := (*Entry)(unsafe.Pointer(uintptr(unsafe.Pointer(lhs)) - unsafe.Offsetof(Entry{}.node)))
	re := (*Entry)(unsafe.Pointer(uintptr(unsafe.Pointer(rhs)) - unsafe.Offsetof(Entry{}.node)))
	return lhs.hcode == rhs.hcode && le.key == re.key
}

func hash(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}
