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

// NewRegistry validates tree, builds indexes, and returns a Registry.
func NewRegistry(tree Tree) (Registry, error) {
	if err := tree.Validate(); err != nil {
		return Registry{}, fmt.Errorf("%w: invalid tree: %w", ErrInvalidRegistry, err)
	}

	registry := Registry{tree: tree}
	if err := registry.rebuildIndexes(); err != nil {
		return Registry{}, err
	}

	if err := registry.Validate(); err != nil {
		return Registry{}, err
	}

	return registry, nil
}

// MustRegistry validates tree, builds indexes, and panics on error.
//
// MustRegistry is intended for static command-definition wiring and tests.
func MustRegistry(tree Tree) Registry {
	registry, err := NewRegistry(tree)
	if err != nil {
		panic(err)
	}

	return registry
}

// NewRegistryFromRoot validates root, builds a Tree, and returns a Registry.
func NewRegistryFromRoot(root Node) (Registry, error) {
	tree, err := NewTree(root)
	if err != nil {
		return Registry{}, fmt.Errorf("%w: invalid root: %w", ErrInvalidRegistry, err)
	}

	return NewRegistry(tree)
}

// MustRegistryFromRoot validates root, builds a Registry, and panics on error.
func MustRegistryFromRoot(root Node) Registry {
	registry, err := NewRegistryFromRoot(root)
	if err != nil {
		panic(err)
	}

	return registry
}

// NewRegistryFromChildren builds a root tree from top-level children and returns
// a Registry.
func NewRegistryFromChildren(children ...Node) (Registry, error) {
	tree, err := NewTreeFromChildren(children...)
	if err != nil {
		return Registry{}, fmt.Errorf("%w: invalid children: %w", ErrInvalidRegistry, err)
	}

	return NewRegistry(tree)
}

// MustRegistryFromChildren builds a Registry from top-level children and panics
// on error.
func MustRegistryFromChildren(children ...Node) Registry {
	registry, err := NewRegistryFromChildren(children...)
	if err != nil {
		panic(err)
	}

	return registry
}

// EmptyRegistry returns a valid registry containing only the root node.
func EmptyRegistry() Registry {
	return MustRegistryFromChildren()
}
