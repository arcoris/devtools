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

func TestNewVisibilityAcceptsKnownValues(t *testing.T) {
	t.Parallel()

	for _, want := range KnownVisibilities() {
		want := want

		t.Run(want.String(), func(t *testing.T) {
			t.Parallel()

			got, err := NewVisibility(want.String())
			if err != nil {
				t.Fatalf("NewVisibility(%q) returned unexpected error: %v", want, err)
			}

			if got != want {
				t.Fatalf("NewVisibility(%q) = %q, want %q", want, got, want)
			}

			if !got.IsValid() {
				t.Fatalf("%q.IsValid() = false, want true", got)
			}

			if !got.IsKnown() {
				t.Fatalf("%q.IsKnown() = false, want true", got)
			}
		})
	}
}

func TestParseVisibilityIsAliasForNewVisibility(t *testing.T) {
	t.Parallel()

	fromNew, err := NewVisibility("hidden")
	if err != nil {
		t.Fatalf("NewVisibility() returned unexpected error: %v", err)
	}

	fromParse, err := ParseVisibility("hidden")
	if err != nil {
		t.Fatalf("ParseVisibility() returned unexpected error: %v", err)
	}

	if fromParse != fromNew {
		t.Fatalf("ParseVisibility() = %q, want %q", fromParse, fromNew)
	}
}

func TestMustVisibilityReturnsKnownValue(t *testing.T) {
	t.Parallel()

	visibility := MustVisibility("internal")

	if got, want := visibility, VisibilityInternal; got != want {
		t.Fatalf("MustVisibility() = %q, want %q", got, want)
	}
}

func TestMustVisibilityPanicsForInvalidValue(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustVisibility("private")
	})
}

func TestDefaultVisibility(t *testing.T) {
	t.Parallel()

	if got, want := DefaultVisibility(), VisibilityPublic; got != want {
		t.Fatalf("DefaultVisibility() = %q, want %q", got, want)
	}
}

func TestKnownVisibilitiesReturnsDetachedCopy(t *testing.T) {
	t.Parallel()

	values := KnownVisibilities()
	values[0] = VisibilityInternal

	again := KnownVisibilities()
	assertVisibilitySlicesEqual(t, again, []Visibility{VisibilityPublic, VisibilityHidden, VisibilityInternal})
}

func TestVisibilityFromHidden(t *testing.T) {
	t.Parallel()

	if got, want := VisibilityFromHidden(false), VisibilityPublic; got != want {
		t.Fatalf("VisibilityFromHidden(false) = %q, want %q", got, want)
	}

	if got, want := VisibilityFromHidden(true), VisibilityHidden; got != want {
		t.Fatalf("VisibilityFromHidden(true) = %q, want %q", got, want)
	}
}
