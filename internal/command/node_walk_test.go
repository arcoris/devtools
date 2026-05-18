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

import (
	"errors"
	"testing"
)

func TestNodeWalk(t *testing.T) {
	t.Parallel()

	root := MustRootNode(
		mustTestFamilyNode(
			t,
			"bench",
			mustTestCommandNode(t, "bench.run"),
			mustTestCommandNode(t, "bench.compare"),
		),
	)

	var got []string
	err := root.Walk(func(node Node) error {
		got = append(got, node.Path().String())
		return nil
	})
	if err != nil {
		t.Fatalf("Walk() returned unexpected error: %v", err)
	}

	want := []string{"", "bench", "bench run", "bench compare"}
	if len(got) != len(want) {
		t.Fatalf("Walk visited %d nodes, want %d: got %v", len(got), len(want), got)
	}

	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("Walk[%d] = %q, want %q", index, got[index], want[index])
		}
	}
}

func TestNodeWalkNilFuncIsNoop(t *testing.T) {
	t.Parallel()

	root := MustRootNode(mustTestCommandNode(t, "check"))

	if err := root.Walk(nil); err != nil {
		t.Fatalf("Walk(nil) error = %v, want nil", err)
	}
}

func TestNodeWalkStopsOnError(t *testing.T) {
	t.Parallel()

	expected := errors.New("stop")
	root := MustRootNode(mustTestCommandNode(t, "check"))

	err := root.Walk(func(Node) error {
		return expected
	})

	if !errors.Is(err, expected) {
		t.Fatalf("Walk() error = %v, want %v", err, expected)
	}
}

func TestNodeFindByID(t *testing.T) {
	t.Parallel()

	root := MustRootNode(
		mustTestFamilyNode(t, "bench", mustTestCommandNode(t, "bench.run")),
		mustTestCommandNode(t, "check"),
	)

	found, ok := root.FindByID(MustID("bench.run"))
	if !ok {
		t.Fatalf("FindByID(%q) ok = false, want true", "bench.run")
	}

	if got, want := found.Path(), MustPath("bench", "run"); !got.Equal(want) {
		t.Fatalf("FindByID path = %q, want %q", got, want)
	}

	if _, ok := root.FindByID(MustID("missing")); ok {
		t.Fatalf("FindByID(missing) ok = true, want false")
	}

	if _, ok := root.FindByID(""); ok {
		t.Fatalf("FindByID(zero) ok = true, want false")
	}
}

func TestNodeFindByPath(t *testing.T) {
	t.Parallel()

	root := MustRootNode(mustTestFamilyNode(t, "bench", mustTestCommandNode(t, "bench.run")))

	found, ok := root.FindByPath(MustPath("bench", "run"))
	if !ok {
		t.Fatalf("FindByPath(%q) ok = false, want true", "bench run")
	}

	if got, want := found.ID(), MustID("bench.run"); got != want {
		t.Fatalf("FindByPath ID = %q, want %q", got, want)
	}

	if found, ok := root.FindByPath(RootPath()); !ok || !found.IsRoot() {
		t.Fatalf("FindByPath(root) = %q, %v; want root, true", found.Path(), ok)
	}

	if _, ok := root.FindByPath(MustPath("profile")); ok {
		t.Fatalf("FindByPath(profile) ok = true, want false")
	}
}
