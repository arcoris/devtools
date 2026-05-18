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

func TestNewGroupAcceptsValidGroups(t *testing.T) {
	t.Parallel()

	tests := []string{
		"quality",
		"benchmark",
		"diagnostics",
		"diagnostics.perf",
		"config.schema",
		"release-notes",
		"a",
		"a1",
		"a-1",
		"a.b.c",
	}

	for _, raw := range tests {
		raw := raw

		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			group, err := NewGroup(raw)
			if err != nil {
				t.Fatalf("NewGroup(%q) returned unexpected error: %v", raw, err)
			}

			if got := group.String(); got != raw {
				t.Fatalf("String() = %q, want %q", got, raw)
			}

			if !group.IsValid() {
				t.Fatalf("IsValid() = false, want true")
			}
		})
	}
}

func TestParseGroupIsAliasForNewGroup(t *testing.T) {
	t.Parallel()

	const raw = "diagnostics.perf"

	fromNew, err := NewGroup(raw)
	if err != nil {
		t.Fatalf("NewGroup(%q) returned unexpected error: %v", raw, err)
	}

	fromParse, err := ParseGroup(raw)
	if err != nil {
		t.Fatalf("ParseGroup(%q) returned unexpected error: %v", raw, err)
	}

	if fromParse != fromNew {
		t.Fatalf("ParseGroup(%q) = %q, want %q", raw, fromParse, fromNew)
	}
}

func TestNewGroupPartsBuildsCanonicalGroup(t *testing.T) {
	t.Parallel()

	group, err := NewGroupParts("diagnostics", "perf")
	if err != nil {
		t.Fatalf("NewGroupParts() returned unexpected error: %v", err)
	}

	if got, want := group, MustGroup("diagnostics.perf"); got != want {
		t.Fatalf("NewGroupParts() = %q, want %q", got, want)
	}
}

func TestMustGroupReturnsValidGroup(t *testing.T) {
	t.Parallel()

	group := MustGroup("diagnostics.perf")

	if got, want := group.String(), "diagnostics.perf"; got != want {
		t.Fatalf("MustGroup returned %q, want %q", got, want)
	}
}

func TestMustGroupPanicsForInvalidGroup(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustGroup("Diagnostics.Perf")
	})
}

func TestMustGroupPartsPanicsForInvalidSegment(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustGroupParts("diagnostics", "Perf")
	})
}
