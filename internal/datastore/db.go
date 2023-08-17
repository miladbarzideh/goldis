package datastore

import (
	"fmt"
	"strings"
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

// ZAdd command pattern: zadd zset score name
func (ds *DataStore) ZAdd(key string, score float64, name string) string {
	entry := NewMapEntry(key)
	node := ds.db.Lookup(&entry.node, EntryEq)
	//update the value
	if node == nil {
		entry.zset = NewZSet()
		ds.db.Insert(&entry.node)
	} else {
		entry = (*MapEntry)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
		//check the type of data and raise an error
	}

	entry.zset.Add(name, score)
	return resOK
}

// ZRemove command pattern: zrem zset name
func (ds *DataStore) ZRemove(key string, name string) string {
	exist, entry := ds.expect(key)
	if !exist {
		return resKO
	}

	node := entry.zset.Pop(name)
	if node != nil {
		node = &ZNode{}
	}
	return resOK
}

// ZScore command pattern: zscore zset name
func (ds *DataStore) ZScore(key string, name string) string {
	exist, entry := ds.expect(key)
	if !exist {
		return resNil
	}
	node := entry.zset.Lookup(name)
	if node == nil {
		return resNil
	}
	return fmt.Sprintf("%v", node.score)
}

// ZQuery command pattern: zquery zset score name offset limit
func (ds *DataStore) ZQuery(key string, score float64, name string, offset uint32, limit uint32) string {
	exist, entry := ds.expect(key)
	if !exist {
		return resNil
	}
	node := entry.zset.Query(score, name, offset)
	result := strings.Builder{}
	n := uint32(0)
	for node != nil && n < limit {
		result.WriteString(fmt.Sprintf("%v => %v", node.name, node.score))
		node = (*ZNode)(containerOf(unsafe.Pointer(node.tree.offset(1)), unsafe.Offsetof(ZNode{}.tree)))
		n += 2
	}
	return result.String()
}

func (ds *DataStore) expect(key string) (bool, *MapEntry) {
	entry := NewMapEntry(key)
	node := ds.db.Lookup(&entry.node, EntryEq)
	if node == nil {
		return false, nil
	}
	entry = (*MapEntry)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
	//TODO: check the type of data and raise an error
	return true, entry
}

type MapEntry struct {
	node  HNode
	zset  *ZSet //which one GOD!! :)) TODO: remove it :))
	key   string
	value string
	//also add a type
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
