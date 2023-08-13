package datastore

import (
	"hash/fnv"
	"unsafe"
)

const (
	resOK  = "OK"
	resKO  = "KO"
	resNil = "(nil)"
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

func newEntry(key string) *Entry {
	return &Entry{
		node: HNode{hcode: hash(key)},
		key:  key,
	}
}

func (ds *DataStore) Get(key string) string {
	entry := newEntry(key)
	node := ds.db.Lookup(&entry.node, entryEq)
	if node == nil {
		return resNil
	}
	return containerOf(node).value
}

func (ds *DataStore) Set(key string, value string) string {
	entry := newEntry(key)
	node := ds.db.Lookup(&entry.node, entryEq)
	//update the value
	if node != nil {
		containerOf(node).value = value
	} else {
		entry.value = value
		ds.db.Insert(&entry.node)
	}
	return resOK
}

func (ds *DataStore) Delete(key string) string {
	entry := newEntry(key)
	node := ds.db.Pop(&entry.node, entryEq)
	if node != nil {
		//containerOf(node) = nil
		return resOK
	}
	return resKO
}

func entryEq(lhs, rhs *HNode) bool {
	le := containerOf(lhs)
	re := containerOf(rhs)
	return lhs.hcode == rhs.hcode && le.key == re.key
}

// We can use the unsafe package to perform pointer arithmetic,
// and have an intrusive data structure
func containerOf(lhs *HNode) *Entry {
	return (*Entry)(unsafe.Pointer(uintptr(unsafe.Pointer(lhs)) - unsafe.Offsetof(Entry{}.node)))
}

func hash(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}
