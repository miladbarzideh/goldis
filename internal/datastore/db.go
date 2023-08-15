package datastore

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
	return mapEntryContainerOf(node).value
}

func (ds *DataStore) Set(key string, value string) string {
	entry := NewMapEntry(key)
	node := ds.db.Lookup(&entry.node, EntryEq)
	//update the value
	if node != nil {
		mapEntryContainerOf(node).value = value
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
