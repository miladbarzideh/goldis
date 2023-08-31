package datastore

import (
	"unsafe"

	"github.com/miladbarzideh/goldis/utils"
)

type MapComparator func(a, b interface{}) bool

func MapEntryComparator(a, b interface{}) bool {
	l := a.(*HNode)
	r := b.(*HNode)
	le := (*MapEntry)(utils.ContainerOf(unsafe.Pointer(l), unsafe.Offsetof(MapEntry{}.node)))
	re := (*MapEntry)(utils.ContainerOf(unsafe.Pointer(r), unsafe.Offsetof(MapEntry{}.node)))
	return l.hcode == r.hcode && le.key == re.key
}

func ZNodeComparator(a, b interface{}) bool {
	key := a.(*HNode)
	node := b.(*HNode)
	if node.hcode != key.hcode {
		return false
	}
	znode := (*ZNode)(utils.ContainerOf(unsafe.Pointer(node), unsafe.Offsetof(ZNode{}.hmap)))
	hkey := (*HKey)(utils.ContainerOf(unsafe.Pointer(key), unsafe.Offsetof(HKey{}.node)))
	return znode.name == hkey.name
}

type TreeComparator func(a, b interface{}) int

func AVLTreeComparator(a, b interface{}) int {
	l := a.(*AVLNode)
	r := b.(*AVLNode)
	le := (*ZNode)(utils.ContainerOf(unsafe.Pointer(l), unsafe.Offsetof(ZNode{}.tree)))
	re := (*ZNode)(utils.ContainerOf(unsafe.Pointer(r), unsafe.Offsetof(ZNode{}.tree)))
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
