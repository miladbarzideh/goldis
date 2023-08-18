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

	entry := NewZNode("milad", 20)
	tree.Insert(&entry.tree, avlEntryEq)

	if tree.root != &entry.tree {
		t.Errorf("Expected root to be %v, got %v", tree.root, entry.tree)
	}
}

func TestAVLTree_Remove(t *testing.T) {
	tree := NewAVLTree()
	entry := NewZNode("milad", 20)
	tree.Insert(&entry.tree, avlEntryEq)

	tree.Remove(&entry.tree, avlEntryEq)

	if tree.root != nil {
		t.Errorf("Expected root to be nil, got %v", tree.root)
	}
}

func TestAVLTree_Search(t *testing.T) {
	tree := NewAVLTree()
	entry1 := NewZNode("milad", 20)
	entry2 := NewZNode("ali", 19)
	tree.Insert(&entry1.tree, avlEntryEq)
	tree.Insert(&entry2.tree, avlEntryEq)

	foundNode := tree.Search(&entry2.tree, avlEntryEq)
	if foundNode == nil {
		t.Errorf("Expected found node to be %v, got nil", &entry2.tree)
	}
	if foundNode != &entry2.tree {
		t.Errorf("Expected found node to be %v, got %v", &entry2.tree, foundNode)
	}
}

func TestAVLTree_DisplayNodes(t *testing.T) {
	tree := NewAVLTree()
	nums := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	for range nums {
		n := nums[r.Intn(len(nums))]
		entry := NewZNode("name", float64(n))
		tree.Insert(&entry.tree, avlEntryEq)
	}
	tree.Traverse()
}
