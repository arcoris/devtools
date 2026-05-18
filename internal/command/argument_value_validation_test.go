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

// TestArgumentValidateValuesAllowsEmptyStringWhenPolicyAllows verifies explicit
// empty string support.
func TestArgumentValidateValuesAllowsEmptyStringWhenPolicyAllows(t *testing.T) {
	t.Parallel()

	argument := MustArgument(ArgumentSpec{
		Name:       "text",
		Kind:       OptionKindString,
		EmptyValue: OptionEmptyValueAllow,
	})

	if err := argument.ValidateValues(""); err != nil {
		t.Fatalf("ValidateValues(empty string) returned unexpected error: %v", err)
	}
}

// TestArgumentValidateValuesForVariadic verifies variadic runtime value
// validation.
func TestArgumentValidateValuesForVariadic(t *testing.T) {
	t.Parallel()

	argument := MustArgument(ArgumentSpec{
		Name:        "package",
		Kind:        OptionKindString,
		Cardinality: ArgumentCardinalityVariadic,
	})

	if err := argument.ValidateValues("./...", "./internal/..."); err != nil {
		t.Fatalf("ValidateValues(variadic) returned unexpected error: %v", err)
	}

	if err := argument.ValidateValues(); err == nil {
		t.Fatalf("required variadic ValidateValues() returned nil error")
	}
}

// TestArgumentValidateValues verifies runtime value validation.
func TestArgumentValidateValues(t *testing.T) {
	t.Parallel()

	argument := MustArgument(ArgumentSpec{
		Name:          "format",
		Kind:          OptionKindEnum,
		AllowedValues: []string{"text", "json"},
	})

	if err := argument.ValidateValues("json"); err != nil {
		t.Fatalf("ValidateValues(valid) returned unexpected error: %v", err)
	}

	tests := []struct {
		name   string
		values []string
	}{
		{
			name:   "missing required",
			values: nil,
		},
		{
			name:   "too many single",
			values: []string{"text", "json"},
		},
		{
			name:   "invalid enum grammar",
			values: []string{"JSON"},
		},
		{
			name:   "not allowed",
			values: []string{"xml"},
		},
		{
			name:   "empty",
			values: []string{""},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := argument.ValidateValues(tt.values...)
			if err == nil {
				t.Fatalf("ValidateValues() returned nil error")
			}

			if !errors.Is(err, ErrInvalidArgumentValue) {
				t.Fatalf("ValidateValues() error = %v, want ErrInvalidArgumentValue", err)
			}
		})
	}
}
