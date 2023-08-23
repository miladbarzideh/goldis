package datastore

import (
	"fmt"
	"log"
	"strings"
	"unsafe"
)

const (
	resOK   = "OK"
	resKO   = "KO"
	resNil  = "(nil)"
	errType = "(error) ERR expect zset"
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
	entry := NewMapEntry(key, STR)
	node := ds.db.Lookup(&entry.node, EntryEq)
	if node == nil {
		return resNil
	}
	return (*MapEntry)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node))).value
}

func (ds *DataStore) Set(key string, value string) string {
	entry := NewMapEntry(key, STR)
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
	entry := NewMapEntry(key, STR)
	node := ds.db.Pop(&entry.node, EntryEq)
	if node != nil {
		//containerOf(node) = nil
		return resOK
	}
	return resKO
}

func (ds *DataStore) Keys() string {
	nodes := ds.db.Keys()
	log.Print("Hashtable key-value pairs:\n")
	for _, node := range nodes {
		entry := (*MapEntry)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
		log.Printf("%s => %s", entry.key, entry.value)
	}
	return resOK
}

// ZAdd command pattern: zadd zset score name
func (ds *DataStore) ZAdd(key string, score float64, name string) string {
	entry := NewMapEntry(key, ZSET)
	node := ds.db.Lookup(&entry.node, EntryEq)
	//update the value
	if node == nil {
		entry.zset = NewZSet()
		ds.db.Insert(&entry.node)
	} else {
		entry = (*MapEntry)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
		if entry.entryType != ZSET {
			return errType
		}
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
func (ds *DataStore) ZQuery(key string, score float64, name string, offset int32, limit uint32) string {
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

func (ds *DataStore) ZShow(key string) string {
	exist, entry := ds.expect(key)
	if !exist {
		return resNil
	}
	entry.zset.Show()
	return resOK
}

func (ds *DataStore) expect(key string) (bool, *MapEntry) {
	entry := NewMapEntry(key, ZSET)
	node := ds.db.Lookup(&entry.node, EntryEq)
	if node == nil {
		return false, nil
	}
	entry = (*MapEntry)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
	//TODO: check the type of data and raise an error
	return true, entry
}

type EntryType int

const (
	ZSET = iota
	STR
)

type MapEntry struct {
	node      HNode
	zset      *ZSet //TODO
	key       string
	value     string
	entryType EntryType
}

func NewMapEntry(key string, entryType EntryType) *MapEntry {
	return &MapEntry{
		node:      HNode{hcode: hash(key)},
		key:       key,
		entryType: entryType,
	}
}

func EntryEq(lhs, rhs *HNode) bool {
	le := (*MapEntry)(containerOf(unsafe.Pointer(lhs), unsafe.Offsetof(MapEntry{}.node)))
	re := (*MapEntry)(containerOf(unsafe.Pointer(rhs), unsafe.Offsetof(MapEntry{}.node)))
	return lhs.hcode == rhs.hcode && le.key == re.key
}
