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

// Validate verifies the tree structure and global invariants.
func (tree Tree) Validate() error {
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

	if err := tree.root.Validate(); err != nil {
		return fmt.Errorf("%w: invalid root node: %w", ErrInvalidTree, err)
	}

	if err := tree.validateUniquePaths(); err != nil {
		return err
	}

	if err := tree.validateUniqueIDs(); err != nil {
		return err
	}

	return nil
}

// validateUniquePaths verifies that no two nodes share the same Path.
//
// Valid Node values make duplicate paths impossible through parent/child and
// sibling checks. This tree-level pass remains as a defensive invariant around
// package-internal construction and future bulk import paths.
func (tree Tree) validateUniquePaths() error {
	seen := make(map[string]Path)

	return tree.walkRawPreOrder(func(path Path, node Node) error {
		key := path.Key()
		if previous, exists := seen[key]; exists {
			return fmt.Errorf(
				"%w: duplicate path %q already registered as %q",
				ErrInvalidTree,
				path,
				previous,
			)
		}

		seen[key] = path

		return nil
	})
}

// validateUniqueIDs verifies that no two non-root nodes share the same ID.
//
// IDs are path-derived for valid non-root nodes, so this check mirrors the path
// uniqueness pass while producing ID-focused diagnostics.
func (tree Tree) validateUniqueIDs() error {
	seen := make(map[ID]Path)

	return tree.walkRawPreOrder(func(path Path, node Node) error {
		id := node.ID()
		if id.IsZero() {
			return nil
		}

		if previousPath, exists := seen[id]; exists {
			return fmt.Errorf(
				"%w: duplicate ID %q already registered at path %q; duplicate path %q",
				ErrInvalidTree,
				id,
				previousPath,
				path,
			)
		}

		seen[id] = path

		return nil
	})
}

// walkRawPreOrder traverses without validating tree first.
//
// It is used by validation itself after basic root checks have passed.
func (tree Tree) walkRawPreOrder(fn TreeWalkFunc) error {
	return walkNodePreOrder(tree.root, fn)
}
