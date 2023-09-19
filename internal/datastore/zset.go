package datastore

import (
	"fmt"
	"log"
	"strings"
	"unsafe"

	"github.com/miladbarzideh/goldis/utils"
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
		hmap: NewHMap(ZNodeComparator),
		tree: NewAVLTree(AVLTreeComparator),
	}
}

func NewZNode(name string, score float64) *ZNode {
	return &ZNode{
		hmap:  HNode{hcode: utils.Hash(name)},
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
		zset.tree.Insert(&node.tree)
		return true
	}
	zset.update(node, score)
	return false
}

func (zset *ZSet) update(node *ZNode, score float64) {
	if node.score == score {
		return
	}
	zset.tree.Remove(&node.tree)
	key := newHKey(node.name)
	zset.hmap.Pop(&key.node)
	newNode := NewZNode(node.name, score)
	zset.hmap.Insert(&newNode.hmap)
	zset.tree.Insert(&newNode.tree)
}

func (zset *ZSet) Lookup(name string) *ZNode {
	if zset.tree == nil {
		return nil
	}
	key := newHKey(name)
	found := zset.hmap.Lookup(&key.node)
	if found == nil {
		return nil
	}
	return (*ZNode)(utils.ContainerOf(unsafe.Pointer(found), unsafe.Offsetof(ZNode{}.hmap)))
}

func (zset *ZSet) Pop(name string) *ZNode {
	if zset.tree == nil {
		return nil
	}
	key := newHKey(name)
	found := zset.hmap.Pop(&key.node)
	if found == nil {
		return nil
	}

	node := (*ZNode)(utils.ContainerOf(unsafe.Pointer(found), unsafe.Offsetof(ZNode{}.hmap)))
	zset.tree.Remove(&node.tree)
	return node
}

func (zset *ZSet) Show() string {
	printHashtable(zset.hmap.Keys())
	return printTreeNode(zset.tree.Traverse())
}

func printHashtable(nodes []*HNode) {
	log.Print("Hashtable name-score pair:\n")
	for _, node := range nodes {
		entry := (*ZNode)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(ZNode{}.hmap)))
		log.Printf("%v => %v", entry.name, entry.score)
	}
}

func printTreeNode(nodes []*AVLNode) string {
	log.Print("AVL Tree Inorder Traversal:\n")
	res := strings.Builder{}
	for i, node := range nodes {
		entry := (*ZNode)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(ZNode{}.tree)))
		sn := fmt.Sprintf("%v) %v => %v\n", i+1, entry.score, entry.name)
		log.Print(sn)
		res.WriteString(sn)
	}
	return res.String()
}

func (zset *ZSet) Query(score float64, name string, offset int32, limit uint32) []*ZNode {
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
		node := zset.tree.Offset(found, offset)
		znode := (*ZNode)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(ZNode{}.tree)))
		res := make([]*ZNode, 0)
		n := uint32(0)
		for znode != nil && n < limit {
			res = append(res, znode)
			node = znode.tree.offset(1)
			if node == nil {
				break
			}
			znode = (*ZNode)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(ZNode{}.tree)))
			n++
		}
		return res
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
		node: HNode{hcode: utils.Hash(name)},
		name: name,
	}
}

// zless compare by the (score, name) tuple
func zless(lhs *AVLNode, score float64, name string) bool {
	zl := (*ZNode)(utils.ContainerOf(unsafe.Pointer(lhs), unsafe.Offsetof(ZNode{}.tree)))
	if zl.score != score {
		return zl.score < score
	}
	return zl.name < name
}
