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

func TestAVLTree_RightRotate(t *testing.T) {
	tree := NewAVLTree()
	entry1 := NewZNode("milad", 3)
	entry2 := NewZNode("ali", 2)
	entry3 := NewZNode("ali", 1)
	tree.Insert(&entry1.tree, avlEntryEq)
	tree.Insert(&entry2.tree, avlEntryEq)
	tree.Insert(&entry3.tree, avlEntryEq)

	h := tree.root.getHeight()
	if h != 1 {
		t.Errorf("Expected height to be 1, got %v", h)
	}
	if tree.root.left == nil {
		t.Errorf("Expected left child not to be null")
	}
	if tree.root.right == nil {
		t.Errorf("Expected ritgh child not to be null")
	}
}

func TestAVLTree_LeftRotate(t *testing.T) {
	tree := NewAVLTree()
	entry1 := NewZNode("milad", 1)
	entry2 := NewZNode("ali", 2)
	entry3 := NewZNode("ali", 3)
	tree.Insert(&entry1.tree, avlEntryEq)
	tree.Insert(&entry2.tree, avlEntryEq)
	tree.Insert(&entry3.tree, avlEntryEq)

	h := tree.root.getHeight()
	if h != 1 {
		t.Errorf("Expected height to be 1, got %v", h)
	}
	if tree.root.left == nil {
		t.Errorf("Expected left child not to be null")
	}
	if tree.root.right == nil {
		t.Errorf("Expected ritgh child not to be null")
	}
}

func TestAVLTree_DoubleLeftRotate(t *testing.T) {
	tree := NewAVLTree()
	entry1 := NewZNode("milad", 7)
	entry2 := NewZNode("ali", 9)
	entry3 := NewZNode("ali", 6)
	tree.Insert(&entry1.tree, avlEntryEq)
	tree.Insert(&entry2.tree, avlEntryEq)
	tree.Insert(&entry3.tree, avlEntryEq)

	h := tree.root.getHeight()
	if h != 1 {
		t.Errorf("Expected height to be 1, got %v", h)
	}
	if tree.root.left == nil {
		t.Errorf("Expected left child not to be null")
	}
	if tree.root.right == nil {
		t.Errorf("Expected ritgh child not to be null")
	}
}

func TestAVLTree_DoubleRightRotate(t *testing.T) {
	tree := NewAVLTree()
	entry1 := NewZNode("milad", 9)
	entry2 := NewZNode("ali", 5)
	entry3 := NewZNode("ali", 7)
	tree.Insert(&entry1.tree, avlEntryEq)
	tree.Insert(&entry2.tree, avlEntryEq)
	tree.Insert(&entry3.tree, avlEntryEq)

	h := tree.root.getHeight()
	if h != 1 {
		t.Errorf("Expected height to be 1, got %v", h)
	}
	if tree.root.left == nil {
		t.Errorf("Expected left child not to be null")
	}
	if tree.root.right == nil {
		t.Errorf("Expected ritgh child not to be null")
	}
}

func TestAVLTree_RemoveRoot(t *testing.T) {
	tree := NewAVLTree()
	entry := NewZNode("milad", 20)
	tree.Insert(&entry.tree, avlEntryEq)

	tree.Remove(&entry.tree, avlEntryEq)

	if tree.root != nil {
		t.Errorf("Expected root to be nil, got %v", tree.root)
	}
}

func TestAVLTree_RemoveOneChild(t *testing.T) {
	tree := NewAVLTree()
	entry1 := NewZNode("milad", 20)
	entry2 := NewZNode("milad", 18)
	tree.Insert(&entry1.tree, avlEntryEq)
	tree.Insert(&entry2.tree, avlEntryEq)

	tree.Remove(&entry1.tree, avlEntryEq)

	if tree.root.getHeight() != 0 {
		t.Errorf("Expected height to be 0, got %v", tree.root.getHeight())
	}
}

func TestAVLTree_RemoveTwoChild(t *testing.T) {
	tree := NewAVLTree()
	entry1 := NewZNode("milad", 20)
	entry2 := NewZNode("milad", 19)
	entry3 := NewZNode("milad", 18)
	tree.Insert(&entry1.tree, avlEntryEq)
	tree.Insert(&entry2.tree, avlEntryEq)
	tree.Insert(&entry3.tree, avlEntryEq)

	tree.Remove(&entry2.tree, avlEntryEq)

	foundNode := tree.Search(&entry2.tree, avlEntryEq)
	if foundNode != nil {
		t.Errorf("Expected foundNode to be nil, got %v", foundNode)
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
