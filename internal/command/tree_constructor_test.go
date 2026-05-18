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

func TestNewTreeFromChildrenBuildsValidTree(t *testing.T) {
	t.Parallel()

	tree, err := NewTreeFromChildren(mustTestBenchFamilyNode(t))
	if err != nil {
		t.Fatalf("NewTreeFromChildren() returned unexpected error: %v", err)
	}

	if tree.IsZero() {
		t.Fatalf("tree IsZero() = true, want false")
	}

	if tree.IsEmpty() {
		t.Fatalf("tree IsEmpty() = true, want false")
	}

	if got, want := tree.Size(), 4; got != want {
		t.Fatalf("Size() = %d, want %d", got, want)
	}

	if got, want := tree.FamilyCount(), 1; got != want {
		t.Fatalf("FamilyCount() = %d, want %d", got, want)
	}

	if got, want := tree.CommandCount(), 2; got != want {
		t.Fatalf("CommandCount() = %d, want %d", got, want)
	}
}

func TestNewTreeRejectsNonRoot(t *testing.T) {
	t.Parallel()

	_, err := NewTree(mustTestCommandNode(t, "check"))
	if err == nil {
		t.Fatalf("NewTree() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTree) {
		t.Fatalf("NewTree() error = %v, want ErrInvalidTree", err)
	}
}

func TestTreeValidateRejectsZeroTree(t *testing.T) {
	t.Parallel()

	var tree Tree

	err := tree.Validate()
	if err == nil {
		t.Fatalf("zero tree Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTree) {
		t.Fatalf("zero tree Validate() error = %v, want ErrInvalidTree", err)
	}
}

func TestMustTreePanicsForInvalidRoot(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustTree(mustTestCommandNode(t, "check"))
	})
}

func TestEmptyTreeIsValid(t *testing.T) {
	t.Parallel()

	tree, err := NewTreeFromChildren()
	if err != nil {
		t.Fatalf("NewTreeFromChildren() returned unexpected error: %v", err)
	}

	if tree.IsZero() {
		t.Fatalf("empty tree IsZero() = true, want false")
	}

	if !tree.IsEmpty() {
		t.Fatalf("empty tree IsEmpty() = false, want true")
	}

	if got, want := tree.Size(), 1; got != want {
		t.Fatalf("Size() = %d, want %d", got, want)
	}

	if got, want := tree.CommandCount(), 0; got != want {
		t.Fatalf("CommandCount() = %d, want %d", got, want)
	}
}

func TestTreeZeroValueIsEmptyButInvalid(t *testing.T) {
	t.Parallel()

	var tree Tree

	if !tree.IsZero() {
		t.Fatalf("zero-value tree IsZero() = false, want true")
	}

	if !tree.IsEmpty() {
		t.Fatalf("zero-value tree IsEmpty() = false, want true")
	}
}

func TestTreeRejectsInvalidChildNode(t *testing.T) {
	t.Parallel()

	invalid := Node{
		kind: NodeCommand,
		id:   MustID("bench.execute"),
		path: MustPath("bench", "run"),
		use:  "run",
	}

	_, err := NewTreeFromChildren(invalid)
	if err == nil {
		t.Fatalf("NewTreeFromChildren() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTree) {
		t.Fatalf("NewTreeFromChildren() error = %v, want ErrInvalidTree", err)
	}

	if !errors.Is(err, ErrInvalidNode) {
		t.Fatalf("NewTreeFromChildren() error = %v, want wrapped ErrInvalidNode", err)
	}
}
