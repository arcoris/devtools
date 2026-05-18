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

// rebuildIndexes builds deterministic lookup indexes from the registry tree.
//
// Index order follows tree pre-order traversal.
func (registry *Registry) rebuildIndexes() error {
	registry.byID = make(map[ID]Node)
	registry.byPath = make(map[string]Node)
	registry.ids = nil
	registry.paths = nil

	err := registry.tree.Walk(func(path Path, node Node) error {
		if err := registry.registerPath(path, node); err != nil {
			return err
		}

		if err := registry.registerID(node.ID(), node); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("%w: failed to build indexes: %w", ErrInvalidRegistry, err)
	}

	return nil
}

// registerPath registers one canonical path in the path index.
func (registry *Registry) registerPath(path Path, node Node) error {
	key := path.Key()

	if existing, exists := registry.byPath[key]; exists {
		return fmt.Errorf(
			"%w: duplicate path %q already registered by node %q; duplicate node %q",
			ErrInvalidRegistry,
			path,
			existing.ID(),
			node.ID(),
		)
	}

	registry.byPath[key] = node
	registry.paths = append(registry.paths, path.clone())

	return nil
}

// registerID registers one ID in the ID index.
//
// The root node uses the zero ID and is indexed as such. Non-root duplicates
// should already be rejected by Tree, but Registry keeps its own check to
// protect index integrity and package-internal construction paths.
func (registry *Registry) registerID(id ID, node Node) error {
	if existing, exists := registry.byID[id]; exists {
		return fmt.Errorf(
			"%w: duplicate ID %q already registered at path %q; duplicate path %q",
			ErrInvalidRegistry,
			id,
			existing.Path(),
			node.Path(),
		)
	}

	registry.byID[id] = node
	registry.ids = append(registry.ids, id)

	return nil
}
