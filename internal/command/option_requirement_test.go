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

// TestOptionRequirementValidation verifies requirement parsing and predicates.
func TestOptionRequirementValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		raw      string
		want     OptionRequirement
		optional bool
		required bool
	}{
		{
			raw:      "optional",
			want:     OptionRequirementOptional,
			optional: true,
		},
		{
			raw:      "required",
			want:     OptionRequirementRequired,
			required: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.raw, func(t *testing.T) {
			t.Parallel()

			got, err := NewOptionRequirement(tt.raw)
			if err != nil {
				t.Fatalf("NewOptionRequirement(%q) returned unexpected error: %v", tt.raw, err)
			}

			if got != tt.want {
				t.Fatalf("NewOptionRequirement(%q) = %q, want %q", tt.raw, got, tt.want)
			}

			if got.IsOptional() != tt.optional {
				t.Fatalf("IsOptional() = %v, want %v", got.IsOptional(), tt.optional)
			}

			if got.IsRequired() != tt.required {
				t.Fatalf("IsRequired() = %v, want %v", got.IsRequired(), tt.required)
			}
		})
	}

	invalid := []string{"", "mandatory", "Optional"}
	for _, raw := range invalid {
		raw := raw

		t.Run("invalid-"+raw, func(t *testing.T) {
			t.Parallel()

			_, err := NewOptionRequirement(raw)
			if err == nil {
				t.Fatalf("NewOptionRequirement(%q) returned nil error", raw)
			}

			if !errors.Is(err, ErrInvalidOptionPolicy) {
				t.Fatalf("NewOptionRequirement(%q) error = %v, want ErrInvalidOptionPolicy", raw, err)
			}
		})
	}
}
