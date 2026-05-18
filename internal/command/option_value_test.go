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

// TestMustOptionValuePanicsForInvalidValue verifies fail-fast construction.
func TestMustOptionValuePanicsForInvalidValue(t *testing.T) {
	t.Parallel()

	assertOptionValuePanics(t, func() {
		_ = MustOptionValue(OptionValueSpec{})
	})
}

// TestNewOptionValueAcceptsListValue verifies list resolved value construction.
func TestNewOptionValueAcceptsListValue(t *testing.T) {
	t.Parallel()

	value, err := NewOptionValue(OptionValueSpec{
		Name:   "package",
		Kind:   OptionKindStringList,
		Source: OptionSourceCommandLine,
		Values: []string{"./...", "./internal/..."},
	})
	if err != nil {
		t.Fatalf("NewOptionValue() returned unexpected error: %v", err)
	}

	if !value.IsList() {
		t.Fatalf("IsList() = false, want true")
	}

	if got, want := value.Len(), 2; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}

	if got, want := value.String(), "./...,./internal/..."; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

// TestNewOptionValueAcceptsScalarValue verifies scalar resolved value construction.
func TestNewOptionValueAcceptsScalarValue(t *testing.T) {
	t.Parallel()

	value, err := NewOptionValue(OptionValueSpec{
		Name:   "timeout",
		Kind:   OptionKindDuration,
		Source: OptionSourceCommandLine,
		Values: []string{"10s"},
	})
	if err != nil {
		t.Fatalf("NewOptionValue() returned unexpected error: %v", err)
	}

	if got, want := value.Name(), MustOptionName("timeout"); got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}

	if got, want := value.Kind(), OptionKindDuration; got != want {
		t.Fatalf("Kind() = %q, want %q", got, want)
	}

	if got, want := value.Source(), OptionSourceCommandLine; got != want {
		t.Fatalf("Source() = %q, want %q", got, want)
	}

	if !value.IsScalar() {
		t.Fatalf("IsScalar() = false, want true")
	}

	if got, want := value.String(), "10s"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

// TestNewOptionValueRejectsInvalidValue verifies generic value validation.
func TestNewOptionValueRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec OptionValueSpec
		err  error
	}{
		{
			name: "invalid name",
			spec: OptionValueSpec{
				Name:   "Timeout",
				Kind:   OptionKindDuration,
				Source: OptionSourceCommandLine,
				Values: []string{"10s"},
			},
			err: ErrInvalidOptionValue,
		},
		{
			name: "invalid kind",
			spec: OptionValueSpec{
				Name:   "timeout",
				Kind:   OptionKind("time"),
				Source: OptionSourceCommandLine,
				Values: []string{"10s"},
			},
			err: ErrInvalidOptionValue,
		},
		{
			name: "invalid source",
			spec: OptionValueSpec{
				Name:   "timeout",
				Kind:   OptionKindDuration,
				Source: OptionSource("cli"),
				Values: []string{"10s"},
			},
			err: ErrInvalidOptionValue,
		},
		{
			name: "no values",
			spec: OptionValueSpec{
				Name:   "timeout",
				Kind:   OptionKindDuration,
				Source: OptionSourceCommandLine,
			},
			err: ErrInvalidOptionValue,
		},
		{
			name: "multiple scalar values",
			spec: OptionValueSpec{
				Name:   "timeout",
				Kind:   OptionKindDuration,
				Source: OptionSourceCommandLine,
				Values: []string{"10s", "20s"},
			},
			err: ErrInvalidOptionValue,
		},
		{
			name: "invalid bool",
			spec: OptionValueSpec{
				Name:   "verbose",
				Kind:   OptionKindBool,
				Source: OptionSourceCommandLine,
				Values: []string{"yes"},
			},
			err: ErrInvalidOptionValue,
		},
		{
			name: "invalid int list element",
			spec: OptionValueSpec{
				Name:   "count",
				Kind:   OptionKindIntList,
				Source: OptionSourceCommandLine,
				Values: []string{"1", "bad"},
			},
			err: ErrInvalidOptionValue,
		},
		{
			name: "invalid enum value",
			spec: OptionValueSpec{
				Name:   "format",
				Kind:   OptionKindEnum,
				Source: OptionSourceCommandLine,
				Values: []string{"JSON"},
			},
			err: ErrInvalidOptionValue,
		},
		{
			name: "control rune",
			spec: OptionValueSpec{
				Name:   "output",
				Kind:   OptionKindString,
				Source: OptionSourceCommandLine,
				Values: []string{"bad\x00value"},
			},
			err: ErrInvalidOptionValue,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewOptionValue(tt.spec)
			if err == nil {
				t.Fatalf("NewOptionValue() returned nil error")
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewOptionValue() error = %v, want %v", err, tt.err)
			}
		})
	}
}

// TestScalarAndListOptionValueConstructorsEnforceShape verifies helper
// constructor contracts are stricter than the generic constructor.
func TestScalarAndListOptionValueConstructorsEnforceShape(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		build func() (OptionValue, error)
	}{
		{
			name: "scalar helper rejects list kind",
			build: func() (OptionValue, error) {
				return NewScalarOptionValue("package", OptionKindStringList, OptionSourceCommandLine, "./...")
			},
		},
		{
			name: "list helper rejects scalar kind",
			build: func() (OptionValue, error) {
				return NewListOptionValue("output", OptionKindString, OptionSourceCommandLine, "out.txt")
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.build()
			if err == nil {
				t.Fatalf("constructor returned nil error")
			}

			if !errors.Is(err, ErrInvalidOptionValue) {
				t.Fatalf("constructor error = %v, want ErrInvalidOptionValue", err)
			}
		})
	}
}
