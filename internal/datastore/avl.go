package datastore

import "github.com/miladbarzideh/goldis/utils"

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
	parent *AVLNode
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
		currNode.right.parent = currNode
	} else if cmp(node, currNode) < 0 {
		currNode.left = currNode.left.insert(node, cmp)
		currNode.left.parent = currNode
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
		currNode.right.parent = currNode.parent
		currNode = currNode.right
	} else if currNode.right == nil {
		currNode.left.parent = currNode.parent
		currNode = currNode.left
	} else { //has two children
		inOrderSuccessor := currNode.right.findSmallest()
		inOrderSuccessor.right = currNode.right.remove(inOrderSuccessor, cmp)
		inOrderSuccessor.left = currNode.left
		inOrderSuccessor.parent = currNode.parent
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
	pos := int32(0)
	for offset != pos {
		if pos < offset && pos+currNode.right.getCount() >= offset {
			// the target is inside the right subtree
			currNode = currNode.right
			pos += currNode.left.getCount() + 1
		} else if pos > offset && pos-currNode.left.getCount() <= offset {
			// the target is inside the left subtree
			currNode = currNode.left
			pos -= currNode.right.getCount() + 1
		} else {
			// go to the parent
			if currNode.parent == nil {
				return nil
			}
			if currNode.parent.right == currNode {
				pos -= currNode.left.getCount() + 1
			} else {
				pos += currNode.right.getCount() + 1
			}
			currNode = currNode.parent
		}
	}
	return currNode
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
	leftChild.parent = node.parent
	node.parent = leftChild
	node.update()
	leftChild.update()
	return leftChild
}

func (node *AVLNode) rotateLeft() *AVLNode {
	rightChild := node.right
	node.right = rightChild.left
	rightChild.left = node
	rightChild.parent = node.parent
	node.parent = rightChild
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
	node.height = utils.Max(node.left.getHeight(), node.right.getHeight()) + 1
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
