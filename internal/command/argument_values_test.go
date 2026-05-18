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

// TestValidateArgumentAllowedValues verifies allowed-value validation.
func TestValidateArgumentAllowedValues(t *testing.T) {
	t.Parallel()

	if err := validateArgumentAllowedValues(OptionKindEnum, []string{"text", "json"}); err != nil {
		t.Fatalf("validateArgumentAllowedValues(valid) returned unexpected error: %v", err)
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

			err := validateArgumentAllowedValues(tt.kind, tt.values)
			if err == nil {
				t.Fatalf("validateArgumentAllowedValues() returned nil error")
			}

			if !errors.Is(err, ErrInvalidArgument) {
				t.Fatalf("validateArgumentAllowedValues() error = %v, want ErrInvalidArgument", err)
			}
		})
	}
}

// TestValidateArgumentMetavar verifies metavar validation.
func TestValidateArgumentMetavar(t *testing.T) {
	t.Parallel()

	if err := validateArgumentMetavar("PATH"); err != nil {
		t.Fatalf("validateArgumentMetavar(valid) returned unexpected error: %v", err)
	}

	invalid := []string{
		"",
		" PATH",
		"OUTPUT PATH",
		"bad\x00value",
		strings.Repeat("X", maxArgumentMetavarLength+1),
	}

	for _, raw := range invalid {
		raw := raw

		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			err := validateArgumentMetavar(raw)
			if err == nil {
				t.Fatalf("validateArgumentMetavar(%q) returned nil error", raw)
			}

			if !errors.Is(err, ErrInvalidArgument) {
				t.Fatalf("validateArgumentMetavar(%q) error = %v, want ErrInvalidArgument", raw, err)
			}
		})
	}

	if err := validateArgumentMetavar("OUTPUT PATH"); !errors.Is(err, textvalidate.ErrInvalidTokenText) {
		t.Fatalf("validateArgumentMetavar(token) error = %v, want ErrInvalidTokenText", err)
	}
}

// TestValidateArgumentRawValue verifies raw value validation.
func TestValidateArgumentRawValue(t *testing.T) {
	t.Parallel()

	if err := validateArgumentRawValue("value", "hello\nworld"); err != nil {
		t.Fatalf("validateArgumentRawValue(valid) returned unexpected error: %v", err)
	}

	invalidUTF8 := string([]byte{0xff, 0xfe})
	if err := validateArgumentRawValue("value", invalidUTF8); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("validateArgumentRawValue(invalid UTF-8) error = %v, want ErrInvalidArgument", err)
	}

	if err := validateArgumentRawValue("value", strings.Repeat("x", maxArgumentValueLength+1)); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("validateArgumentRawValue(too long) error = %v, want ErrInvalidArgument", err)
	}

	if err := validateArgumentRawValue("value", "bad\x00value"); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("validateArgumentRawValue(control) error = %v, want ErrInvalidArgument", err)
	}

	if err := validateArgumentRawValue("value", "bad\x00value"); !errors.Is(err, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("validateArgumentRawValue(control) error = %v, want ErrInvalidCompactText", err)
	}
}
