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

// Subtree returns a Tree rooted at path.
//
// The returned tree keeps the selected node as its root only if path is the
// command-tree root. For non-root paths, the selected node is wrapped into a
// synthetic root so the returned value remains a valid Tree.
func (tree Tree) Subtree(path Path) (Tree, bool) {
	node, ok := tree.FindByPath(path)
	if !ok {
		return Tree{}, false
	}

	if node.IsRoot() {
		return tree, true
	}

	subtree, err := NewTreeFromChildren(node)
	if err != nil {
		return Tree{}, false
	}

	return subtree, true
}
