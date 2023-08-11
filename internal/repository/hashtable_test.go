package repository

import "testing"

func TestNewHMap(t *testing.T) {
	hmap := NewHMap()

	if hmap.tab1.size != 0 {
		t.Errorf("Expected tab1 size to be 0, got %d", hmap.tab1.size)
	}
	if hmap.tab2.tab != nil {
		t.Error("Expected tab2 to be nil")
	}
}

func TestHMapInsert(t *testing.T) {
	hmap := NewHMap()

	for i := 0; i < 3; i++ {
		node := &HNode{hcode: hash("key")}
		hmap.insert(node)
		expectedSize := i + 1
		if hmap.tab1.size != expectedSize {
			t.Errorf("Expected tab1 size to be %d, got %d", expectedSize, hmap.tab1.size)
		}
	}
}

func TestHMapLookup(t *testing.T) {
	hmap := NewHMap()
	node := &HNode{hcode: hash("key")}
	hmap.insert(node)

	foundNode := hmap.lookup(node, entryEq)

	if foundNode == nil {
		t.Errorf("Expected node to be found, got  nil")
	}
	if foundNode != node {
		t.Errorf("Expected node to be %v, got %v", foundNode, node)
	}
}

func TestHMapDelete(t *testing.T) {
	hmap := NewHMap()
	node := &HNode{hcode: hash("key")}
	hmap.insert(node)

	deleteNode := hmap.pop(node, entryEq)

	if hmap.tab1.size != 0 {
		t.Errorf("Expected tab1 size to be 0, got %v", hmap.tab1.size)
	}
	if deleteNode != node {
		t.Errorf("Expected deleted node to be %v, got %v", deleteNode, node)
	}
}

func TestHMaoDestroy(t *testing.T) {
	hmap := NewHMap()
	node := &HNode{hcode: hash("key")}
	hmap.insert(node)

	hmap.destroy()

	if hmap.tab1.size != 0 {
		t.Errorf("Expected tab1 size to be 0, got %v", hmap.tab1.size)
	}
	if hmap.tab2.tab != nil {
		t.Error("Expected tab2 to be nil")
	}
}
