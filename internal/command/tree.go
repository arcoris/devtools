// Copyright 2026 The ARCORIS Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidTree reports that a command tree definition is malformed.
	ErrInvalidTree = errors.New("command tree is invalid")
)

// Tree is a validated framework-neutral command tree.
//
// Tree owns one root Node. The root Node must have kind NodeRoot. Node-level
// invariants are delegated to Node.Validate. Tree additionally verifies global
// invariants that require looking at the whole tree, such as unique IDs and
// unique Paths.
//
// Tree is an immutable-style value object. Methods that expose nodes or slices
// return detached values. Node also preserves detached child slices.
type Tree struct {
	root Node
}

// NewTree validates root and returns a Tree.
//
// root must be a NodeRoot node. Use NewTreeFromChildren when constructing a
// tree from top-level command nodes.
func NewTree(root Node) (Tree, error) {
	tree := Tree{root: root}
	if err := tree.Validate(); err != nil {
		return Tree{}, err
	}

	return tree, nil
}

// MustTree validates root and returns a Tree.
//
// MustTree panics on invalid input. It is intended for static command
// definitions and tests where invalid trees are programmer errors.
func MustTree(root Node) Tree {
	tree, err := NewTree(root)
	if err != nil {
		panic(err)
	}

	return tree
}

// NewTreeFromChildren constructs a root node from children and returns a Tree.
func NewTreeFromChildren(children ...Node) (Tree, error) {
	root, err := NewRootNode(children...)
	if err != nil {
		return Tree{}, fmt.Errorf("%w: invalid root node: %w", ErrInvalidTree, err)
	}

	return NewTree(root)
}

// MustTreeFromChildren constructs a Tree from top-level children and panics on
// invalid input.
func MustTreeFromChildren(children ...Node) Tree {
	tree, err := NewTreeFromChildren(children...)
	if err != nil {
		panic(err)
	}

	return tree
}

// Root returns the tree root node.
func (tree Tree) Root() Node {
	return tree.root
}

// IsZero reports whether the tree has not been initialized.
//
// A zero Tree is invalid because it does not contain a valid root node.
func (tree Tree) IsZero() bool {
	return tree.root.Kind() == ""
}

// IsEmpty reports whether the tree has no top-level command nodes.
//
// A valid empty tree still has a root node.
func (tree Tree) IsEmpty() bool {
	if tree.IsZero() {
		return true
	}

	return tree.root.ChildCount() == 0
}

// TopLevel returns the root's direct children.
func (tree Tree) TopLevel() []Node {
	if tree.IsZero() {
		return nil
	}

	return tree.root.Children()
}
