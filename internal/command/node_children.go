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

// Children returns a detached copy of child nodes.
func (node Node) Children() []Node {
	return cloneNodes(node.children)
}

// ChildCount returns the number of child nodes.
func (node Node) ChildCount() int {
	return len(node.children)
}

// FindChildByUse returns the direct child with the given canonical Use segment.
//
// Aliases are intentionally ignored. Use FindChildBySegment when adapter-level
// alias resolution is desired.
func (node Node) FindChildByUse(use string) (Node, bool) {
	for _, child := range node.children {
		if child.use == use {
			return child, true
		}
	}

	return Node{}, false
}

// FindChildBySegment returns the direct child matching segment by Use or alias.
//
// Validation prevents sibling Use/alias conflicts, so this method returns at
// most one child in a valid tree.
func (node Node) FindChildBySegment(segment string) (Node, bool) {
	for _, child := range node.children {
		if child.MatchesSegment(segment) {
			return child, true
		}
	}

	return Node{}, false
}

// WithChildren returns a copy of node with children replaced.
//
// WithChildren never modifies the receiver. The returned node is validated.
func (node Node) WithChildren(children ...Node) (Node, error) {
	spec := node.spec()
	spec.Children = children

	return NewNode(spec)
}

// MustWithChildren returns a copy of node with children replaced and panics on
// validation errors.
func (node Node) MustWithChildren(children ...Node) Node {
	next, err := node.WithChildren(children...)
	if err != nil {
		panic(err)
	}

	return next
}

// AppendChild returns a copy of node with child appended.
//
// AppendChild never modifies the receiver. The returned node is validated.
func (node Node) AppendChild(child Node) (Node, error) {
	children := make([]Node, 0, len(node.children)+1)
	children = append(children, node.children...)
	children = append(children, child)

	return node.WithChildren(children...)
}

// MustAppendChild returns a copy of node with child appended and panics on
// validation errors.
func (node Node) MustAppendChild(child Node) Node {
	next, err := node.AppendChild(child)
	if err != nil {
		panic(err)
	}

	return next
}
