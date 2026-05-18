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

// TestOptionEmptyValuePolicyValidation verifies empty-value policy parsing.
func TestOptionEmptyValuePolicyValidation(t *testing.T) {
	t.Parallel()

	reject, err := NewOptionEmptyValuePolicy("reject-empty")
	if err != nil {
		t.Fatalf("NewOptionEmptyValuePolicy(reject-empty) returned unexpected error: %v", err)
	}

	if !reject.RejectsEmpty() {
		t.Fatalf("reject-empty RejectsEmpty() = false, want true")
	}

	allow, err := NewOptionEmptyValuePolicy("allow-empty")
	if err != nil {
		t.Fatalf("NewOptionEmptyValuePolicy(allow-empty) returned unexpected error: %v", err)
	}

	if !allow.AllowsEmpty() {
		t.Fatalf("allow-empty AllowsEmpty() = false, want true")
	}

	if got, want := OptionEmptyValuePolicy("").OrDefault(), OptionEmptyValueReject; got != want {
		t.Fatalf("zero OrDefault() = %q, want %q", got, want)
	}

	invalid := []string{"", "allow", "AllowEmpty"}
	for _, raw := range invalid {
		raw := raw

		t.Run("invalid-"+raw, func(t *testing.T) {
			t.Parallel()

			_, err := NewOptionEmptyValuePolicy(raw)
			if err == nil {
				t.Fatalf("NewOptionEmptyValuePolicy(%q) returned nil error", raw)
			}

			if !errors.Is(err, ErrInvalidOptionPolicy) {
				t.Fatalf("NewOptionEmptyValuePolicy(%q) error = %v, want ErrInvalidOptionPolicy", raw, err)
			}
		})
	}
}
