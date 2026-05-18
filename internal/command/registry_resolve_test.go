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

func TestRegistryResolveUsesAliases(t *testing.T) {
	t.Parallel()

	registry := mustTestAliasRegistry(t)

	node, ok := registry.Resolve("bench", "execute")
	if !ok {
		t.Fatalf("Resolve() ok = false, want true")
	}

	if got, want := node.ID(), MustID("bench.run"); got != want {
		t.Fatalf("Resolve() ID() = %q, want %q", got, want)
	}

	if _, ok := registry.Resolve("bench", "missing"); ok {
		t.Fatalf("Resolve(missing) ok = true, want false")
	}
}

func TestRegistryResolveCanReturnRoot(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	node, ok := registry.Resolve()
	if !ok {
		t.Fatalf("Resolve() ok = false, want true")
	}

	if !node.IsRoot() {
		t.Fatalf("Resolve() returned non-root node")
	}
}

func TestRegistryResolvePathUsesCanonicalSegmentsOnly(t *testing.T) {
	t.Parallel()

	registry := mustTestAliasRegistry(t)

	node, ok := registry.ResolvePath(MustPath("bench", "run"))
	if !ok {
		t.Fatalf("ResolvePath(canonical) ok = false, want true")
	}

	if got, want := node.ID(), MustID("bench.run"); got != want {
		t.Fatalf("ResolvePath(canonical) ID() = %q, want %q", got, want)
	}

	if _, ok := registry.ResolvePath(MustPath("bench", "execute")); ok {
		t.Fatalf("ResolvePath(alias) ok = true, want false")
	}
}

func TestZeroRegistryResolveMisses(t *testing.T) {
	t.Parallel()

	var registry Registry

	if _, ok := registry.Resolve("bench"); ok {
		t.Fatalf("zero registry Resolve() ok = true, want false")
	}

	if _, ok := registry.ResolvePath(MustPath("bench")); ok {
		t.Fatalf("zero registry ResolvePath() ok = true, want false")
	}
}
