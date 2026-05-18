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

func TestTreeCollections(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	if got, want := len(tree.Nodes()), 4; got != want {
		t.Fatalf("len(Nodes()) = %d, want %d", got, want)
	}

	if got, want := len(tree.Families()), 1; got != want {
		t.Fatalf("len(Families()) = %d, want %d", got, want)
	}

	if got, want := len(tree.Commands()), 2; got != want {
		t.Fatalf("len(Commands()) = %d, want %d", got, want)
	}
}

func TestTreeCollectionOrder(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	assertStringSlicesEqual(t, pathStrings(tree.Paths()), []string{"", "bench", "bench run", "bench compare"})
	assertStringSlicesEqual(t, idStrings(tree.IDs()), []string{"bench", "bench.run", "bench.compare"})
}

func TestTopLevelReturnsDetachedCopy(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	top := tree.TopLevel()
	if got, want := len(top), 1; got != want {
		t.Fatalf("len(TopLevel()) = %d, want %d", got, want)
	}

	top[0] = mustTestCommandNode(t, "check")

	again := tree.TopLevel()
	if got, want := again[0].ID(), MustID("bench"); got != want {
		t.Fatalf("TopLevel returned mutable state: got %q, want %q", got, want)
	}
}

func TestZeroTreeCollectionsAreEmpty(t *testing.T) {
	t.Parallel()

	var tree Tree

	if got := tree.Nodes(); len(got) != 0 {
		t.Fatalf("zero tree Nodes() = %v, want empty", got)
	}

	if got := tree.Families(); len(got) != 0 {
		t.Fatalf("zero tree Families() = %v, want empty", got)
	}

	if got := tree.Commands(); len(got) != 0 {
		t.Fatalf("zero tree Commands() = %v, want empty", got)
	}

	if got := tree.TopLevel(); len(got) != 0 {
		t.Fatalf("zero tree TopLevel() = %v, want empty", got)
	}
}
