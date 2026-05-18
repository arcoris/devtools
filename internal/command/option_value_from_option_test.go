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

// TestOptionValueFromOptionAcceptsExplicitDefaultValues verifies explicit
// default-source values when they match the declaration defaults exactly.
func TestOptionValueFromOptionAcceptsExplicitDefaultValues(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name:          "format",
		Kind:          OptionKindEnum,
		AllowedValues: []string{"text", "json"},
		DefaultValues: []string{"text"},
	})

	value, err := NewOptionValueFromOption(option, OptionSourceDefault, "text")
	if err != nil {
		t.Fatalf("NewOptionValueFromOption() returned unexpected error: %v", err)
	}

	if !value.IsDefault() {
		t.Fatalf("IsDefault() = false, want true")
	}

	if got, want := value.MustValue(), "text"; got != want {
		t.Fatalf("MustValue() = %q, want %q", got, want)
	}
}

// TestOptionValueFromOptionAllowsEmptyStringWhenPolicyAllows verifies explicit empty string values.
func TestOptionValueFromOptionAllowsEmptyStringWhenPolicyAllows(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name: "output",
		Kind: OptionKindString,
		Policy: MustOptionPolicyForKind(OptionKindString, OptionPolicySpec{
			EmptyValue: OptionEmptyValueAllow,
		}),
	})

	value, err := NewOptionValueFromOption(option, OptionSourceCommandLine, "")
	if err != nil {
		t.Fatalf("NewOptionValueFromOption() returned unexpected error: %v", err)
	}

	if got, want := value.MustValue(), ""; got != want {
		t.Fatalf("MustValue() = %q, want empty", got)
	}
}

// TestOptionValueFromOptionRejectsDefaultMismatch verifies default provenance enforcement.
func TestOptionValueFromOptionRejectsDefaultMismatch(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name:          "format",
		Kind:          OptionKindEnum,
		AllowedValues: []string{"text", "json"},
		DefaultValues: []string{"text"},
	})

	_, err := NewOptionValueFromOption(option, OptionSourceDefault, "json")
	if err == nil {
		t.Fatalf("NewOptionValueFromOption() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("NewOptionValueFromOption() error = %v, want ErrInvalidOptionValue", err)
	}
}

// TestOptionValueFromOptionRejectsDefaultWithoutDeclaration verifies default
// source values cannot be invented when the option has no declared default.
func TestOptionValueFromOptionRejectsDefaultWithoutDeclaration(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name: "output",
		Kind: OptionKindString,
	})

	_, err := NewDefaultOptionValue(option)
	if err == nil {
		t.Fatalf("NewDefaultOptionValue() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("NewDefaultOptionValue() error = %v, want ErrInvalidOptionValue", err)
	}
}

// TestOptionValueFromOptionRejectsDisallowedValue verifies option allowed values.
func TestOptionValueFromOptionRejectsDisallowedValue(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name:          "format",
		Kind:          OptionKindEnum,
		AllowedValues: []string{"text", "json"},
	})

	_, err := NewOptionValueFromOption(option, OptionSourceCommandLine, "xml")
	if err == nil {
		t.Fatalf("NewOptionValueFromOption() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("NewOptionValueFromOption() error = %v, want ErrInvalidOptionValue", err)
	}
}

// TestOptionValueFromOptionRejectsEmptyValueByPolicy verifies empty-value policy.
func TestOptionValueFromOptionRejectsEmptyValueByPolicy(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name: "output",
		Kind: OptionKindString,
	})

	_, err := NewOptionValueFromOption(option, OptionSourceCommandLine, "")
	if err == nil {
		t.Fatalf("NewOptionValueFromOption() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("NewOptionValueFromOption() error = %v, want ErrInvalidOptionValue", err)
	}
}

// TestOptionValueFromOptionRejectsInvalidSource verifies policy source enforcement.
func TestOptionValueFromOptionRejectsInvalidSource(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name: "output",
		Kind: OptionKindString,
		Policy: MustOptionPolicy(OptionPolicySpec{
			AllowedSources: []OptionSource{OptionSourceCommandLine},
		}),
	})

	_, err := NewOptionValueFromOption(option, OptionSourceEnvironment, "out.txt")
	if err == nil {
		t.Fatalf("NewOptionValueFromOption() returned nil error")
	}

	if !errors.Is(err, ErrInvalidOptionValue) {
		t.Fatalf("NewOptionValueFromOption() error = %v, want ErrInvalidOptionValue", err)
	}
}

// TestOptionValueFromOptionUsesDefaultValues verifies default-source value construction.
func TestOptionValueFromOptionUsesDefaultValues(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{
		Name:          "format",
		Kind:          OptionKindEnum,
		AllowedValues: []string{"text", "json"},
		DefaultValues: []string{"text"},
	})

	value, err := NewDefaultOptionValue(option)
	if err != nil {
		t.Fatalf("NewDefaultOptionValue() returned unexpected error: %v", err)
	}

	if !value.IsDefault() {
		t.Fatalf("IsDefault() = false, want true")
	}

	if got, want := value.MustValue(), "text"; got != want {
		t.Fatalf("MustValue() = %q, want %q", got, want)
	}
}
