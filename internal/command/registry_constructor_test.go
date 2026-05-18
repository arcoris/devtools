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

func TestNewRegistryFromChildrenBuildsIndexedRegistry(t *testing.T) {
	t.Parallel()

	registry, err := NewRegistryFromChildren(mustTestBenchFamilyNode(t))
	if err != nil {
		t.Fatalf("NewRegistryFromChildren() returned unexpected error: %v", err)
	}

	if registry.IsZero() {
		t.Fatalf("registry IsZero() = true, want false")
	}

	if registry.IsEmpty() {
		t.Fatalf("registry IsEmpty() = true, want false")
	}

	if got, want := registry.Size(), 4; got != want {
		t.Fatalf("Size() = %d, want %d", got, want)
	}

	if got, want := registry.FamilyCount(), 1; got != want {
		t.Fatalf("FamilyCount() = %d, want %d", got, want)
	}

	if got, want := registry.CommandCount(), 2; got != want {
		t.Fatalf("CommandCount() = %d, want %d", got, want)
	}
}

func TestEmptyRegistry(t *testing.T) {
	t.Parallel()

	registry := EmptyRegistry()

	if registry.IsZero() {
		t.Fatalf("EmptyRegistry().IsZero() = true, want false")
	}

	if !registry.IsEmpty() {
		t.Fatalf("EmptyRegistry().IsEmpty() = false, want true")
	}

	if got, want := registry.Size(), 1; got != want {
		t.Fatalf("Size() = %d, want %d", got, want)
	}

	if err := registry.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}
}

func TestNewRegistryRejectsInvalidTree(t *testing.T) {
	t.Parallel()

	var tree Tree

	_, err := NewRegistry(tree)
	if err == nil {
		t.Fatalf("NewRegistry() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("NewRegistry() error = %v, want ErrInvalidRegistry", err)
	}

	if !errors.Is(err, ErrInvalidTree) {
		t.Fatalf("NewRegistry() error = %v, want wrapped ErrInvalidTree", err)
	}
}

func TestNewRegistryFromRootBuildsIndexedRegistry(t *testing.T) {
	t.Parallel()

	root := MustRootNode(mustTestBenchFamilyNode(t))

	registry, err := NewRegistryFromRoot(root)
	if err != nil {
		t.Fatalf("NewRegistryFromRoot() returned unexpected error: %v", err)
	}

	if got, want := registry.Size(), 4; got != want {
		t.Fatalf("Size() = %d, want %d", got, want)
	}
}

func TestNewRegistryFromRootRejectsNonRoot(t *testing.T) {
	t.Parallel()

	_, err := NewRegistryFromRoot(mustTestCommandNode(t, "check"))
	if err == nil {
		t.Fatalf("NewRegistryFromRoot() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("NewRegistryFromRoot() error = %v, want ErrInvalidRegistry", err)
	}

	if !errors.Is(err, ErrInvalidTree) {
		t.Fatalf("NewRegistryFromRoot() error = %v, want wrapped ErrInvalidTree", err)
	}
}

func TestMustRegistryPanicsForInvalidTree(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustRegistry(Tree{})
	})
}

func TestZeroRegistryValidateRejects(t *testing.T) {
	t.Parallel()

	var registry Registry

	err := registry.Validate()
	if err == nil {
		t.Fatalf("zero registry Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("zero registry Validate() error = %v, want ErrInvalidRegistry", err)
	}
}
