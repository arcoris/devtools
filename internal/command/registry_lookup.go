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

// FindByID returns the node with id.
//
// The zero ID resolves to the root node.
func (registry Registry) FindByID(id ID) (Node, bool) {
	if registry.byID == nil {
		return Node{}, false
	}

	node, ok := registry.byID[id]

	return node, ok
}

// MustFindByID returns the node with id and panics when it is not found.
func (registry Registry) MustFindByID(id ID) Node {
	node, ok := registry.FindByID(id)
	if !ok {
		panic(fmt.Sprintf("command registry ID %q was not found", id))
	}

	return node
}

// ContainsID reports whether the registry contains id.
func (registry Registry) ContainsID(id ID) bool {
	_, ok := registry.FindByID(id)

	return ok
}

// FindByPath returns the node with canonical path.
func (registry Registry) FindByPath(path Path) (Node, bool) {
	if registry.byPath == nil {
		return Node{}, false
	}

	node, ok := registry.byPath[path.Key()]

	return node, ok
}

// MustFindByPath returns the node with path and panics when it is not found.
func (registry Registry) MustFindByPath(path Path) Node {
	node, ok := registry.FindByPath(path)
	if !ok {
		panic(fmt.Sprintf("command registry path %q was not found", path))
	}

	return node
}

// ContainsPath reports whether the registry contains path.
func (registry Registry) ContainsPath(path Path) bool {
	_, ok := registry.FindByPath(path)

	return ok
}
