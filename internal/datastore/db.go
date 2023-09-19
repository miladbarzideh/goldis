package datastore

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unsafe"

	"github.com/miladbarzideh/goldis/utils"
)

const (
	resOK   = "OK"
	resKO   = "KO"
	resNil  = "(nil)"
	errType = "(error) ERR expect zset"
)

const (
	maxWorks           = 200
	largeContainerSize = 10000
)

type DataStore struct {
	db   *HMap
	heap *MinHeap
}

func NewDataStore() *DataStore {
	return &DataStore{
		db:   NewHMap(MapEntryComparator),
		heap: NewMinHeap(),
	}
}

func (ds *DataStore) Get(key string) string {
	entry := NewMapEntry(key, STR)
	node := ds.db.Lookup(&entry.node)
	if node == nil {
		return resNil
	}
	return (*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node))).value
}

func (ds *DataStore) Set(key string, value string) string {
	entry := NewMapEntry(key, STR)
	node := ds.db.Lookup(&entry.node)
	// update the value
	if node != nil {
		(*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node))).value = value
	} else {
		entry.value = value
		ds.db.Insert(&entry.node)
	}
	return resOK
}

func (ds *DataStore) Delete(key string) string {
	entry := NewMapEntry(key, ZSET)
	node := ds.db.Pop(&entry.node)
	if node != nil {
		// containerOf(node) = nil
		entry := (*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
		ds.setEntryTtl(entry, -1)
		entryDel(entry)
		return resOK
	}
	return resKO
}

func entryDel(entry *MapEntry) {
	if entry.entryType == ZSET {
		if entry.zset.hmap.Size() > largeContainerSize { //too big
			log.Println("Perform async action to delete entry")
			utils.GetThreadPoolInstance().ThreadPoolQueue(entryDelAsync(), entry)
		} else {
			entry.zset.Dispose()
		}
	}
}

func entryDelAsync() func(arg interface{}) {
	return func(arg interface{}) {
		entry := arg.(*MapEntry)
		entry.zset.Dispose()
	}
}

func (ds *DataStore) Keys() string {
	nodes := ds.db.Keys()
	log.Print("Hashtable key-value pairs:\n")
	res := strings.Builder{}
	for i, node := range nodes {
		entry := (*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
		value := entry.value
		if entry.entryType == ZSET {
			value = "ZSET"
		}
		kv := fmt.Sprintf("%v) %s => %s\n", i+1, entry.key, value)
		log.Print(kv)
		res.WriteString(kv)
	}
	return res.String()
}

// ZAdd command pattern: zadd zset score name
func (ds *DataStore) ZAdd(key string, score float64, name string) string {
	entry := NewMapEntry(key, ZSET)
	node := ds.db.Lookup(&entry.node)
	// update the value
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

	entry.zset.Pop(name)
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

func (ds *DataStore) Expire(key string, ttl int64) string {
	entry := NewMapEntry(key, STR)
	node := ds.db.Lookup(&entry.node)
	if node == nil {
		return resNil
	}
	entry = (*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
	ds.setEntryTtl(entry, ttl)
	return resOK
}

func (ds *DataStore) Ttl(key string) string {
	entry := NewMapEntry(key, STR)
	node := ds.db.Lookup(&entry.node)
	if node == nil {
		return resNil
	}
	entry = (*MapEntry)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(MapEntry{}.node)))
	heapIndex := entry.heapIndex
	if heapIndex == -1 {
		return resNil
	}
	item := ds.heap.Get(entry.heapIndex)
	now := time.Now().UnixMilli()
	expireAt := int64(0)
	if item.value > now {
		expireAt = item.value - now
	}
	return fmt.Sprintf("(int) %v", expireAt)
}

func (ds *DataStore) setEntryTtl(entry *MapEntry, ttl int64) {
	if ttl < 0 && entry.heapIndex != -1 {
		ds.heap.Remove(entry.heapIndex)
		entry.heapIndex = -1
	} else if ttl > 0 {
		now := time.Now().UnixMilli()
		expireTime := now + ttl
		if entry.heapIndex == -1 {
			item := HeapItem{value: expireTime, ref: &entry.heapIndex}
			ds.heap.Insert(item)
		} else {
			ds.heap.Update(entry.heapIndex, expireTime)
		}
	}
}

func (ds *DataStore) RemoveExpiredKeys() {
	now := time.Now().UnixMilli()
	works := 0
	for ds.heap.Get(0) != nil && ds.heap.Get(0).value < now {
		ref := ds.heap.Get(0).ref
		entry := (*MapEntry)(utils.ContainerOf(unsafe.Pointer(ref), unsafe.Offsetof(MapEntry{}.heapIndex)))
		ds.heap.Remove(0)
		ds.db.Pop(&entry.node)
		if works > maxWorks {
			// don't stall the server if too many keys are expiring at once
			break
		}
		works++
	}
}

func (ds *DataStore) expect(key string) (bool, *MapEntry) {
	entry := NewMapEntry(key, ZSET)
	node := ds.db.Lookup(&entry.node)
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
	heapIndex int32
}

func NewMapEntry(key string, entryType EntryType) *MapEntry {
	return &MapEntry{
		node:      HNode{hcode: utils.Hash(key)},
		key:       key,
		entryType: entryType,
		heapIndex: -1,
	}
}
