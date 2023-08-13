package datastore

import (
	"math/rand"
	"testing"
	"time"
)

func TestNewAVLTree(t *testing.T) {
	tree := NewAVLTree()

	if tree.root != nil {
		t.Errorf("Expected root to be nil, got %v", tree.root)
	}
}

func TestAVLTree_Insert(t *testing.T) {
	tree := NewAVLTree()

	entry := NewAVLEntry(1)
	tree.Insert(&entry.node, avlEntryEq)

	if tree.root != &entry.node {
		t.Errorf("Expected root to be %v, got %v", tree.root, entry.node)
	}
}

func TestAVLTree_Remove(t *testing.T) {
	tree := NewAVLTree()
	entry := NewAVLEntry(1)
	tree.Insert(&entry.node, avlEntryEq)

	tree.Remove(&entry.node, avlEntryEq)

	if tree.root != nil {
		t.Errorf("Expected root to be nil, got %v", tree.root)
	}
}

func TestAVLTree_Search(t *testing.T) {
	tree := NewAVLTree()
	entry1 := NewAVLEntry(1)
	entry2 := NewAVLEntry(8)
	tree.Insert(&entry1.node, avlEntryEq)
	tree.Insert(&entry2.node, avlEntryEq)

	foundNode := tree.Search(&entry2.node, avlEntryEq)
	if foundNode == nil {
		t.Errorf("Expected found node to be %v, got nil", &entry2.node)
	}
	if foundNode != &entry2.node {
		t.Errorf("Expected found node to be %v, got %v", &entry2.node, foundNode)
	}
}

func TestAVLTree_DisplayNodes(t *testing.T) {
	tree := NewAVLTree()
	nums := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	for range nums {
		n := nums[r.Intn(len(nums))]
		entry := *NewAVLEntry(n)
		tree.Insert(&entry.node, avlEntryEq)
	}

	tree.DisplayNodes()
}
