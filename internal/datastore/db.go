package datastore

import "unsafe"

const (
	resOK  = "OK"
	resKO  = "KO"
	resNil = "(nil)"
)

type DataStore struct {
	db *HMap
}

func NewDataStore() *DataStore {
	return &DataStore{
		db: NewHMap(),
	}
}

func (ds *DataStore) Get(key string) string {
	entry := NewMapEntry(key)
	node := ds.db.Lookup(&entry.node, EntryEq)
	if node == nil {
		return resNil
	}
	return (*MapEntry)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node))).value
}

func (ds *DataStore) Set(key string, value string) string {
	entry := NewMapEntry(key)
	node := ds.db.Lookup(&entry.node, EntryEq)
	//update the value
	if node != nil {
		(*MapEntry)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node))).value = value
	} else {
		entry.value = value
		ds.db.Insert(&entry.node)
	}
	return resOK
}

func (ds *DataStore) Delete(key string) string {
	entry := NewMapEntry(key)
	node := ds.db.Pop(&entry.node, EntryEq)
	if node != nil {
		//containerOf(node) = nil
		return resOK
	}
	return resKO
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

func EntryEq(lhs, rhs *HNode) bool {
	le := (*MapEntry)(containerOf(unsafe.Pointer(lhs), unsafe.Offsetof(MapEntry{}.node)))
	re := (*MapEntry)(containerOf(unsafe.Pointer(rhs), unsafe.Offsetof(MapEntry{}.node)))
	return lhs.hcode == rhs.hcode && le.key == re.key
}
