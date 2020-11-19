// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package redblacktree implements a red-black tree.
//
// Used by TreeSet and TreeMap.
//
// Structure is not thread safe.
//
// References: http://en.wikipedia.org/wiki/Red%E2%80%93black_tree
package pkg

import (
	"fmt"

	"github.com/emirpasic/gods/utils"
)

type color bool

const (
	black, red color = true, false
)

// Tree holds elements of the red-black tree
type Tree struct {
	Root       *Node
	size       int
	Comparator utils.Comparator
}

type Value interface {
	ID() string
}

// Node is a single element within the tree
type Node struct {
	Key    interface{}
	Values []Value
	color  color
	Left   *Node
	Right  *Node
	Parent *Node
}

// NewWith instantiates a red-black tree with the custom comparator.
func NewWith(comparator utils.Comparator) *Tree {
	return &Tree{Comparator: comparator}
}

// Put inserts node into the tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Put(key interface{}, value Value) {
	var insertedNode *Node
	if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		tree.Comparator(key, key)
		tree.Root = &Node{Key: key, Values: []Value{value}, color: red}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		for loop {
			compare := tree.Comparator(key, node.Key)
			switch {
			case compare == 0:
				for _, v := range node.Values {
					if v.ID() == value.ID() { // duplicated
						return
					}
				}

				node.Key = key
				node.Values = append(node.Values, value)
				tree.size++
				return
			case compare < 0:
				if node.Left == nil {
					node.Left = &Node{Key: key, Values: []Value{value}, color: red}
					insertedNode = node.Left
					loop = false
				} else {
					node = node.Left
				}
			case compare > 0:
				if node.Right == nil {
					node.Right = &Node{Key: key, Values: []Value{value}, color: red}
					insertedNode = node.Right
					loop = false
				} else {
					node = node.Right
				}
			}
		}
		insertedNode.Parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Get(key interface{}) (values []Value, found bool) {
	node := tree.lookup(key)
	if node != nil {
		return node.Values, true
	}
	return nil, false
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Remove(key interface{}) {
	node := tree.lookup(key)
	tree.remove(node)
}

func (tree *Tree) remove(node *Node) {
	if node == nil {
		return
	}
	tree.size -= len(node.Values)

	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		node.Key = pred.Key
		node.Values = pred.Values
		node = pred
	}

	var child *Node
	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.color == black {
			node.color = nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = black
		}
	}
}

// Empty returns true if tree does not contain any nodes
func (tree *Tree) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *Tree) Size() int {
	return tree.size
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *Tree) left() *Node {
	if tree.Empty() {
		return nil
	}

	var parent *Node
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *Tree) Left() (key interface{}, values []Value) {
	node := tree.left()
	if node == nil {
		return
	}

	return node.Key, node.Values
}

func (tree *Tree) PopLeft() (key interface{}, values []Value) {
	node := tree.left()
	if node == nil {
		return
	}

	key, values = node.Key, node.Values

	tree.remove(node)
	return
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *Tree) right() *Node {
	if tree.Empty() {
		return nil
	}

	var parent *Node
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *Tree) Right() (key interface{}, values []Value) {
	node := tree.right()
	if node == nil {
		return
	}

	return node.Key, node.Values
}

func (tree *Tree) PopRight() (key interface{}, values []Value) {
	node := tree.right()
	if node == nil {
		return
	}

	key, values = node.Key, node.Values

	tree.remove(node)
	return
}

// Clear removes all nodes from the tree.
func (tree *Tree) Clear() {
	tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *Tree) String() string {
	str := "RedBlackTree\n"
	if !tree.Empty() {
		output(tree.Root, "", true, &str)
	}
	return str
}

func (node *Node) String() string {
	return fmt.Sprintf("%v", node.Key)
}

func output(node *Node, prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.Left, newPrefix, true, str)
	}
}

func (tree *Tree) lookup(key interface{}) *Node {
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

func (node *Node) grandparent() *Node {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *Node) uncle() *Node {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *Node) sibling() *Node {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (tree *Tree) rotateLeft(node *Node) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *Tree) rotateRight(node *Node) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *Tree) replaceNode(old *Node, new *Node) {
	if old.Parent == nil {
		tree.Root = new
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = new
		} else {
			old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent = old.Parent
	}
}

func (tree *Tree) insertCase1(node *Node) {
	if node.Parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *Tree) insertCase2(node *Node) {
	if nodeColor(node.Parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *Tree) insertCase3(node *Node) {
	uncle := node.uncle()
	if nodeColor(uncle) == red {
		node.Parent.color = black
		uncle.color = black
		node.grandparent().color = red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *Tree) insertCase4(node *Node) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(node.Parent)
		node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *Tree) insertCase5(node *Node) {
	node.Parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *Node) maximumNode() *Node {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *Tree) deleteCase1(node *Node) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *Tree) deleteCase2(node *Node) {
	sibling := node.sibling()
	if nodeColor(sibling) == red {
		node.Parent.color = red
		sibling.color = black
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *Tree) deleteCase3(node *Node) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *Tree) deleteCase4(node *Node) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		node.Parent.color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *Tree) deleteCase5(node *Node) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == red &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		sibling.Left.color = black
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Right) == red &&
		nodeColor(sibling.Left) == black {
		sibling.color = red
		sibling.Right.color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *Tree) deleteCase6(node *Node) {
	sibling := node.sibling()
	sibling.color = nodeColor(node.Parent)
	node.Parent.color = black
	if node == node.Parent.Left && nodeColor(sibling.Right) == red {
		sibling.Right.color = black
		tree.rotateLeft(node.Parent)
	} else if nodeColor(sibling.Left) == red {
		sibling.Left.color = black
		tree.rotateRight(node.Parent)
	}
}

func nodeColor(node *Node) color {
	if node == nil {
		return black
	}
	return node.color
}
