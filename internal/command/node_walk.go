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

// WalkFunc is called once for each visited command-tree node.
type WalkFunc func(Node) error

// Walk performs a pre-order traversal over node and its descendants.
func (node Node) Walk(fn WalkFunc) error {
	if fn == nil {
		return nil
	}

	if err := fn(node); err != nil {
		return err
	}

	for _, child := range node.children {
		if err := child.Walk(fn); err != nil {
			return err
		}
	}

	return nil
}

// FindByID returns the first node in the subtree with id.
//
// The zero ID is not a command identity and never matches, including the root.
func (node Node) FindByID(id ID) (Node, bool) {
	if id.IsZero() {
		return Node{}, false
	}

	if node.id == id {
		return node, true
	}

	for _, child := range node.children {
		if found, ok := child.FindByID(id); ok {
			return found, true
		}
	}

	return Node{}, false
}

// FindByPath returns the first node in the subtree with path.
func (node Node) FindByPath(path Path) (Node, bool) {
	if node.path.Equal(path) {
		return node, true
	}

	for _, child := range node.children {
		if found, ok := child.FindByPath(path); ok {
			return found, true
		}
	}

	return Node{}, false
}
