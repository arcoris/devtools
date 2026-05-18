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

// IsZero reports whether the registry has not been initialized.
func (registry Registry) IsZero() bool {
	return registry.tree.IsZero()
}

// IsEmpty reports whether the registry has no top-level command nodes.
func (registry Registry) IsEmpty() bool {
	if registry.IsZero() {
		return true
	}

	return registry.tree.IsEmpty()
}

// Tree returns the registry command tree.
func (registry Registry) Tree() Tree {
	return registry.tree
}

// Root returns the registry root node.
func (registry Registry) Root() Node {
	return registry.tree.Root()
}
