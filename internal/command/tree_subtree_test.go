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

import "testing"

func TestTreeSubtreeForRootReturnsTree(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	subtree, ok := tree.Subtree(RootPath())
	if !ok {
		t.Fatalf("Subtree(root) ok = false, want true")
	}

	if got, want := subtree.Size(), tree.Size(); got != want {
		t.Fatalf("Subtree(root) Size() = %d, want %d", got, want)
	}
}

func TestTreeSubtreeWrapsNonRootNode(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	subtree, ok := tree.Subtree(MustPath("bench"))
	if !ok {
		t.Fatalf("Subtree(bench) ok = false, want true")
	}

	if got, want := subtree.Size(), 4; got != want {
		t.Fatalf("Subtree(bench) Size() = %d, want %d", got, want)
	}

	if !subtree.Root().IsRoot() {
		t.Fatalf("Subtree(bench) root is not a synthetic root")
	}

	if !subtree.ContainsPath(MustPath("bench", "run")) {
		t.Fatalf("Subtree(bench) does not contain bench run")
	}
}

func TestTreeSubtreeMissingPath(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	if _, ok := tree.Subtree(MustPath("missing")); ok {
		t.Fatalf("Subtree(missing) ok = true, want false")
	}
}

func TestZeroTreeSubtreeMisses(t *testing.T) {
	t.Parallel()

	var tree Tree

	if _, ok := tree.Subtree(RootPath()); ok {
		t.Fatalf("zero tree Subtree(root) ok = true, want false")
	}
}
