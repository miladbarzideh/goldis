package datastore

import (
	"fmt"
	"log"
	"strings"
	"unsafe"

	"github.com/miladbarzideh/goldis/utils"
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
	return (*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node))).value
}

func (ds *DataStore) Set(key string, value string) string {
	entry := NewMapEntry(key, STR)
	node := ds.db.Lookup(&entry.node, EntryEq)
	//update the value
	if node != nil {
		(*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node))).value = value
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
	res := strings.Builder{}
	for i, node := range nodes {
		entry := (*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
		kv := fmt.Sprintf("%v) %s => %s\n", i+1, entry.key, entry.value)
		log.Printf(kv)
		res.WriteString(kv)
	}
	return res.String()
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
		entry = (*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
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
	znodes := entry.zset.Query(score, name, offset, limit)
	result := strings.Builder{}
	for i, znode := range znodes {
		result.WriteString(fmt.Sprintf("%v) %v => %v\n", i, znode.name, znode.score))
	}
	return result.String()
}

func (ds *DataStore) ZShow(key string) string {
	exist, entry := ds.expect(key)
	if !exist {
		return resNil
	}
	return entry.zset.Show()
}

func (ds *DataStore) expect(key string) (bool, *MapEntry) {
	entry := NewMapEntry(key, ZSET)
	node := ds.db.Lookup(&entry.node, EntryEq)
	if node == nil {
		return false, nil
	}
	entry = (*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
	if entry.entryType != ZSET {
		return false, nil
	}
	return true, entry
}

type EntryType int

const (
	ZSET = iota
	STR
)

type MapEntry struct {
	node      HNode
	zset      *ZSet
	key       string
	value     string
	entryType EntryType
}

func NewMapEntry(key string, entryType EntryType) *MapEntry {
	return &MapEntry{
		node:      HNode{hcode: utils.Hash(key)},
		key:       key,
		entryType: entryType,
	}
}

func EntryEq(lhs, rhs *HNode) bool {
	le := (*MapEntry)(utils.ContainerOf(unsafe.Pointer(lhs), unsafe.Offsetof(MapEntry{}.node)))
	re := (*MapEntry)(utils.ContainerOf(unsafe.Pointer(rhs), unsafe.Offsetof(MapEntry{}.node)))
	return lhs.hcode == rhs.hcode && le.key == re.key
}
