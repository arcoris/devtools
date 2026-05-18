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

// TestMustOptionKindPanicsForInvalidKind verifies fail-fast invalid behavior.
func TestMustOptionKindPanicsForInvalidKind(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustOptionKind did not panic")
		}
	}()

	_ = MustOptionKind("path")
}

// TestMustOptionKindReturnsKnownKind verifies fail-fast static construction.
func TestMustOptionKindReturnsKnownKind(t *testing.T) {
	t.Parallel()

	kind := MustOptionKind("int64")

	if got, want := kind, OptionKindInt64; got != want {
		t.Fatalf("MustOptionKind() = %q, want %q", got, want)
	}
}

// TestNewOptionKindAcceptsKnownValues verifies every declared option kind.
func TestNewOptionKindAcceptsKnownValues(t *testing.T) {
	t.Parallel()

	for _, want := range KnownOptionKinds() {
		want := want

		t.Run(want.String(), func(t *testing.T) {
			t.Parallel()

			got, err := NewOptionKind(want.String())
			if err != nil {
				t.Fatalf("NewOptionKind(%q) returned unexpected error: %v", want, err)
			}

			if got != want {
				t.Fatalf("NewOptionKind(%q) = %q, want %q", want, got, want)
			}

			if !got.IsKnown() {
				t.Fatalf("%q.IsKnown() = false, want true", got)
			}

			if !got.IsValid() {
				t.Fatalf("%q.IsValid() = false, want true", got)
			}
		})
	}
}

// TestNewOptionKindRejectsInvalidValues verifies invalid kind handling.
func TestNewOptionKindRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		err  error
	}{
		{
			name: "empty",
			raw:  "",
			err:  ErrEmptyOptionKind,
		},
		{
			name: "uppercase",
			raw:  "String",
			err:  ErrInvalidOptionKind,
		},
		{
			name: "unknown",
			raw:  "path",
			err:  ErrInvalidOptionKind,
		},
		{
			name: "underscore",
			raw:  "string_list",
			err:  ErrInvalidOptionKind,
		},
		{
			name: "space",
			raw:  "string list",
			err:  ErrInvalidOptionKind,
		},
		{
			name: "plural bool",
			raw:  "bool-list",
			err:  ErrInvalidOptionKind,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewOptionKind(tt.raw)
			if err == nil {
				t.Fatalf("NewOptionKind(%q) returned nil error and kind %q", tt.raw, got)
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewOptionKind(%q) error = %v, want errors.Is(..., %v)", tt.raw, err, tt.err)
			}

			if OptionKind(tt.raw).IsValid() {
				t.Fatalf("OptionKind(%q).IsValid() = true, want false", tt.raw)
			}
		})
	}
}

// TestOptionKindString verifies canonical string rendering.
func TestOptionKindString(t *testing.T) {
	t.Parallel()

	if got, want := OptionKindDurationList.String(), "duration-list"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

// TestParseOptionKindIsAliasForNewOptionKind verifies parse behavior.
func TestParseOptionKindIsAliasForNewOptionKind(t *testing.T) {
	t.Parallel()

	fromNew, err := NewOptionKind("duration-list")
	if err != nil {
		t.Fatalf("NewOptionKind() returned unexpected error: %v", err)
	}

	fromParse, err := ParseOptionKind("duration-list")
	if err != nil {
		t.Fatalf("ParseOptionKind() returned unexpected error: %v", err)
	}

	if fromParse != fromNew {
		t.Fatalf("ParseOptionKind() = %q, want %q", fromParse, fromNew)
	}
}
