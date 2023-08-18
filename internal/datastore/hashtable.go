package datastore

const (
	resizingWork  = 128
	maxLoadFactor = 8
	bucketSize    = 4
)

type HNode struct {
	next  *HNode
	hcode uint64
}

type HTab struct {
	tab  []*HNode
	mask uint64
	size int
}

type HMap struct {
	tab1        HTab
	tab2        HTab
	resizingPos uint64
}

func initHTab(htab *HTab, n uint64) {
	htab.tab = make([]*HNode, n)
	htab.mask = n - 1
	htab.size = 0
}

func (htab *HTab) insert(node *HNode) {
	pos := node.hcode & htab.mask
	node.next = htab.tab[pos]
	htab.tab[pos] = node
	htab.size++
}

func (htab *HTab) lookup(key *HNode, cmp func(node1 *HNode, node2 *HNode) bool) **HNode {
	if htab.tab == nil || htab.size == 0 {
		return nil
	}
	pos := key.hcode & htab.mask
	head := &htab.tab[pos]
	for *head != nil {
		if cmp(key, *head) {
			return head
		}
		head = &((*head).next)
	}
	return nil
}

func (htab *HTab) detach(from **HNode) *HNode {
	node := *from
	*from = node.next
	node.next = nil
	htab.size--
	return node
}

func (htab *HTab) freeHTab() {
	for i := range htab.tab {
		htab.tab[i] = nil
	}
	htab.tab = nil
}

func (htab *HTab) keys() []*HNode {
	if htab.size == 0 {
		return make([]*HNode, 0)
	}
	nodes := make([]*HNode, 0)
	for _, node := range htab.tab {
		for node != nil {
			nodes = append(nodes, node)
			node = node.next
		}
	}
	return nodes
}

func NewHMap() *HMap {
	hmap := HMap{}
	initHTab(&hmap.tab1, bucketSize)
	return &hmap
}

func (hmap *HMap) Lookup(key *HNode, cmp func(node1 *HNode, node2 *HNode) bool) *HNode {
	hmap.helpResizing()
	node := hmap.tab1.lookup(key, cmp)
	if node == nil {
		node = hmap.tab2.lookup(key, cmp)
	}
	if node != nil {
		return *node
	}
	return nil
}

func (hmap *HMap) Insert(node *HNode) {
	hmap.tab1.insert(node)
	if hmap.tab2.tab == nil {
		loadFactor := uint64(hmap.tab1.size)/hmap.tab1.mask + 1
		if loadFactor > maxLoadFactor {
			hmap.startResizing()
		}
	}
	hmap.helpResizing()
}

func (hmap *HMap) Pop(key *HNode, cmp func(node1 *HNode, node2 *HNode) bool) *HNode {
	hmap.helpResizing()
	node := hmap.tab1.lookup(key, cmp)
	if node != nil {
		return hmap.tab1.detach(node)
	}
	node = hmap.tab2.lookup(key, cmp)
	if node != nil {
		return hmap.tab2.detach(node)
	}
	return nil
}

func (hmap *HMap) Keys() []*HNode {
	t1 := hmap.tab1.keys()
	t2 := hmap.tab2.keys()
	return append(t1, t2...)
}

func (hmap *HMap) Destroy() {
	hmap.tab1.freeHTab()
	hmap.tab2.freeHTab()
	hmap.tab1 = HTab{}
	hmap.tab2 = HTab{}
}

func (hmap *HMap) helpResizing() {
	if hmap.tab2.tab == nil {
		return
	}

	nwork := 0
	for nwork < resizingWork && hmap.tab2.size > 0 {
		from := &hmap.tab2.tab[hmap.resizingPos]
		if *from == nil {
			hmap.resizingPos++
			continue
		}
		hmap.tab1.insert(hmap.tab2.detach(from))
		nwork++
	}

	if hmap.tab2.size == 0 {
		hmap.tab2.freeHTab()
		hmap.tab2 = HTab{}
	}
}

func (hmap *HMap) startResizing() {
	// create a bigger hashtable and swap them
	hmap.tab2 = hmap.tab1
	newSize := (hmap.tab1.mask + 1) * 2
	initHTab(&hmap.tab1, newSize)
	hmap.resizingPos = 0
}
