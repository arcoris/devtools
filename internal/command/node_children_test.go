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

func TestNodeAliases(t *testing.T) {
	t.Parallel()

	aliases := []string{"execute", "start"}

	node := MustNode(NodeSpec{
		Kind:    NodeCommand,
		ID:      MustID("bench.run"),
		Path:    MustPath("bench", "run"),
		Use:     "run",
		Aliases: aliases,
	})

	aliases[0] = "changed"

	if node.HasAlias("changed") {
		t.Fatalf("node alias changed through input slice mutation")
	}

	if !node.HasAlias("execute") {
		t.Fatalf("HasAlias(%q) = false, want true", "execute")
	}

	if !node.MatchesSegment("run") {
		t.Fatalf("MatchesSegment(%q) = false, want true", "run")
	}

	if !node.MatchesSegment("execute") {
		t.Fatalf("MatchesSegment(%q) = false, want true", "execute")
	}

	if RootPathNode := MustRootNode(); RootPathNode.MatchesSegment("") {
		t.Fatalf("root MatchesSegment(empty) = true, want false")
	}

	out := node.Aliases()
	out[0] = "changed"

	if node.HasAlias("changed") {
		t.Fatalf("node alias changed through output slice mutation")
	}
}

func TestNodeChildrenReturnsDetachedCopy(t *testing.T) {
	t.Parallel()

	family := mustTestFamilyNode(t, "bench", mustTestCommandNode(t, "bench.run"))

	children := family.Children()
	children[0] = mustTestCommandNode(t, "bench.compare")

	again := family.Children()
	if got, want := again[0].ID(), MustID("bench.run"); got != want {
		t.Fatalf("child mutated through returned slice: got %q, want %q", got, want)
	}
}

func TestNodeSpecReturnsDetachedCopies(t *testing.T) {
	t.Parallel()

	node := mustTestFamilyNode(t, "bench", mustTestCommandNode(t, "bench.run", "execute"))

	spec := node.Spec()
	spec.Aliases = []string{"changed"}
	spec.Children[0] = mustTestCommandNode(t, "bench.compare")

	again := node.Spec()
	if got, want := again.Children[0].ID(), MustID("bench.run"); got != want {
		t.Fatalf("Spec child mutated node: got %q, want %q", got, want)
	}
}

func TestFindChildByUse(t *testing.T) {
	t.Parallel()

	family := mustTestFamilyNode(
		t,
		"bench",
		mustTestCommandNode(t, "bench.run"),
		mustTestCommandNode(t, "bench.compare"),
	)

	child, ok := family.FindChildByUse("run")
	if !ok {
		t.Fatalf("FindChildByUse(%q) ok = false, want true", "run")
	}

	if got, want := child.ID(), MustID("bench.run"); got != want {
		t.Fatalf("FindChildByUse(%q) ID = %q, want %q", "run", got, want)
	}

	if _, ok := family.FindChildByUse("missing"); ok {
		t.Fatalf("FindChildByUse(%q) ok = true, want false", "missing")
	}
}

func TestFindChildBySegment(t *testing.T) {
	t.Parallel()

	family := mustTestFamilyNode(t, "bench", mustTestCommandNode(t, "bench.run", "execute"))

	child, ok := family.FindChildBySegment("execute")
	if !ok {
		t.Fatalf("FindChildBySegment(%q) ok = false, want true", "execute")
	}

	if got, want := child.ID(), MustID("bench.run"); got != want {
		t.Fatalf("FindChildBySegment(%q) ID = %q, want %q", "execute", got, want)
	}
}

func TestNodeAppendChild(t *testing.T) {
	t.Parallel()

	run := mustTestCommandNode(t, "bench.run")
	compare := mustTestCommandNode(t, "bench.compare")
	family := mustTestFamilyNode(t, "bench", run)

	next, err := family.AppendChild(compare)
	if err != nil {
		t.Fatalf("AppendChild() returned unexpected error: %v", err)
	}

	if got, want := family.ChildCount(), 1; got != want {
		t.Fatalf("AppendChild mutated original child count: got %d, want %d", got, want)
	}

	if got, want := next.ChildCount(), 2; got != want {
		t.Fatalf("AppendChild result child count = %d, want %d", got, want)
	}
}

func TestNodeWithChildren(t *testing.T) {
	t.Parallel()

	run := mustTestCommandNode(t, "bench.run")
	compare := mustTestCommandNode(t, "bench.compare")
	family := mustTestFamilyNode(t, "bench", run)

	next, err := family.WithChildren(compare)
	if err != nil {
		t.Fatalf("WithChildren() returned unexpected error: %v", err)
	}

	child, ok := next.FindChildByUse("compare")
	if !ok {
		t.Fatalf("WithChildren result does not contain compare child")
	}

	if got, want := child.ID(), MustID("bench.compare"); got != want {
		t.Fatalf("WithChildren child ID = %q, want %q", got, want)
	}

	if _, ok := family.FindChildByUse("compare"); ok {
		t.Fatalf("WithChildren mutated original family")
	}
}

func TestNodeWithChildrenRejectsInvalidResults(t *testing.T) {
	t.Parallel()

	t.Run("family empty", func(t *testing.T) {
		t.Parallel()

		family := mustTestFamilyNode(t, "bench", mustTestCommandNode(t, "bench.run"))

		_, err := family.WithChildren()
		if !errors.Is(err, ErrInvalidNode) {
			t.Fatalf("WithChildren() error = %v, want ErrInvalidNode", err)
		}
	})

	t.Run("command with children", func(t *testing.T) {
		t.Parallel()

		command := mustTestCommandNode(t, "bench.run")
		child := mustTestCommandNode(t, "bench.run.child")

		_, err := command.WithChildren(child)
		if !errors.Is(err, ErrInvalidNode) {
			t.Fatalf("WithChildren() error = %v, want ErrInvalidNode", err)
		}
	})
}

func TestNodeMustChildMutatorsPanicForInvalidResults(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = mustTestCommandNode(t, "bench.run").MustAppendChild(mustTestCommandNode(t, "bench.run.child"))
	})

	assertPanics(t, func() {
		_ = mustTestFamilyNode(t, "bench", mustTestCommandNode(t, "bench.run")).MustWithChildren()
	})
}
