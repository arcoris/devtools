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

import "fmt"

// WalkOrder defines how a Tree traversal visits nodes.
type WalkOrder string

const (
	// WalkPreOrder visits a node before its children.
	WalkPreOrder WalkOrder = "pre-order"

	// WalkPostOrder visits a node after its children.
	WalkPostOrder WalkOrder = "post-order"
)

// TreeWalkFunc is called for every node visited during tree traversal.
//
// The Path argument is the node's own path and is provided for convenience. It
// is equivalent to node.Path().
type TreeWalkFunc func(path Path, node Node) error

// String returns the stable text form of order.
func (order WalkOrder) String() string {
	return string(order)
}

// IsValid reports whether order is a known walk order.
func (order WalkOrder) IsValid() bool {
	return validateWalkOrder(order) == nil
}

// Walk visits all nodes in pre-order.
//
// Walk is equivalent to WalkOrder(WalkPreOrder, fn).
func (tree Tree) Walk(fn TreeWalkFunc) error {
	return tree.WalkOrder(WalkPreOrder, fn)
}

// WalkOrder visits all nodes in the requested order.
func (tree Tree) WalkOrder(order WalkOrder, fn TreeWalkFunc) error {
	if fn == nil {
		return nil
	}

	if err := tree.ensureValidForTraversal(); err != nil {
		return err
	}

	if err := validateWalkOrder(order); err != nil {
		return err
	}

	switch order {
	case WalkPreOrder:
		return walkNodePreOrder(tree.root, fn)
	case WalkPostOrder:
		return walkNodePostOrder(tree.root, fn)
	default:
		return fmt.Errorf("%w: unsupported walk order %q", ErrInvalidTree, order)
	}
}

// validateWalkOrder validates traversal order.
func validateWalkOrder(order WalkOrder) error {
	switch order {
	case WalkPreOrder, WalkPostOrder:
		return nil
	default:
		return fmt.Errorf("%w: unknown walk order %q", ErrInvalidTree, order)
	}
}

// ensureValidForTraversal verifies that traversal can safely operate on tree.
func (tree Tree) ensureValidForTraversal() error {
	if tree.IsZero() {
		return fmt.Errorf("%w: root node is not set", ErrInvalidTree)
	}

	if !tree.root.IsRoot() {
		return fmt.Errorf(
			"%w: root node must have kind %q, got %q",
			ErrInvalidTree,
			NodeRoot,
			tree.root.Kind(),
		)
	}

	return nil
}

// walkNodePreOrder walks node before its descendants.
func walkNodePreOrder(node Node, fn TreeWalkFunc) error {
	if err := fn(node.Path(), node); err != nil {
		return err
	}

	for _, child := range node.children {
		if err := walkNodePreOrder(child, fn); err != nil {
			return err
		}
	}

	return nil
}

// walkNodePostOrder walks node after its descendants.
func walkNodePostOrder(node Node, fn TreeWalkFunc) error {
	for _, child := range node.children {
		if err := walkNodePostOrder(child, fn); err != nil {
			return err
		}
	}

	if err := fn(node.Path(), node); err != nil {
		return err
	}

	return nil
}
