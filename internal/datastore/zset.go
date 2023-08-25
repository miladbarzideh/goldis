package datastore

import (
	"fmt"
	"log"
	"strings"
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
		tree:  AVLNode{count: 1},
		score: score,
		name:  name,
	}
}

func (zset *ZSet) Add(name string, score float64) bool {
	node := zset.Lookup(name)
	if node == nil {
		node = NewZNode(name, score)
		zset.hmap.Insert(&node.hmap)
		zset.tree.Insert(&node.tree, avlEntryEq)
		return true
	} else {
		zset.update(node, score)
		return false
	}
}

func (zset *ZSet) update(node *ZNode, score float64) {
	if node.score == score {
		return
	}
	zset.tree.Remove(&node.tree, avlEntryEq)
	key := newHKey(node.name)
	zset.hmap.Pop(&key.node, entryEq)
	newNode := NewZNode(node.name, score)
	zset.hmap.Insert(&newNode.hmap)
	zset.tree.Insert(&newNode.tree, avlEntryEq)
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

func (zset *ZSet) Show() string {
	printHashtable(zset.hmap.Keys())
	return printTreeNode(zset.tree.Traverse())
}

func printHashtable(nodes []*HNode) {
	log.Print("Hashtable name-score pair:\n")
	for _, node := range nodes {
		entry := (*ZNode)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(ZNode{}.hmap)))
		log.Printf("%v => %v", entry.name, entry.score)
	}
}

func printTreeNode(nodes []*AVLNode) string {
	log.Print("AVL Tree Inorder Traversal:\n")
	res := strings.Builder{}
	for i, node := range nodes {
		entry := (*ZNode)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(ZNode{}.tree)))
		sn := fmt.Sprintf("%v) %v => %v\n", i+1, entry.score, entry.name)
		log.Printf(sn)
		res.WriteString(sn)
	}
	return res.String()
}

func (zset *ZSet) Query(score float64, name string, offset int32) *ZNode {
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

// HKey a helper structure for the hashtable lookup
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

func entryEq(key, node *HNode) bool {
	if node.hcode != key.hcode {
		return false
	}
	znode := (*ZNode)(containerOf(unsafe.Pointer(node), unsafe.Offsetof(ZNode{}.hmap)))
	hkey := (*HKey)(containerOf(unsafe.Pointer(key), unsafe.Offsetof(HKey{}.node)))
	return znode.name == hkey.name
}

// zless compare by the (score, name) tuple
func zless(lhs *AVLNode, score float64, name string) bool {
	zl := (*ZNode)(containerOf(unsafe.Pointer(lhs), unsafe.Offsetof(ZNode{}.tree)))
	if zl.score != score {
		return zl.score < score
	}
	return zl.name < name
}

func avlEntryEq(l, r *AVLNode) int {
	le := (*ZNode)(containerOf(unsafe.Pointer(l), unsafe.Offsetof(ZNode{}.tree)))
	re := (*ZNode)(containerOf(unsafe.Pointer(r), unsafe.Offsetof(ZNode{}.tree)))
	if le.score > re.score {
		return 1
	} else if le.score < re.score {
		return -1
	}
	if le.name < re.name {
		return 1
	} else if le.name > re.name {
		return -1
	}
	return 0
}
