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

func TestRegistryWalk(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	var got []string
	err := registry.Walk(func(path Path, _ Node) error {
		got = append(got, path.String())
		return nil
	})
	if err != nil {
		t.Fatalf("Walk() returned unexpected error: %v", err)
	}

	assertStringSlicesEqual(t, got, []string{"", "bench", "bench run", "bench compare"})
}

func TestRegistryWalkOrder(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	var got []string
	err := registry.WalkOrder(WalkPostOrder, func(path Path, _ Node) error {
		got = append(got, path.String())
		return nil
	})
	if err != nil {
		t.Fatalf("WalkOrder() returned unexpected error: %v", err)
	}

	assertStringSlicesEqual(t, got, []string{"bench run", "bench compare", "bench", ""})
}

func TestRegistryWalkNilFuncIsNoop(t *testing.T) {
	t.Parallel()

	registry := mustTestRegistry(t)

	if err := registry.Walk(nil); err != nil {
		t.Fatalf("Walk(nil) error = %v, want nil", err)
	}
}
