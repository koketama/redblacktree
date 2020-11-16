package redblacktree

import (
	"fmt"
	"sync"

	"github.com/emirpasic/gods/utils"
	"github.com/pkg/errors"
)

type color bool

const (
	black, red color = true, false
)

type Value interface {
	ID() string
	String() string
}

// Tree holds elements of the red-black tree
type Tree struct {
	sync.RWMutex
	root       *Node
	size       int
	comparator utils.Comparator
}

// Node is a single element within the tree
type Node struct {
	key    interface{}
	values []Value
	color  color
	left   *Node
	right  *Node
	parent *Node
}

// NewWith instantiates a red-black tree with the custom comparator.
func NewWith(comparator utils.Comparator) (*Tree, error) {
	if comparator == nil {
		return nil, errors.New("comparator required")
	}

	return &Tree{
		comparator: comparator,
	}, nil
}

// Put inserts node into the tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Put(key interface{}, value Value) {
	if key == nil || value == nil {
		return
	}

	tree.Lock()
	defer tree.Unlock()

	var insertedNode *Node
	if tree.root == nil {
		// Assert key is of comparator's type for initial tree
		tree.comparator(key, key)
		tree.size++
		tree.root = &Node{key: key, values: []Value{value}, color: red}
		insertedNode = tree.root
	} else {
		node := tree.root
		loop := true
		for loop {
			compare := tree.comparator(key, node.key)
			switch {
			case compare == 0:
				for _, v := range node.values {
					if v.ID() == value.ID() {
						return // duplicated
					}
				}

				tree.size++
				node.values = append(node.values, value)
				return

			case compare < 0:
				if node.left == nil {
					tree.size++
					node.left = &Node{key: key, values: []Value{value}, color: red}
					insertedNode = node.left
					loop = false
				} else {
					node = node.left
				}

			case compare > 0:
				if node.right == nil {
					tree.size++
					node.right = &Node{key: key, values: []Value{value}, color: red}
					insertedNode = node.right
					loop = false
				} else {
					node = node.right
				}
			}
		}
		insertedNode.parent = node
	}
	tree.insertCase1(insertedNode)
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Get(key interface{}) (values []Value, found bool) {
	if key == nil {
		return nil, false
	}

	tree.RLock()
	defer tree.RUnlock()

	node := tree.lookup(key)
	if node != nil {
		return node.values, true
	}
	return nil, false
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Remove(key interface{}) {
	if key == nil {
		return
	}

	tree.Lock()
	defer tree.Unlock()

	node := tree.lookup(key)
	if node == nil {
		return
	}
	tree.size -= len(node.values)

	if node.left != nil && node.right != nil {
		pred := node.left.maximumNode()
		node.key = pred.key
		node.values = pred.values
		node = pred
	}

	var child *Node
	if node.left == nil || node.right == nil {
		if node.right == nil {
			child = node.left
		} else {
			child = node.right
		}
		if node.color == black {
			node.color = nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.parent == nil && child != nil {
			child.color = black
		}
	}
}

// Empty returns true if tree does not contain any nodes
func (tree *Tree) Empty() bool {
	tree.RLock()
	defer tree.RUnlock()

	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *Tree) Size() int {
	tree.RLock()
	defer tree.RUnlock()

	return tree.size
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *Tree) Left() (values []Value) {
	tree.RLock()
	defer tree.RUnlock()

	var parent *Node
	for current := tree.root; current != nil; {
		parent = current
		current = current.left
	}

	if parent != nil {
		values = parent.values
	}
	return
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *Tree) Right() (values []Value) {
	tree.RLock()
	defer tree.RUnlock()

	var parent *Node
	for current := tree.root; current != nil; {
		parent = current
		current = current.right
	}

	if parent != nil {
		values = parent.values
	}
	return
}

// Floor Finds floor node of the input key, return the floor node or nil if no floor is found.
// Second return parameter is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Floor(key interface{}) (floor []Value, found bool) {
	tree.RLock()
	defer tree.RUnlock()

	found = false
	node := tree.root
	for node != nil {
		compare := tree.comparator(key, node.key)
		switch {
		case compare == 0:
			return node.values, true
		case compare < 0:
			node = node.left
		case compare > 0:
			floor, found = node.values, true
			node = node.right
		}
	}
	if found {
		return floor, true
	}
	return nil, false
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling is found.
// Second return parameter is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Ceiling(key interface{}) (ceiling []Value, found bool) {
	tree.RLock()
	defer tree.RUnlock()

	found = false
	node := tree.root
	for node != nil {
		compare := tree.comparator(key, node.key)
		switch {
		case compare == 0:
			return node.values, true
		case compare < 0:
			ceiling, found = node.values, true
			node = node.left
		case compare > 0:
			node = node.right
		}
	}
	if found {
		return ceiling, true
	}
	return nil, false
}

// Clear removes all nodes from the tree.
func (tree *Tree) Clear() {
	tree.Lock()
	defer tree.Unlock()

	tree.root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *Tree) String() string {
	tree.RLock()
	defer tree.RUnlock()

	str := "RedBlackTree\n"
	if !tree.Empty() {
		output(tree.root, "", true, &str)
	}
	return str
}

func (node *Node) String() string {
	return fmt.Sprintf("%+v", node.key)
}

func output(node *Node, prefix string, isTail bool, str *string) {
	if node.right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.left, newPrefix, true, str)
	}
}

func (tree *Tree) lookup(key interface{}) *Node {
	node := tree.root
	for node != nil {
		compare := tree.comparator(key, node.key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.left
		case compare > 0:
			node = node.right
		}
	}
	return nil
}

func (node *Node) grandparent() *Node {
	if node != nil && node.parent != nil {
		return node.parent.parent
	}
	return nil
}

func (node *Node) uncle() *Node {
	if node == nil || node.parent == nil || node.parent.parent == nil {
		return nil
	}
	return node.parent.sibling()
}

func (node *Node) sibling() *Node {
	if node == nil || node.parent == nil {
		return nil
	}
	if node == node.parent.left {
		return node.parent.right
	}
	return node.parent.left
}

func (tree *Tree) rotateLeft(node *Node) {
	right := node.right
	tree.replaceNode(node, right)
	node.right = right.left
	if right.left != nil {
		right.left.parent = node
	}
	right.left = node
	node.parent = right
}

func (tree *Tree) rotateRight(node *Node) {
	left := node.left
	tree.replaceNode(node, left)
	node.left = left.right
	if left.right != nil {
		left.right.parent = node
	}
	left.right = node
	node.parent = left
}

func (tree *Tree) replaceNode(old *Node, new *Node) {
	if old.parent == nil {
		tree.root = new
	} else {
		if old == old.parent.left {
			old.parent.left = new
		} else {
			old.parent.right = new
		}
	}
	if new != nil {
		new.parent = old.parent
	}
}

func (tree *Tree) insertCase1(node *Node) {
	if node.parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *Tree) insertCase2(node *Node) {
	if nodeColor(node.parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *Tree) insertCase3(node *Node) {
	uncle := node.uncle()
	if nodeColor(uncle) == red {
		node.parent.color = black
		uncle.color = black
		node.grandparent().color = red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *Tree) insertCase4(node *Node) {
	grandparent := node.grandparent()
	if node == node.parent.right && node.parent == grandparent.left {
		tree.rotateLeft(node.parent)
		node = node.left
	} else if node == node.parent.left && node.parent == grandparent.right {
		tree.rotateRight(node.parent)
		node = node.right
	}
	tree.insertCase5(node)
}

func (tree *Tree) insertCase5(node *Node) {
	node.parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.parent.left && node.parent == grandparent.left {
		tree.rotateRight(grandparent)
	} else if node == node.parent.right && node.parent == grandparent.right {
		tree.rotateLeft(grandparent)
	}
}

func (node *Node) maximumNode() *Node {
	if node == nil {
		return nil
	}
	for node.right != nil {
		node = node.right
	}
	return node
}

func (tree *Tree) deleteCase1(node *Node) {
	if node.parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *Tree) deleteCase2(node *Node) {
	sibling := node.sibling()
	if nodeColor(sibling) == red {
		node.parent.color = red
		sibling.color = black
		if node == node.parent.left {
			tree.rotateLeft(node.parent)
		} else {
			tree.rotateRight(node.parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *Tree) deleteCase3(node *Node) {
	sibling := node.sibling()
	if nodeColor(node.parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.left) == black &&
		nodeColor(sibling.right) == black {
		sibling.color = red
		tree.deleteCase1(node.parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *Tree) deleteCase4(node *Node) {
	sibling := node.sibling()
	if nodeColor(node.parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.left) == black &&
		nodeColor(sibling.right) == black {
		sibling.color = red
		node.parent.color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *Tree) deleteCase5(node *Node) {
	sibling := node.sibling()
	if node == node.parent.left &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.left) == red &&
		nodeColor(sibling.right) == black {
		sibling.color = red
		sibling.left.color = black
		tree.rotateRight(sibling)
	} else if node == node.parent.right &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.right) == red &&
		nodeColor(sibling.left) == black {
		sibling.color = red
		sibling.right.color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *Tree) deleteCase6(node *Node) {
	sibling := node.sibling()
	sibling.color = nodeColor(node.parent)
	node.parent.color = black
	if node == node.parent.left && nodeColor(sibling.right) == red {
		sibling.right.color = black
		tree.rotateLeft(node.parent)
	} else if nodeColor(sibling.left) == red {
		sibling.left.color = black
		tree.rotateRight(node.parent)
	}
}

func nodeColor(node *Node) color {
	if node == nil {
		return black
	}
	return node.color
}
