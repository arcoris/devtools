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
	"strings"
	"testing"

	"arcoris.dev/devtools/internal/textvalidate"
)

func TestNewDeprecationAcceptsValidDeprecation(t *testing.T) {
	t.Parallel()

	deprecation, err := NewDeprecation(DeprecationSpec{
		Since:       "v0.2.0",
		Message:     "Use bench run instead.",
		Replacement: MustPath("bench", "run"),
	})
	if err != nil {
		t.Fatalf("NewDeprecation() returned unexpected error: %v", err)
	}

	if got, want := deprecation.Since(), "v0.2.0"; got != want {
		t.Fatalf("Since() = %q, want %q", got, want)
	}

	if got, want := deprecation.Message(), "Use bench run instead."; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}

	replacement, ok := deprecation.Replacement()
	if !ok {
		t.Fatalf("Replacement() ok = false, want true")
	}

	if got, want := replacement, MustPath("bench", "run"); !got.Equal(want) {
		t.Fatalf("Replacement() = %q, want %q", got, want)
	}
}

func TestNewDeprecationAcceptsNoReplacement(t *testing.T) {
	t.Parallel()

	deprecation := MustDeprecation(DeprecationSpec{Message: "This command is deprecated."})

	if deprecation.HasReplacement() {
		t.Fatalf("HasReplacement() = true, want false")
	}

	replacement, ok := deprecation.Replacement()
	if ok {
		t.Fatalf("Replacement() ok = true, want false")
	}

	if !replacement.IsRoot() {
		t.Fatalf("Replacement() path = %q, want root", replacement)
	}
}

func TestNewDeprecationRejectsInvalidDeprecation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec DeprecationSpec
	}{
		{name: "blank message", spec: DeprecationSpec{Message: "   "}},
		{name: "too long message", spec: DeprecationSpec{Message: strings.Repeat("x", maxDeprecationMessageLength+1)}},
		{name: "too long since", spec: DeprecationSpec{Since: strings.Repeat("x", maxMetadataTextLength+1), Message: "Deprecated."}},
		{name: "invalid replacement", spec: DeprecationSpec{Message: "Deprecated.", Replacement: Path{segments: []string{"Invalid"}}}},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewDeprecation(test.spec)
			if err == nil {
				t.Fatalf("NewDeprecation() returned nil error")
			}

			if !errors.Is(err, ErrInvalidDeprecation) {
				t.Fatalf("NewDeprecation() error = %v, want ErrInvalidDeprecation", err)
			}
		})
	}
}

func TestNewDeprecationWrapsTextvalidateErrorWithoutMetadataSentinel(t *testing.T) {
	t.Parallel()

	_, err := NewDeprecation(DeprecationSpec{
		Since:   "bad\x00value",
		Message: "Deprecated.",
	})
	if err == nil {
		t.Fatalf("NewDeprecation() returned nil error")
	}

	if !errors.Is(err, ErrInvalidDeprecation) {
		t.Fatalf("NewDeprecation() error = %v, want ErrInvalidDeprecation", err)
	}

	if !errors.Is(err, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("NewDeprecation() error = %v, want ErrInvalidCompactText", err)
	}

	if errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("NewDeprecation() error = %v, must not wrap ErrInvalidMetadata", err)
	}
}

func TestMustDeprecationPanicsForInvalidDeprecation(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustDeprecation(DeprecationSpec{})
	})
}
