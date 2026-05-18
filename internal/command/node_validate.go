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

// Validate verifies the node definition and its structural invariants.
//
// Validate recursively validates child nodes, direct parent/child path
// relationships, and sibling Use/alias conflicts. Global registry validation can
// add cross-tree invariants later without weakening this local model.
func (node Node) Validate() error {
	if err := node.validateSelf(); err != nil {
		return err
	}

	if err := node.validateChildren(); err != nil {
		return err
	}

	return nil
}

// validateSelf validates node-local invariants that do not require sibling
// context.
func (node Node) validateSelf() error {
	if err := validateNodeKind(node.kind); err != nil {
		return err
	}

	switch node.kind {
	case NodeRoot:
		return node.validateRoot()
	case NodeFamily:
		return node.validateFamily()
	case NodeCommand:
		return node.validateCommand()
	default:
		return fmt.Errorf("%w: unsupported kind %q", ErrInvalidNode, node.kind)
	}
}

// validateRoot validates root-node invariants.
func (node Node) validateRoot() error {
	if !node.id.IsZero() {
		return fmt.Errorf("%w: root node must not have an ID", ErrInvalidNode)
	}

	if !node.path.IsRoot() {
		return fmt.Errorf("%w: root node path must be root", ErrInvalidNode)
	}

	if node.use != "" {
		return fmt.Errorf("%w: root node must not have a use segment", ErrInvalidNode)
	}

	if len(node.aliases) != 0 {
		return fmt.Errorf("%w: root node must not have aliases", ErrInvalidNode)
	}

	return nil
}

// validateFamily validates family-node invariants.
func (node Node) validateFamily() error {
	if err := node.validateNonRootIdentity(); err != nil {
		return err
	}

	if node.IsLeaf() {
		return fmt.Errorf("%w: family node %q must have children", ErrInvalidNode, node.path)
	}

	return nil
}

// validateCommand validates command-node invariants.
func (node Node) validateCommand() error {
	if err := node.validateNonRootIdentity(); err != nil {
		return err
	}

	if node.HasChildren() {
		return fmt.Errorf("%w: command node %q must not have children", ErrInvalidNode, node.path)
	}

	return nil
}

// validateNonRootIdentity validates common non-root invariants.
func (node Node) validateNonRootIdentity() error {
	if node.id.IsZero() {
		return fmt.Errorf("%w: non-root node must have an ID", ErrInvalidNode)
	}

	if err := node.id.Validate(); err != nil {
		return fmt.Errorf("%w: invalid node ID %q: %w", ErrInvalidNode, node.id, err)
	}

	if node.path.IsRoot() {
		return fmt.Errorf("%w: non-root node %q must have a non-root path", ErrInvalidNode, node.id)
	}

	if err := node.path.Validate(); err != nil {
		return fmt.Errorf("%w: invalid node path %q: %w", ErrInvalidNode, node.path, err)
	}

	pathID, err := node.path.ID()
	if err != nil {
		return fmt.Errorf("%w: invalid node path %q: %w", ErrInvalidNode, node.path, err)
	}

	if pathID != node.id {
		return fmt.Errorf(
			"%w: node ID %q must match path-derived ID %q",
			ErrInvalidNode,
			node.id,
			pathID,
		)
	}

	if node.use == "" {
		return fmt.Errorf("%w: non-root node %q must have a use segment", ErrInvalidNode, node.id)
	}

	if err := validateCommandNameSegment(node.use); err != nil {
		return fmt.Errorf("%w: invalid use segment %q: %w", ErrInvalidNode, node.use, err)
	}

	if node.path.Leaf() != node.use {
		return fmt.Errorf(
			"%w: node %q use segment %q must match path leaf %q",
			ErrInvalidNode,
			node.id,
			node.use,
			node.path.Leaf(),
		)
	}

	return node.validateAliases()
}

// validateAliases validates aliases for this node only.
func (node Node) validateAliases() error {
	seen := make(map[string]struct{}, len(node.aliases)+1)
	seen[node.use] = struct{}{}

	for _, alias := range node.aliases {
		if err := validateCommandNameSegment(alias); err != nil {
			return fmt.Errorf("%w: invalid alias %q: %w", ErrInvalidNode, alias, err)
		}

		if _, exists := seen[alias]; exists {
			return fmt.Errorf("%w: duplicate alias or use segment %q", ErrInvalidNode, alias)
		}

		seen[alias] = struct{}{}
	}

	return nil
}

// validateChildren validates direct children and sibling segment conflicts.
func (node Node) validateChildren() error {
	segments := make(map[string]Path)

	for _, child := range node.children {
		if child.kind == NodeRoot {
			return fmt.Errorf("%w: child of %q must not be a root node", ErrInvalidNode, node.path)
		}

		if err := child.Validate(); err != nil {
			return fmt.Errorf("%w: invalid child %q: %w", ErrInvalidNode, child.path, err)
		}

		if err := validateChildParent(node, child); err != nil {
			return err
		}

		if err := registerChildSegments(segments, child); err != nil {
			return err
		}
	}

	return nil
}

// validateChildParent verifies that child's parent path is node's path.
func validateChildParent(parent Node, child Node) error {
	childParent, ok := child.path.Parent()
	if !ok {
		return fmt.Errorf("%w: child %q has no parent path", ErrInvalidNode, child.path)
	}

	if !childParent.Equal(parent.path) {
		return fmt.Errorf(
			"%w: child %q parent path must be %q, got %q",
			ErrInvalidNode,
			child.path,
			parent.path,
			childParent,
		)
	}

	return nil
}

// registerChildSegments tracks direct child Use and alias segments to prevent
// ambiguous adapter-level resolution.
func registerChildSegments(segments map[string]Path, child Node) error {
	if err := registerChildSegment(segments, child.use, child.path); err != nil {
		return err
	}

	for _, alias := range child.aliases {
		if err := registerChildSegment(segments, alias, child.path); err != nil {
			return err
		}
	}

	return nil
}

// registerChildSegment registers one sibling segment.
func registerChildSegment(segments map[string]Path, segment string, childPath Path) error {
	if existing, exists := segments[segment]; exists {
		return fmt.Errorf(
			"%w: duplicate sibling command segment %q near %q and %q",
			ErrInvalidNode,
			segment,
			existing,
			childPath,
		)
	}

	segments[segment] = childPath

	return nil
}
