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

// Validate verifies registry indexes and the underlying tree.
func (registry Registry) Validate() error {
	if registry.IsZero() {
		return fmt.Errorf("%w: tree is not set", ErrInvalidRegistry)
	}

	if err := registry.tree.Validate(); err != nil {
		return fmt.Errorf("%w: invalid tree: %w", ErrInvalidRegistry, err)
	}

	if registry.byID == nil {
		return fmt.Errorf("%w: ID index is not initialized", ErrInvalidRegistry)
	}

	if registry.byPath == nil {
		return fmt.Errorf("%w: path index is not initialized", ErrInvalidRegistry)
	}

	if err := registry.validateIndexes(); err != nil {
		return err
	}

	return nil
}

// validateIndexes verifies that all indexes match the underlying tree exactly.
func (registry Registry) validateIndexes() error {
	expectedSize := registry.tree.Size()

	if len(registry.byID) != expectedSize {
		return fmt.Errorf(
			"%w: ID index size %d does not match tree size %d",
			ErrInvalidRegistry,
			len(registry.byID),
			expectedSize,
		)
	}

	if len(registry.byPath) != expectedSize {
		return fmt.Errorf(
			"%w: path index size %d does not match tree size %d",
			ErrInvalidRegistry,
			len(registry.byPath),
			expectedSize,
		)
	}

	expectedIDs, expectedPaths, err := registry.expectedIndexOrder()
	if err != nil {
		return err
	}

	if err := registry.validateOrderedIDs(expectedIDs); err != nil {
		return err
	}

	if err := registry.validateOrderedPaths(expectedPaths); err != nil {
		return err
	}

	return registry.validateIndexedNodes()
}

// expectedIndexOrder returns the canonical traversal order for index slices.
func (registry Registry) expectedIndexOrder() ([]ID, []Path, error) {
	ids := make([]ID, 0, registry.tree.Size())
	paths := make([]Path, 0, registry.tree.Size())

	err := registry.tree.Walk(func(path Path, node Node) error {
		ids = append(ids, node.ID())
		paths = append(paths, path.clone())

		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%w: cannot compute expected index order: %w", ErrInvalidRegistry, err)
	}

	return ids, paths, nil
}

// validateOrderedIDs verifies the stable pre-order ID list.
func (registry Registry) validateOrderedIDs(expected []ID) error {
	if len(registry.ids) != len(expected) {
		return fmt.Errorf(
			"%w: ordered ID list size %d does not match tree size %d",
			ErrInvalidRegistry,
			len(registry.ids),
			len(expected),
		)
	}

	for index, want := range expected {
		if got := registry.ids[index]; got != want {
			return fmt.Errorf(
				"%w: ordered ID list mismatch at index %d: got %q, want %q",
				ErrInvalidRegistry,
				index,
				got,
				want,
			)
		}
	}

	return nil
}

// validateOrderedPaths verifies the stable pre-order path list.
func (registry Registry) validateOrderedPaths(expected []Path) error {
	if len(registry.paths) != len(expected) {
		return fmt.Errorf(
			"%w: ordered path list size %d does not match tree size %d",
			ErrInvalidRegistry,
			len(registry.paths),
			len(expected),
		)
	}

	for index, want := range expected {
		if got := registry.paths[index]; !got.Equal(want) {
			return fmt.Errorf(
				"%w: ordered path list mismatch at index %d: got %q, want %q",
				ErrInvalidRegistry,
				index,
				got,
				want,
			)
		}
	}

	return nil
}

// validateIndexedNodes verifies every tree node is addressable through both
// maps and that map entries point back to the same logical node.
func (registry Registry) validateIndexedNodes() error {
	return registry.tree.Walk(func(path Path, node Node) error {
		byPath, ok := registry.byPath[path.Key()]
		if !ok {
			return fmt.Errorf("%w: path %q is missing from path index", ErrInvalidRegistry, path)
		}

		if !sameNodeIdentity(byPath, node) {
			return fmt.Errorf(
				"%w: path index mismatch for %q: got ID %q path %q, want ID %q path %q",
				ErrInvalidRegistry,
				path,
				byPath.ID(),
				byPath.Path(),
				node.ID(),
				node.Path(),
			)
		}

		byID, ok := registry.byID[node.ID()]
		if !ok {
			return fmt.Errorf("%w: ID %q is missing from ID index", ErrInvalidRegistry, node.ID())
		}

		if !sameNodeIdentity(byID, node) {
			return fmt.Errorf(
				"%w: ID index mismatch for %q: got ID %q path %q, want ID %q path %q",
				ErrInvalidRegistry,
				node.ID(),
				byID.ID(),
				byID.Path(),
				node.ID(),
				node.Path(),
			)
		}

		return nil
	})
}

// sameNodeIdentity compares the stable registry identity of two nodes.
func sameNodeIdentity(left Node, right Node) bool {
	return left.ID() == right.ID() && left.Path().Equal(right.Path())
}
