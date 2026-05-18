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

// FindByPath returns the node with path.
func (tree Tree) FindByPath(path Path) (Node, bool) {
	if tree.IsZero() {
		return Node{}, false
	}

	return tree.root.FindByPath(path)
}

// MustFindByPath returns the node with path and panics when it is not found.
func (tree Tree) MustFindByPath(path Path) Node {
	node, ok := tree.FindByPath(path)
	if !ok {
		panic(fmt.Sprintf("command tree path %q was not found", path))
	}

	return node
}

// ContainsPath reports whether the tree contains path.
func (tree Tree) ContainsPath(path Path) bool {
	_, ok := tree.FindByPath(path)

	return ok
}

// FindByID returns the node with id.
//
// The zero ID is treated as the root ID for lookup purposes.
func (tree Tree) FindByID(id ID) (Node, bool) {
	if tree.IsZero() {
		return Node{}, false
	}

	if id.IsZero() {
		return tree.root, true
	}

	return tree.root.FindByID(id)
}

// MustFindByID returns the node with id and panics when it is not found.
func (tree Tree) MustFindByID(id ID) Node {
	node, ok := tree.FindByID(id)
	if !ok {
		panic(fmt.Sprintf("command tree ID %q was not found", id))
	}

	return node
}

// ContainsID reports whether the tree contains id.
func (tree Tree) ContainsID(id ID) bool {
	_, ok := tree.FindByID(id)

	return ok
}

// Resolve resolves a path by walking through child Use and alias segments.
//
// Resolve is adapter-oriented: it accepts aliases. FindByPath is
// structure-oriented: it matches canonical paths only.
func (tree Tree) Resolve(segments ...string) (Node, bool) {
	if tree.IsZero() {
		return Node{}, false
	}

	current := tree.root
	for _, segment := range segments {
		next, ok := current.FindChildBySegment(segment)
		if !ok {
			return Node{}, false
		}

		current = next
	}

	return current, true
}

// ResolvePath resolves a Path by canonical Use segments only.
//
// ResolvePath does not use aliases. Use Resolve when adapter-level alias
// resolution is desired.
func (tree Tree) ResolvePath(path Path) (Node, bool) {
	if tree.IsZero() {
		return Node{}, false
	}

	current := tree.root
	for _, segment := range path.segments {
		next, ok := current.FindChildByUse(segment)
		if !ok {
			return Node{}, false
		}

		current = next
	}

	return current, true
}
