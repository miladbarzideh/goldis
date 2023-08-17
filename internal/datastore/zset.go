package datastore

import (
	"unsafe"
)

type ZSet struct {
	hmap *HMap
	tree *AVLTree
}

type ZNode struct {
	hmap  HNode
	tree  AVLNode
	score float64
	name  string
}

func NewZSet() *ZSet {
	return &ZSet{
		hmap: NewHMap(),
		tree: NewAVLTree(),
	}
}

func NewZNode(name string, score float64) *ZNode {
	return &ZNode{
		hmap:  HNode{hcode: hash(name)},
		tree:  AVLNode{},
		score: score,
		name:  name,
	}
}

func (zset *ZSet) Add(name string, score float64) bool {
	node := zset.Lookup(name)
	node = NewZNode(name, score)
	zset.hmap.Insert(&node.hmap)
	zset.tree.Insert(&node.tree, avlEntryEq)
	return true //TODO: always return true, return false in case of update
}

func (zset *ZSet) Lookup(name string) *ZNode {
	if zset.tree == nil {
		return nil
	}
	key := newHKey(name)
	found := zset.hmap.Lookup(&key.node, entryEq)
	if found == nil {
		return nil
	}
	return (*ZNode)(containerOf(unsafe.Pointer(found), unsafe.Offsetof(ZNode{}.hmap)))
}

func (zset *ZSet) Pop(name string) *ZNode {
	if zset.tree == nil {
		return nil
	}
	key := newHKey(name)
	found := zset.hmap.Pop(&key.node, entryEq)
	if found == nil {
		return nil
	}

	node := (*ZNode)(containerOf(unsafe.Pointer(found), unsafe.Offsetof(ZNode{}.hmap)))
	zset.tree.Remove(&node.tree, avlEntryEq)
	return node
}

func (zset *ZSet) Query(score float64, name string, offset uint32) *ZNode {
	var found *AVLNode
	cur := zset.tree.root
	for cur != nil {
		if zless(cur, score, name) {
			cur = cur.right
		} else {
			found = cur
			cur = cur.left
		}
	}

	if found != nil {
		found = zset.tree.Offset(found, offset)
		return (*ZNode)(containerOf(unsafe.Pointer(found), unsafe.Offsetof(ZNode{}.tree)))
	}
	return nil
}

func (zset *ZSet) Dispose() {
	zset.hmap.Destroy()
	zset.tree.Dispose()
}

// zless compare by the (score, name) tuple
func zless(lhs *AVLNode, score float64, name string) bool {
	zl := (*ZNode)(containerOf(unsafe.Pointer(lhs), unsafe.Offsetof(ZNode{}.tree)))
	if zl.score != score {
		return zl.score < score
	}
	return zl.name < name
}

func less(lhs *AVLNode, rhs *AVLNode) bool {
	zr := (*ZNode)(containerOf(unsafe.Pointer(rhs), unsafe.Offsetof(ZNode{}.tree)))
	return zless(lhs, zr.score, zr.name)
}

// a helper structure for the hashtable lookup (TODO: same as hmapEntry)
type HKey struct {
	node HNode
	name string
}

func newHKey(name string) *HKey {
	return &HKey{
		node: HNode{hcode: hash(name)},
		name: name,
	}
}

func entryEq(node, key *HNode) bool {
	if node.hcode != key.hcode {
		return false
	}
	znode := (*ZNode)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(ZNode{}.hmap)))
	hkey := (*HKey)(containerOf(unsafe.Pointer(key), unsafe.Offsetof(HKey{}.node)))
	return znode.name == hkey.name
}
