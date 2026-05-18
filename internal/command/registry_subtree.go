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

// Subtree returns a Registry containing the subtree rooted at path.
//
// The returned registry is re-indexed and validated. For non-root paths, the
// selected node is wrapped into a synthetic root by Tree.Subtree.
func (registry Registry) Subtree(path Path) (Registry, bool) {
	subtree, ok := registry.tree.Subtree(path)
	if !ok {
		return Registry{}, false
	}

	next, err := NewRegistry(subtree)
	if err != nil {
		return Registry{}, false
	}

	return next, true
}
