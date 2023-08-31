package datastore

import (
	"strconv"
	"testing"

	"github.com/miladbarzideh/goldis/utils"
)

func TestNewHMap(t *testing.T) {
	hmap := NewHMap(MapEntryComparator)

	if hmap.tab1.size != 0 {
		t.Errorf("Expected tab1 size to be 0, got %d", hmap.tab1.size)
	}
	if hmap.tab2.tab != nil {
		t.Error("Expected tab2 to be nil")
	}
}

func TestHMapInsert(t *testing.T) {
	hmap := NewHMap(MapEntryComparator)

	for i := 0; i < 3; i++ {
		node := &HNode{hcode: utils.Hash("key")}
		hmap.Insert(node)
		expectedSize := i + 1
		if hmap.tab1.size != expectedSize {
			t.Errorf("Expected tab1 size to be %d, got %d", expectedSize, hmap.tab1.size)
		}
	}
}

func TestHMapLookup(t *testing.T) {
	hmap := NewHMap(MapEntryComparator)
	node := &HNode{hcode: utils.Hash("key")}
	hmap.Insert(node)

	foundNode := hmap.Lookup(node)

	if foundNode == nil {
		t.Errorf("Expected node to be found, got  nil")
	}
	if foundNode != node {
		t.Errorf("Expected node to be %v, got %v", foundNode, node)
	}
}

func TestHMapDelete(t *testing.T) {
	hmap := NewHMap(MapEntryComparator)
	node := &HNode{hcode: utils.Hash("key")}
	hmap.Insert(node)

	deleteNode := hmap.Pop(node)

	if hmap.tab1.size != 0 {
		t.Errorf("Expected tab1 size to be 0, got %v", hmap.tab1.size)
	}
	if deleteNode != node {
		t.Errorf("Expected deleted node to be %v, got %v", deleteNode, node)
	}
	found := hmap.Lookup(node)
	if found != nil {
		t.Errorf("Expected found to be nil, got %v", found)
	}
}

func TestHMaoDestroy(t *testing.T) {
	hmap := NewHMap(MapEntryComparator)
	node := &HNode{hcode: utils.Hash("key")}
	hmap.Insert(node)

	hmap.Destroy()

	if hmap.tab1.size != 0 {
		t.Errorf("Expected tab1 size to be 0, got %v", hmap.tab1.size)
	}
	if hmap.tab2.tab != nil {
		t.Error("Expected tab2 to be nil")
	}
}

func TestHMapResizing(t *testing.T) {
	hmap := NewHMap(MapEntryComparator)

	for i := 0; i < 10; i++ {
		node := &HNode{hcode: utils.Hash("key" + strconv.Itoa(i))}
		hmap.Insert(node)
	}

	if hmap.tab1.size != 10 {
		t.Errorf("Expected tab1 size to be 10, got %v", hmap.tab1.size)
	}
}
