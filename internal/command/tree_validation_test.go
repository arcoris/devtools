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

func TestTreeValidateUniquePathsRejectsDuplicateRawNodes(t *testing.T) {
	t.Parallel()

	duplicate := Node{
		kind: NodeCommand,
		id:   MustID("bench.run"),
		path: MustPath("bench", "run"),
		use:  "run",
	}

	tree := Tree{
		root: Node{
			kind:     NodeRoot,
			path:     RootPath(),
			children: []Node{duplicate, duplicate},
		},
	}

	err := tree.validateUniquePaths()
	if err == nil {
		t.Fatalf("validateUniquePaths() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTree) {
		t.Fatalf("validateUniquePaths() error = %v, want ErrInvalidTree", err)
	}
}

func TestTreeValidateUniqueIDsRejectsDuplicateRawNodes(t *testing.T) {
	t.Parallel()

	first := Node{
		kind: NodeCommand,
		id:   MustID("bench.run"),
		path: MustPath("bench", "run"),
		use:  "run",
	}

	second := Node{
		kind: NodeCommand,
		id:   MustID("bench.run"),
		path: MustPath("bench", "execute"),
		use:  "execute",
	}

	tree := Tree{
		root: Node{
			kind:     NodeRoot,
			path:     RootPath(),
			children: []Node{first, second},
		},
	}

	err := tree.validateUniqueIDs()
	if err == nil {
		t.Fatalf("validateUniqueIDs() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTree) {
		t.Fatalf("validateUniqueIDs() error = %v, want ErrInvalidTree", err)
	}
}
