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

func TestRegistrySubtreeForRootReturnsRegistry(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	subtree, ok := registry.Subtree(RootPath())
	if !ok {
		t.Fatalf("Subtree(root) ok = false, want true")
	}

	if got, want := subtree.Size(), registry.Size(); got != want {
		t.Fatalf("Subtree(root) Size() = %d, want %d", got, want)
	}
}

func TestRegistrySubtreeWrapsNonRootNode(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	subtree, ok := registry.Subtree(MustPath("bench"))
	if !ok {
		t.Fatalf("Subtree(bench) ok = false, want true")
	}

	if got, want := subtree.Size(), 4; got != want {
		t.Fatalf("Subtree(bench) Size() = %d, want %d", got, want)
	}

	if !subtree.ContainsID(MustID("bench.run")) {
		t.Fatalf("subtree does not contain bench.run")
	}

	if !subtree.Root().IsRoot() {
		t.Fatalf("subtree root is not synthetic root")
	}
}

func TestRegistrySubtreeMissingPath(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	if _, ok := registry.Subtree(MustPath("missing")); ok {
		t.Fatalf("Subtree(missing) ok = true, want false")
	}
}

func TestZeroRegistrySubtreeMisses(t *testing.T) {
	t.Parallel()

	var registry Registry

	if _, ok := registry.Subtree(RootPath()); ok {
		t.Fatalf("zero registry Subtree(root) ok = true, want false")
	}
}
