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
)

// TestIsZeroOptionPolicy verifies zero-policy detection.
func TestIsZeroOptionPolicy(t *testing.T) {
	t.Parallel()

	if !isZeroOptionPolicy(OptionPolicy{}) {
		t.Fatalf("zero OptionPolicy was not detected as zero")
	}

	if isZeroOptionPolicy(DefaultOptionPolicy()) {
		t.Fatalf("DefaultOptionPolicy() was detected as zero")
	}
}

// TestValidateOptionAllowedValues verifies allowed-value validation.
func TestValidateOptionAllowedValues(t *testing.T) {
	t.Parallel()

	if err := validateOptionAllowedValues(OptionKindEnum, []string{"text", "json"}); err != nil {
		t.Fatalf("validateOptionAllowedValues(valid) returned unexpected error: %v", err)
	}

	tests := []struct {
		name   string
		kind   OptionKind
		values []string
	}{
		{
			name:   "enum missing values",
			kind:   OptionKindEnum,
			values: nil,
		},
		{
			name:   "invalid value",
			kind:   OptionKindEnum,
			values: []string{"JSON"},
		},
		{
			name:   "duplicate value",
			kind:   OptionKindEnum,
			values: []string{"json", "json"},
		},
		{
			name:   "numeric values",
			kind:   OptionKindInt,
			values: []string{"one"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateOptionAllowedValues(tt.kind, tt.values)
			if err == nil {
				t.Fatalf("validateOptionAllowedValues() returned nil error")
			}

			if !errors.Is(err, ErrInvalidOption) {
				t.Fatalf("validateOptionAllowedValues() error = %v, want ErrInvalidOption", err)
			}
		})
	}
}

// TestValidateOptionMetavar verifies metavar validation.
func TestValidateOptionMetavar(t *testing.T) {
	t.Parallel()

	if err := validateOptionMetavar("PATH"); err != nil {
		t.Fatalf("validateOptionMetavar(valid) returned unexpected error: %v", err)
	}

	invalid := []string{
		"",
		" PATH",
		"OUTPUT PATH",
		"bad\x00value",
		strings.Repeat("X", maxOptionMetavarLength+1),
	}

	for _, raw := range invalid {
		raw := raw

		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			err := validateOptionMetavar(raw)
			if err == nil {
				t.Fatalf("validateOptionMetavar(%q) returned nil error", raw)
			}

			if !errors.Is(err, ErrInvalidOption) {
				t.Fatalf("validateOptionMetavar(%q) error = %v, want ErrInvalidOption", raw, err)
			}
		})
	}
}

// TestValidateOptionRawValue verifies raw value validation.
func TestValidateOptionRawValue(t *testing.T) {
	t.Parallel()

	if err := validateOptionRawValue("value", "hello\nworld"); err != nil {
		t.Fatalf("validateOptionRawValue(valid) returned unexpected error: %v", err)
	}

	invalidUTF8 := string([]byte{0xff, 0xfe})
	if err := validateOptionRawValue("value", invalidUTF8); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("validateOptionRawValue(invalid UTF-8) error = %v, want ErrInvalidOption", err)
	}

	if err := validateOptionRawValue("value", strings.Repeat("x", maxOptionValueLength+1)); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("validateOptionRawValue(too long) error = %v, want ErrInvalidOption", err)
	}

	if err := validateOptionRawValue("value", "bad\x00value"); !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("validateOptionRawValue(control) error = %v, want ErrInvalidOption", err)
	}
}

// TestValidateOptionShorthand verifies shorthand validation.
func TestValidateOptionShorthand(t *testing.T) {
	t.Parallel()

	valid := []string{"a", "Z", "1"}
	for _, raw := range valid {
		raw := raw

		t.Run("valid-"+raw, func(t *testing.T) {
			t.Parallel()

			if err := validateOptionShorthand(raw); err != nil {
				t.Fatalf("validateOptionShorthand(%q) returned unexpected error: %v", raw, err)
			}
		})
	}

	invalid := []string{"ab", "-", "_", "ф"}
	for _, raw := range invalid {
		raw := raw

		t.Run("invalid-"+raw, func(t *testing.T) {
			t.Parallel()

			err := validateOptionShorthand(raw)
			if err == nil {
				t.Fatalf("validateOptionShorthand(%q) returned nil error", raw)
			}

			if !errors.Is(err, ErrInvalidOption) {
				t.Fatalf("validateOptionShorthand(%q) error = %v, want ErrInvalidOption", raw, err)
			}
		})
	}
}

// TestValidateOptionValueForKind verifies scalar default parsing.
func TestValidateOptionValueForKind(t *testing.T) {
	t.Parallel()

	valid := []struct {
		kind OptionKind
		raw  string
	}{
		{kind: OptionKindBool, raw: "true"},
		{kind: OptionKindString, raw: "anything"},
		{kind: OptionKindEnum, raw: "json"},
		{kind: OptionKindInt, raw: "-10"},
		{kind: OptionKindInt64, raw: "-9223372036854775808"},
		{kind: OptionKindUint, raw: "10"},
		{kind: OptionKindUint64, raw: "18446744073709551615"},
		{kind: OptionKindFloat64, raw: "3.14"},
		{kind: OptionKindDuration, raw: "10s"},
	}

	for _, tt := range valid {
		tt := tt

		t.Run("valid-"+tt.kind.String(), func(t *testing.T) {
			t.Parallel()

			if err := validateOptionValueForKind(tt.kind, tt.raw); err != nil {
				t.Fatalf("validateOptionValueForKind(%q, %q) returned unexpected error: %v", tt.kind, tt.raw, err)
			}
		})
	}

	invalid := []struct {
		kind OptionKind
		raw  string
	}{
		{kind: OptionKindBool, raw: "yes"},
		{kind: OptionKindInt, raw: "abc"},
		{kind: OptionKindUint, raw: "-1"},
		{kind: OptionKindFloat64, raw: "NaN"},
		{kind: OptionKindFloat64, raw: "+Inf"},
		{kind: OptionKindDuration, raw: "soon"},
		{kind: OptionKindStringList, raw: "value"},
	}

	for _, tt := range invalid {
		tt := tt

		t.Run("invalid-"+tt.kind.String()+"-"+tt.raw, func(t *testing.T) {
			t.Parallel()

			if err := validateOptionValueForKind(tt.kind, tt.raw); err == nil {
				t.Fatalf("validateOptionValueForKind(%q, %q) returned nil error", tt.kind, tt.raw)
			}
		})
	}
}
