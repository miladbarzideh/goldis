package datastore

type LNode struct {
	next     *LNode
	previous *LNode
}

type DList struct {
	head *LNode
}

func NewDList() *DList {
	return &DList{}
}

func (dl *DList) IsEmpty() bool {
	return dl.head == nil
}

func (dl *DList) Detach(node *LNode, cmp func(node1, node2 *LNode) bool) {
	if dl.IsEmpty() {
		return
	}
	if cmp(node, dl.head) {
		dl.head = node.next
	} else {
		node.previous.next = node.next
		if node.next != nil {
			node.next.previous = node.previous
		}
	}
	node.next = nil
	node.previous = nil
}

func (dl *DList) GetHead() *LNode {
	return dl.head
}

func (dl *DList) InsertBefore(newNode *LNode) {
	if dl.IsEmpty() {
		dl.head = newNode
	} else {
		newNode.next = dl.head
		dl.head.previous = newNode
		dl.head = newNode
	}
}

func (dl *DList) Iterator() func() *LNode {
	curr := dl.head
	return func() *LNode {
		if curr == nil {
			return nil
		}
		node := curr
		curr = curr.next
		return node
	}
}
