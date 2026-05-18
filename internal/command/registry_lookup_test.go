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

func TestRegistryFindByID(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	node, ok := registry.FindByID(MustID("bench.run"))
	if !ok {
		t.Fatalf("FindByID() ok = false, want true")
	}

	if got, want := node.Path(), MustPath("bench", "run"); !got.Equal(want) {
		t.Fatalf("FindByID() Path() = %q, want %q", got, want)
	}

	if _, ok := registry.FindByID(MustID("profile.cpu")); ok {
		t.Fatalf("FindByID(missing) ok = true, want false")
	}
}

func TestRegistryFindByIDCanFindRoot(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	node, ok := registry.FindByID("")
	if !ok {
		t.Fatalf("FindByID(zero) ok = false, want true")
	}

	if !node.IsRoot() {
		t.Fatalf("FindByID(zero) returned non-root node")
	}
}

func TestRegistryMustFindByIDPanics(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	assertPanics(t, func() {
		_ = registry.MustFindByID(MustID("missing"))
	})
}

func TestRegistryFindByPath(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	node, ok := registry.FindByPath(MustPath("bench", "compare"))
	if !ok {
		t.Fatalf("FindByPath() ok = false, want true")
	}

	if got, want := node.ID(), MustID("bench.compare"); got != want {
		t.Fatalf("FindByPath() ID() = %q, want %q", got, want)
	}

	if _, ok := registry.FindByPath(MustPath("profile", "cpu")); ok {
		t.Fatalf("FindByPath(missing) ok = true, want false")
	}
}

func TestRegistryMustFindByPathPanics(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	assertPanics(t, func() {
		_ = registry.MustFindByPath(MustPath("missing"))
	})
}

func TestRegistryContainsHelpers(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	if !registry.ContainsID(MustID("bench.run")) {
		t.Fatalf("ContainsID(bench.run) = false, want true")
	}

	if registry.ContainsID(MustID("missing")) {
		t.Fatalf("ContainsID(missing) = true, want false")
	}

	if !registry.ContainsPath(MustPath("bench", "run")) {
		t.Fatalf("ContainsPath(bench run) = false, want true")
	}

	if registry.ContainsPath(MustPath("missing")) {
		t.Fatalf("ContainsPath(missing) = true, want false")
	}
}

func TestZeroRegistryLookupsMiss(t *testing.T) {
	t.Parallel()

	var registry Registry

	if _, ok := registry.FindByID(MustID("bench")); ok {
		t.Fatalf("zero registry FindByID() ok = true, want false")
	}

	if _, ok := registry.FindByPath(MustPath("bench")); ok {
		t.Fatalf("zero registry FindByPath() ok = true, want false")
	}
}
