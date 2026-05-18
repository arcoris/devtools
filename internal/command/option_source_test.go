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

// TestMustOptionSourcePanicsForInvalidSource verifies fail-fast invalid behavior.
func TestMustOptionSourcePanicsForInvalidSource(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustOptionSource did not panic")
		}
	}()

	_ = MustOptionSource("cli")
}

// TestMustOptionSourceReturnsKnownSource verifies fail-fast static construction.
func TestMustOptionSourceReturnsKnownSource(t *testing.T) {
	t.Parallel()

	source := MustOptionSource("environment")

	if got, want := source, OptionSourceEnvironment; got != want {
		t.Fatalf("MustOptionSource() = %q, want %q", got, want)
	}
}

// TestNewOptionSourceAcceptsKnownValues verifies every declared option source.
func TestNewOptionSourceAcceptsKnownValues(t *testing.T) {
	t.Parallel()

	for _, want := range KnownOptionSources() {
		want := want

		t.Run(want.String(), func(t *testing.T) {
			t.Parallel()

			got, err := NewOptionSource(want.String())
			if err != nil {
				t.Fatalf("NewOptionSource(%q) returned unexpected error: %v", want, err)
			}

			if got != want {
				t.Fatalf("NewOptionSource(%q) = %q, want %q", want, got, want)
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

// TestNewOptionSourceRejectsInvalidValues verifies invalid source handling.
func TestNewOptionSourceRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		err  error
	}{
		{
			name: "empty",
			raw:  "",
			err:  ErrEmptyOptionSource,
		},
		{
			name: "uppercase",
			raw:  "Default",
			err:  ErrInvalidOptionSource,
		},
		{
			name: "unknown",
			raw:  "file",
			err:  ErrInvalidOptionSource,
		},
		{
			name: "underscore",
			raw:  "command_line",
			err:  ErrInvalidOptionSource,
		},
		{
			name: "space",
			raw:  "command line",
			err:  ErrInvalidOptionSource,
		},
		{
			name: "short cli",
			raw:  "cli",
			err:  ErrInvalidOptionSource,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewOptionSource(tt.raw)
			if err == nil {
				t.Fatalf("NewOptionSource(%q) returned nil error and source %q", tt.raw, got)
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewOptionSource(%q) error = %v, want errors.Is(..., %v)", tt.raw, err, tt.err)
			}

			if OptionSource(tt.raw).IsValid() {
				t.Fatalf("OptionSource(%q).IsValid() = true, want false", tt.raw)
			}
		})
	}
}

// TestParseOptionSourceIsAliasForNewOptionSource verifies parse behavior.
func TestParseOptionSourceIsAliasForNewOptionSource(t *testing.T) {
	t.Parallel()

	fromNew, err := NewOptionSource("command-line")
	if err != nil {
		t.Fatalf("NewOptionSource() returned unexpected error: %v", err)
	}

	fromParse, err := ParseOptionSource("command-line")
	if err != nil {
		t.Fatalf("ParseOptionSource() returned unexpected error: %v", err)
	}

	if fromParse != fromNew {
		t.Fatalf("ParseOptionSource() = %q, want %q", fromParse, fromNew)
	}
}
