package datastore

// AVLTree (Height-balanced BST)
type AVLTree struct {
	root *AVLNode
}

func NewAVLTree() *AVLTree {
	return &AVLTree{}
}

func (t *AVLTree) Insert(node *AVLNode, cmp func(node1 *AVLNode, node2 *AVLNode) int) {
	t.root = t.root.insert(node, cmp)
}

func (t *AVLTree) Remove(node *AVLNode, cmp func(node1 *AVLNode, node2 *AVLNode) int) {
	t.root = t.root.remove(node, cmp)
}

func (t *AVLTree) Search(node *AVLNode, cmp func(node1 *AVLNode, node2 *AVLNode) int) *AVLNode {
	return t.root.search(node, cmp)
}

func (t *AVLTree) Traverse() []*AVLNode {
	res := make([]*AVLNode, 0)
	traverse(t.root, &res)
	return res
}

func traverse(node *AVLNode, result *[]*AVLNode) {
	if node != nil {
		traverse(node.left, result)
		*result = append(*result, node)
		traverse(node.right, result)
	}
}

func (t *AVLTree) Offset(node *AVLNode, offset int32) *AVLNode {
	return node.offset(offset)
}

func (t *AVLTree) Dispose() {
	t.root.dispose()
}

type AVLNode struct {
	height int32
	count  int32
	left   *AVLNode
	right  *AVLNode
}

func (currNode *AVLNode) search(node *AVLNode, cmp func(node1 *AVLNode, node2 *AVLNode) int) *AVLNode {
	if currNode == nil {
		return nil
	}
	if cmp(node, currNode) >= 1 {
		return currNode.right.search(node, cmp)
	} else if cmp(node, currNode) <= -1 {
		return currNode.left.search(node, cmp)
	} else { //cpm == 1
		return node
	}
}

func (currNode *AVLNode) insert(node *AVLNode, cmp func(node1 *AVLNode, node2 *AVLNode) int) *AVLNode {
	if currNode == nil {
		return node
	} else if cmp(node, currNode) > 0 {
		currNode.right = currNode.right.insert(node, cmp)
	} else if cmp(node, currNode) < 0 {
		currNode.left = currNode.left.insert(node, cmp)
	} else {
		currNode = node
	}
	currNode.update()
	return currNode.rebalance()
}

func (currNode *AVLNode) remove(node *AVLNode, cmp func(node1 *AVLNode, node2 *AVLNode) int) *AVLNode {
	if currNode == nil {
		return nil
	}
	if cmp(node, currNode) > 0 {
		currNode.right = currNode.right.remove(node, cmp)
	} else if cmp(node, currNode) < 0 {
		currNode.left = currNode.left.remove(node, cmp)
	} else if currNode.left == nil && currNode.right == nil {
		currNode = nil
	} else if currNode.left == nil {
		currNode = currNode.right
	} else if currNode.right == nil {
		currNode = currNode.left
	} else { //has two children
		inOrderSuccessor := currNode.right.findSmallest()
		inOrderSuccessor.right = currNode.right.remove(inOrderSuccessor, cmp)
		inOrderSuccessor.left = currNode.left
		currNode = inOrderSuccessor
	}
	currNode.update()
	return currNode.rebalance()
}

func (currNode *AVLNode) findSmallest() *AVLNode {
	if currNode.left != nil {
		return currNode.left.findSmallest()
	} else {
		return currNode
	}
}

func (currNode *AVLNode) dispose() {
	if currNode == nil {
		return
	}
	currNode.left.dispose()
	currNode.right.dispose()
	*currNode = AVLNode{}
}

func (currNode *AVLNode) offset(offset int32) *AVLNode {
	if currNode == nil {
		return nil
	}
	node := currNode
	pos := node.left.getCount() + 1
	for offset != pos {
		if offset < pos {
			// The target is inside the left subtree
			node = node.left
			pos -= node.right.getCount() + 1
		} else {
			// The target is inside the right subtree
			node = node.right
			pos += node.left.getCount() + 1
		}
		if node == nil {
			return nil
		}
	}

	return node
}

func (node *AVLNode) rebalance() *AVLNode {
	if node == nil {
		return node
	}
	if node.balanceFactor() < -1 {
		//right-heavy
		node = node.fixRight()
	} else if node.balanceFactor() > 1 {
		//left-heavy
		node = node.fixLeft()
	}
	return node
}

// fixLeft re-balance a right-heavy node (double or single right rotate)
func (node *AVLNode) fixLeft() *AVLNode {
	if node.left.balanceFactor() < 0 {
		node.left = node.left.rotateLeft()
	}
	return node.rotateRight()
}

// rotateRight balances a left-heavy node (double or single left rotate)
func (node *AVLNode) fixRight() *AVLNode {
	if node.right.balanceFactor() > 0 {
		node.right = node.right.rotateRight()
	}
	return node.rotateLeft()
}

func (node *AVLNode) rotateRight() *AVLNode {
	leftChild := node.left
	node.left = leftChild.right
	leftChild.right = node
	node.update()
	leftChild.update()
	return leftChild
}

func (node *AVLNode) rotateLeft() *AVLNode {
	rightChild := node.right
	node.right = rightChild.left
	rightChild.left = node
	node.update()
	rightChild.update()
	return rightChild
}

func (node *AVLNode) balanceFactor() int32 {
	return node.left.getHeight() - node.right.getHeight()
}

func (node *AVLNode) update() {
	if node == nil {
		return
	}
	node.height = max(node.left.getHeight(), node.right.getHeight()) + 1
	node.count = node.left.getCount() + node.right.getCount() + 1
}

func (node *AVLNode) getHeight() int32 {
	if node == nil {
		return -1
	}
	return node.height
}

func (node *AVLNode) getCount() int32 {
	if node == nil {
		return 0
	}
	return node.count
}
