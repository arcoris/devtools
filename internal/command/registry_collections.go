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

// Size returns the total number of nodes, including the root node.
func (registry Registry) Size() int {
	return registry.tree.Size()
}

// CommandCount returns the number of leaf command nodes.
func (registry Registry) CommandCount() int {
	return registry.tree.CommandCount()
}

// FamilyCount returns the number of family nodes.
//
// The root node is not counted as a family.
func (registry Registry) FamilyCount() int {
	return registry.tree.FamilyCount()
}

// TopLevel returns the root's direct children.
func (registry Registry) TopLevel() []Node {
	return registry.tree.TopLevel()
}

// Nodes returns all nodes in deterministic pre-order.
func (registry Registry) Nodes() []Node {
	return registry.tree.Nodes()
}

// Commands returns all leaf command nodes in deterministic pre-order.
func (registry Registry) Commands() []Node {
	return registry.tree.Commands()
}

// Families returns all family nodes in deterministic pre-order.
func (registry Registry) Families() []Node {
	return registry.tree.Families()
}

// IDs returns all indexed IDs in deterministic pre-order.
//
// The zero ID used by the root node is included as the first entry. Use
// NonRootIDs when only command/family IDs are needed.
func (registry Registry) IDs() []ID {
	out := make([]ID, len(registry.ids))
	copy(out, registry.ids)

	return out
}

// NonRootIDs returns all non-zero IDs in deterministic pre-order.
func (registry Registry) NonRootIDs() []ID {
	ids := make([]ID, 0, len(registry.ids))

	for _, id := range registry.ids {
		if id.IsZero() {
			continue
		}

		ids = append(ids, id)
	}

	return ids
}

// Paths returns all indexed canonical paths in deterministic pre-order.
//
// The root path is included as the first entry.
func (registry Registry) Paths() []Path {
	return clonePaths(registry.paths)
}
