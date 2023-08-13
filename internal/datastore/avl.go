package datastore

import "log"

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

func (t *AVLTree) DisplayNodes() {
	t.root.displayNodes()
}

type AVLNode struct {
	height uint32
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
	} else if cmp(node, currNode) >= 1 {
		currNode.right = currNode.right.insert(node, cmp)
	} else if cmp(node, currNode) <= -1 {
		currNode.left = currNode.left.insert(node, cmp)
	} else {
		currNode = node
	}
	return currNode.rebalance()
}

func (currNode *AVLNode) remove(node *AVLNode, cmp func(node1 *AVLNode, node2 *AVLNode) int) *AVLNode {
	if currNode == nil {
		return nil
	}
	if cmp(node, currNode) >= 1 {
		currNode.right = currNode.right.remove(node, cmp)
	} else if cmp(node, currNode) <= -1 {
		currNode.left = currNode.left.remove(node, cmp)
	} else if currNode.left == nil && currNode.right == nil {
		currNode = nil
	} else if currNode.left == nil {
		currNode = currNode.right
	} else if currNode.right == nil {
		currNode = currNode.left
	} else { //has two children
		inOrderSuccessor := currNode.right.findSmallest()
		currNode = inOrderSuccessor
		currNode.right = currNode.right.remove(inOrderSuccessor, cmp)
	}
	return currNode.rebalance()
}

func (currNode *AVLNode) findSmallest() *AVLNode {
	if currNode.left != nil {
		return currNode.left.findSmallest()
	} else {
		return currNode
	}
}

func (currNode *AVLNode) displayNodes() {
	if currNode == nil {
		return
	}
	if currNode.left != nil {
		currNode.left.displayNodes()
	}
	log.Printf("%v", avlEntryContainerOf(currNode).value)
	if currNode.right != nil {
		currNode.right.displayNodes()
	}
}

func (node *AVLNode) rebalance() *AVLNode {
	if node == nil {
		return node
	}
	if node.balanceFactor() < -1 {
		//left-heavy
		node = node.fixRight()
	} else if node.balanceFactor() > 1 {
		//right-heavy
		node = node.fixLeft()
	}
	return node
}

// fixLeft re-balance a right-heavy node (double or single right rotate)
func (node *AVLNode) fixLeft() *AVLNode {
	if node.left.balanceFactor() > 0 {
		node.left = node.rotateLeft()
	}
	return node.rotateRight()
}

// rotateRight balances a left-heavy node (double or single left rotate)
func (node *AVLNode) fixRight() *AVLNode {
	if node.right.balanceFactor() < 0 {
		node.right = node.rotateRight()
	}
	return node.rotateLeft()
}

func (node *AVLNode) rotateRight() *AVLNode {
	leftChild := node.left
	node.left = leftChild.right
	leftChild.right = node
	node.updateHeight()
	leftChild.updateHeight()
	return leftChild
}

func (node *AVLNode) rotateLeft() *AVLNode {
	rightChild := node.right
	node.right = rightChild.left
	rightChild.left = node
	node.updateHeight()
	rightChild.updateHeight()
	return rightChild
}

func (node *AVLNode) balanceFactor() int32 {
	return int32(node.left.getHeight() - node.right.getHeight())
}

func (node *AVLNode) updateHeight() {
	node.height = max(node.left.getHeight(), node.right.getHeight()) + 1
}

func (node *AVLNode) getHeight() uint32 {
	if node == nil {
		return 0
	}
	return node.height
}

func max(a uint32, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}
