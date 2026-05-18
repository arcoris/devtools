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

// Nodes returns all tree nodes in pre-order.
func (tree Tree) Nodes() []Node {
	var nodes []Node

	_ = tree.Walk(func(_ Path, node Node) error {
		nodes = append(nodes, node)
		return nil
	})

	return nodes
}

// Families returns all family nodes in pre-order.
//
// The root node is not included because it is a root, not a family.
func (tree Tree) Families() []Node {
	var families []Node

	_ = tree.Walk(func(_ Path, node Node) error {
		if node.IsFamily() {
			families = append(families, node)
		}

		return nil
	})

	return families
}

// Commands returns all leaf command nodes in pre-order.
func (tree Tree) Commands() []Node {
	var commands []Node

	_ = tree.Walk(func(_ Path, node Node) error {
		if node.IsCommand() {
			commands = append(commands, node)
		}

		return nil
	})

	return commands
}

// Size returns the total number of nodes, including the root node.
func (tree Tree) Size() int {
	count := 0

	_ = tree.Walk(func(_ Path, _ Node) error {
		count++
		return nil
	})

	return count
}

// FamilyCount returns the number of family nodes.
//
// The root node is not counted as a family.
func (tree Tree) FamilyCount() int {
	count := 0

	_ = tree.Walk(func(_ Path, node Node) error {
		if node.IsFamily() {
			count++
		}

		return nil
	})

	return count
}

// CommandCount returns the number of leaf command nodes.
func (tree Tree) CommandCount() int {
	count := 0

	_ = tree.Walk(func(_ Path, node Node) error {
		if node.IsCommand() {
			count++
		}

		return nil
	})

	return count
}

// IDs returns all non-zero node IDs in pre-order.
//
// The root node has a zero ID and is intentionally omitted.
func (tree Tree) IDs() []ID {
	var ids []ID

	_ = tree.Walk(func(_ Path, node Node) error {
		if !node.ID().IsZero() {
			ids = append(ids, node.ID())
		}

		return nil
	})

	return ids
}

// Paths returns all node paths in pre-order, including the root path.
func (tree Tree) Paths() []Path {
	var paths []Path

	_ = tree.Walk(func(path Path, _ Node) error {
		paths = append(paths, path)
		return nil
	})

	return paths
}
