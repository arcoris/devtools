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

func TestTreeFindByPath(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	node, ok := tree.FindByPath(MustPath("bench", "run"))
	if !ok {
		t.Fatalf("FindByPath() ok = false, want true")
	}

	if got, want := node.ID(), MustID("bench.run"); got != want {
		t.Fatalf("FindByPath() ID = %q, want %q", got, want)
	}

	if _, ok := tree.FindByPath(MustPath("profile", "cpu")); ok {
		t.Fatalf("FindByPath() ok = true, want false")
	}
}

func TestTreeMustFindByPathPanics(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	assertPanics(t, func() {
		_ = tree.MustFindByPath(MustPath("missing"))
	})
}

func TestTreeFindByID(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	node, ok := tree.FindByID(MustID("bench.compare"))
	if !ok {
		t.Fatalf("FindByID() ok = false, want true")
	}

	if got, want := node.Path(), MustPath("bench", "compare"); !got.Equal(want) {
		t.Fatalf("FindByID() Path = %q, want %q", got, want)
	}

	if _, ok := tree.FindByID(MustID("profile.cpu")); ok {
		t.Fatalf("FindByID() ok = true, want false")
	}
}

func TestTreeFindByIDCanFindRoot(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	node, ok := tree.FindByID("")
	if !ok {
		t.Fatalf("FindByID(zero) ok = false, want true")
	}

	if !node.IsRoot() {
		t.Fatalf("FindByID(zero) returned non-root node")
	}
}

func TestTreeMustFindByIDPanics(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	assertPanics(t, func() {
		_ = tree.MustFindByID(MustID("missing"))
	})
}

func TestTreeResolveUsesAliases(t *testing.T) {
	t.Parallel()

	run := MustNode(NodeSpec{
		Kind:    NodeCommand,
		ID:      MustID("bench.run"),
		Path:    MustPath("bench", "run"),
		Use:     "run",
		Aliases: []string{"execute"},
	})

	bench := MustFamilyNode(MustID("bench"), MustPath("bench"), "bench", run)
	tree := MustTreeFromChildren(bench)

	node, ok := tree.Resolve("bench", "execute")
	if !ok {
		t.Fatalf("Resolve() ok = false, want true")
	}

	if got, want := node.ID(), MustID("bench.run"); got != want {
		t.Fatalf("Resolve() ID = %q, want %q", got, want)
	}

	if _, ok := tree.Resolve("bench", "missing"); ok {
		t.Fatalf("Resolve() ok = true, want false")
	}
}

func TestTreeResolveCanReturnRoot(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	node, ok := tree.Resolve()
	if !ok {
		t.Fatalf("Resolve() ok = false, want true")
	}

	if !node.IsRoot() {
		t.Fatalf("Resolve() returned non-root node")
	}
}

func TestTreeResolvePathUsesCanonicalSegmentsOnly(t *testing.T) {
	t.Parallel()

	run := MustNode(NodeSpec{
		Kind:    NodeCommand,
		ID:      MustID("bench.run"),
		Path:    MustPath("bench", "run"),
		Use:     "run",
		Aliases: []string{"execute"},
	})

	bench := MustFamilyNode(MustID("bench"), MustPath("bench"), "bench", run)
	tree := MustTreeFromChildren(bench)

	node, ok := tree.ResolvePath(MustPath("bench", "run"))
	if !ok {
		t.Fatalf("ResolvePath(canonical) ok = false, want true")
	}

	if got, want := node.ID(), MustID("bench.run"); got != want {
		t.Fatalf("ResolvePath(canonical) ID = %q, want %q", got, want)
	}

	if _, ok := tree.ResolvePath(MustPath("bench", "execute")); ok {
		t.Fatalf("ResolvePath(alias) ok = true, want false")
	}
}

func TestTreeContainsHelpers(t *testing.T) {
	t.Parallel()

	tree := mustTestCommandTree(t)

	if !tree.ContainsPath(MustPath("bench", "run")) {
		t.Fatalf("ContainsPath(bench run) = false, want true")
	}

	if tree.ContainsPath(MustPath("missing")) {
		t.Fatalf("ContainsPath(missing) = true, want false")
	}

	if !tree.ContainsID(MustID("bench.run")) {
		t.Fatalf("ContainsID(bench.run) = false, want true")
	}

	if tree.ContainsID(MustID("missing")) {
		t.Fatalf("ContainsID(missing) = true, want false")
	}
}

func TestZeroTreeLookupsMiss(t *testing.T) {
	t.Parallel()

	var tree Tree

	if _, ok := tree.FindByPath(MustPath("bench")); ok {
		t.Fatalf("zero tree FindByPath() ok = true, want false")
	}

	if _, ok := tree.FindByID(MustID("bench")); ok {
		t.Fatalf("zero tree FindByID() ok = true, want false")
	}

	if _, ok := tree.Resolve("bench"); ok {
		t.Fatalf("zero tree Resolve() ok = true, want false")
	}

	if _, ok := tree.ResolvePath(MustPath("bench")); ok {
		t.Fatalf("zero tree ResolvePath() ok = true, want false")
	}
}
