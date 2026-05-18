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

// TestOptionOccurrenceValidation verifies occurrence parsing and defaults.
func TestOptionOccurrenceValidation(t *testing.T) {
	t.Parallel()

	single, err := NewOptionOccurrence("single")
	if err != nil {
		t.Fatalf("NewOptionOccurrence(single) returned unexpected error: %v", err)
	}

	if !single.IsSingle() {
		t.Fatalf("single IsSingle() = false, want true")
	}

	multiple, err := NewOptionOccurrence("multiple")
	if err != nil {
		t.Fatalf("NewOptionOccurrence(multiple) returned unexpected error: %v", err)
	}

	if !multiple.IsMultiple() {
		t.Fatalf("multiple IsMultiple() = false, want true")
	}

	if got, want := OptionOccurrence("").OrDefaultForKind(OptionKindString), OptionOccurrenceSingle; got != want {
		t.Fatalf("scalar default occurrence = %q, want %q", got, want)
	}

	if got, want := OptionOccurrence("").OrDefaultForKind(OptionKindStringList), OptionOccurrenceMultiple; got != want {
		t.Fatalf("list default occurrence = %q, want %q", got, want)
	}

	invalid := []string{"", "many", "Multiple"}
	for _, raw := range invalid {
		raw := raw

		t.Run("invalid-"+raw, func(t *testing.T) {
			t.Parallel()

			_, err := NewOptionOccurrence(raw)
			if err == nil {
				t.Fatalf("NewOptionOccurrence(%q) returned nil error", raw)
			}

			if !errors.Is(err, ErrInvalidOptionPolicy) {
				t.Fatalf("NewOptionOccurrence(%q) error = %v, want ErrInvalidOptionPolicy", raw, err)
			}
		})
	}
}
