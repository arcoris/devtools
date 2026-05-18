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

func TestRegistryCollections(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	if got, want := len(registry.Nodes()), 4; got != want {
		t.Fatalf("len(Nodes()) = %d, want %d", got, want)
	}

	if got, want := len(registry.Families()), 1; got != want {
		t.Fatalf("len(Families()) = %d, want %d", got, want)
	}

	if got, want := len(registry.Commands()), 2; got != want {
		t.Fatalf("len(Commands()) = %d, want %d", got, want)
	}
}

func TestRegistryIDs(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	assertStringSlicesEqual(t, idStrings(registry.IDs()), []string{"", "bench", "bench.run", "bench.compare"})
}

func TestRegistryNonRootIDs(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	assertStringSlicesEqual(t, idStrings(registry.NonRootIDs()), []string{"bench", "bench.run", "bench.compare"})
}

func TestRegistryPaths(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	assertStringSlicesEqual(t, pathStrings(registry.Paths()), []string{"", "bench", "bench run", "bench compare"})
}

func TestRegistryTopLevelReturnsDetachedCopy(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	topLevel := registry.TopLevel()
	if got, want := len(topLevel), 1; got != want {
		t.Fatalf("len(TopLevel()) = %d, want %d", got, want)
	}

	topLevel[0] = mustTestCommandNode(t, "check")

	again := registry.TopLevel()
	if got, want := again[0].ID(), MustID("bench"); got != want {
		t.Fatalf("TopLevel() returned mutable state: got %q, want %q", got, want)
	}
}

func TestRegistryIDsReturnsDetachedCopy(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	ids := registry.IDs()
	ids[1] = MustID("changed")

	again := registry.IDs()
	if got, want := again[1], MustID("bench"); got != want {
		t.Fatalf("IDs() returned mutable state: got %q, want %q", got, want)
	}
}

func TestRegistryPathsReturnsDeepDetachedCopy(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	paths := registry.Paths()
	paths[1].segments[0] = "changed"

	again := registry.Paths()
	if got, want := again[1], MustPath("bench"); !got.Equal(want) {
		t.Fatalf("Paths() returned mutable state: got %q, want %q", got, want)
	}
}

func TestZeroRegistryCollectionsAreEmpty(t *testing.T) {
	t.Parallel()

	var registry Registry

	if got := registry.IDs(); len(got) != 0 {
		t.Fatalf("zero registry IDs() = %v, want empty", got)
	}

	if got := registry.Paths(); len(got) != 0 {
		t.Fatalf("zero registry Paths() = %v, want empty", got)
	}

	if got := registry.TopLevel(); len(got) != 0 {
		t.Fatalf("zero registry TopLevel() = %v, want empty", got)
	}
}
