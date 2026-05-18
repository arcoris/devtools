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

func TestRegistryValidateWrapsTreeError(t *testing.T) {
	t.Parallel()

	registry := Registry{
		tree:   Tree{root: mustTestCommandNode(t, "check")},
		byID:   map[ID]Node{},
		byPath: map[string]Node{},
	}

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}

	if !errors.Is(err, ErrInvalidTree) {
		t.Fatalf("Validate() error = %v, want wrapped ErrInvalidTree", err)
	}
}

func TestRegistryValidateRejectsNilIDIndex(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)
	registry.byID = nil

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryValidateRejectsNilPathIndex(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)
	registry.byPath = nil

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryValidateRejectsTamperedIDIndex(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)
	delete(registry.byID, MustID("bench.run"))

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryValidateRejectsTamperedPathIndex(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)
	delete(registry.byPath, MustPath("bench", "run").Key())

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryValidateRejectsMismatchedIDIndexEntry(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)
	registry.byID[MustID("bench.run")] = mustTestCommandNode(t, "bench.compare")

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryValidateRejectsMismatchedPathIndexEntry(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)
	registry.byPath[MustPath("bench", "run").Key()] = mustTestCommandNode(t, "bench.compare")

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryValidateRejectsTamperedOrderedIDs(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)
	registry.ids = registry.ids[:len(registry.ids)-1]

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryValidateRejectsTamperedOrderedIDValue(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)
	registry.ids[1] = MustID("bench.run")

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryValidateRejectsTamperedOrderedPaths(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)
	registry.paths = registry.paths[:len(registry.paths)-1]

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryValidateRejectsTamperedOrderedPathValue(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)
	registry.paths[1] = MustPath("bench", "run")

	err := registry.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("Validate() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryRegisterPathRejectsDuplicate(t *testing.T) {
	t.Parallel()

	node := mustTestCommandNode(t, "bench.run")
	registry := Registry{byPath: make(map[string]Node)}

	if err := registry.registerPath(node.Path(), node); err != nil {
		t.Fatalf("registerPath() returned unexpected error: %v", err)
	}

	err := registry.registerPath(node.Path(), node)
	if err == nil {
		t.Fatalf("registerPath() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("registerPath() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryRegisterIDRejectsDuplicate(t *testing.T) {
	t.Parallel()

	node := mustTestCommandNode(t, "bench.run")
	registry := Registry{byID: make(map[ID]Node)}

	if err := registry.registerID(node.ID(), node); err != nil {
		t.Fatalf("registerID() returned unexpected error: %v", err)
	}

	err := registry.registerID(node.ID(), node)
	if err == nil {
		t.Fatalf("registerID() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("registerID() error = %v, want ErrInvalidRegistry", err)
	}
}

func TestRegistryRebuildIndexesRejectsDuplicateRawTree(t *testing.T) {
	t.Parallel()

	duplicate := Node{
		kind: NodeCommand,
		id:   MustID("bench.run"),
		path: MustPath("bench", "run"),
		use:  "run",
	}

	registry := Registry{
		tree: Tree{
			root: Node{
				kind:     NodeRoot,
				path:     RootPath(),
				children: []Node{duplicate, duplicate},
			},
		},
	}

	err := registry.rebuildIndexes()
	if err == nil {
		t.Fatalf("rebuildIndexes() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRegistry) {
		t.Fatalf("rebuildIndexes() error = %v, want ErrInvalidRegistry", err)
	}
}
