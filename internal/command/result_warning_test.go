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

// TestMustResultWarningPanicsForInvalidWarning verifies fail-fast warning construction.
func TestMustResultWarningPanicsForInvalidWarning(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustResultWarning did not panic")
		}
	}()

	_ = MustResultWarning(ResultWarningSpec{})
}

// TestNewResultWarningAcceptsValidWarning verifies warning construction.
func TestNewResultWarningAcceptsValidWarning(t *testing.T) {
	t.Parallel()

	warning, err := NewResultWarning(ResultWarningSpec{
		Kind:    "missing-tool",
		Message: "benchstat was not found.",
		Hints:   []string{"Install golang.org/x/perf/cmd/benchstat."},
	})
	if err != nil {
		t.Fatalf("NewResultWarning() returned unexpected error: %v", err)
	}

	if got, want := warning.Kind(), "missing-tool"; got != want {
		t.Fatalf("Kind() = %q, want %q", got, want)
	}

	if got, want := warning.Message(), "benchstat was not found."; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}

	if !warning.HasHints() {
		t.Fatalf("HasHints() = false, want true")
	}
}

// TestNewResultWarningNormalizesHints verifies hint normalization.
func TestNewResultWarningNormalizesHints(t *testing.T) {
	t.Parallel()

	warning := MustResultWarning(ResultWarningSpec{
		Kind:    "partial",
		Message: "Partial.",
		Hints: []string{
			"  First hint.  ",
			" ",
			"  Second hint.  ",
		},
	})

	resultTestAssertStrings(t, warning.Hints(), []string{"First hint.", "Second hint."})
}

// TestNewResultWarningRejectsInvalidWarning verifies warning validation.
func TestNewResultWarningRejectsInvalidWarning(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec ResultWarningSpec
	}{
		{
			name: "empty kind",
			spec: ResultWarningSpec{
				Message: "Message.",
			},
		},
		{
			name: "invalid kind",
			spec: ResultWarningSpec{
				Kind:    "Bad",
				Message: "Message.",
			},
		},
		{
			name: "empty message",
			spec: ResultWarningSpec{
				Kind: "partial",
			},
		},
		{
			name: "invalid message",
			spec: ResultWarningSpec{
				Kind:    "partial",
				Message: "bad\x00message",
			},
		},
		{
			name: "duplicate hints",
			spec: ResultWarningSpec{
				Kind:    "partial",
				Message: "Message.",
				Hints:   []string{"same", " same "},
			},
		},
		{
			name: "invalid hint",
			spec: ResultWarningSpec{
				Kind:    "partial",
				Message: "Message.",
				Hints:   []string{"bad\x00hint"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewResultWarning(tt.spec)
			if err == nil {
				t.Fatalf("NewResultWarning() returned nil error")
			}

			if !errors.Is(err, ErrInvalidResultWarning) {
				t.Fatalf("NewResultWarning() error = %v, want ErrInvalidResultWarning", err)
			}
		})
	}
}

// TestResultWarningCopySemantics verifies detached hint slices.
func TestResultWarningCopySemantics(t *testing.T) {
	t.Parallel()

	hints := []string{"hint"}

	warning := MustResultWarning(ResultWarningSpec{
		Kind:    "partial",
		Message: "Partial.",
		Hints:   hints,
	})

	hints[0] = "changed"

	if got, want := warning.Hints()[0], "hint"; got != want {
		t.Fatalf("hint changed through input slice: got %q, want %q", got, want)
	}

	out := warning.Hints()
	out[0] = "changed"

	if got, want := warning.Hints()[0], "hint"; got != want {
		t.Fatalf("hint changed through output slice: got %q, want %q", got, want)
	}
}

// TestResultWarningWithHint verifies immutable-style warning hint append.
func TestResultWarningWithHint(t *testing.T) {
	t.Parallel()

	warning := MustResultWarning(ResultWarningSpec{
		Kind:    "partial",
		Message: "Partial.",
	}).MustWithHint("Run the full suite.")

	if got, want := len(warning.Hints()), 1; got != want {
		t.Fatalf("len(Hints()) = %d, want %d", got, want)
	}
}
