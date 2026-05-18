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

// TestOptionScopeValidation verifies scope parsing and predicates.
func TestOptionScopeValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		raw       string
		want      OptionScope
		local     bool
		subtree   bool
		global    bool
		inherited bool
	}{
		{
			raw:   "local",
			want:  OptionScopeLocal,
			local: true,
		},
		{
			raw:       "subtree",
			want:      OptionScopeSubtree,
			subtree:   true,
			inherited: true,
		},
		{
			raw:       "global",
			want:      OptionScopeGlobal,
			global:    true,
			inherited: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.raw, func(t *testing.T) {
			t.Parallel()

			got, err := NewOptionScope(tt.raw)
			if err != nil {
				t.Fatalf("NewOptionScope(%q) returned unexpected error: %v", tt.raw, err)
			}

			if got != tt.want {
				t.Fatalf("NewOptionScope(%q) = %q, want %q", tt.raw, got, tt.want)
			}

			if got.IsLocal() != tt.local {
				t.Fatalf("IsLocal() = %v, want %v", got.IsLocal(), tt.local)
			}

			if got.IsSubtree() != tt.subtree {
				t.Fatalf("IsSubtree() = %v, want %v", got.IsSubtree(), tt.subtree)
			}

			if got.IsGlobal() != tt.global {
				t.Fatalf("IsGlobal() = %v, want %v", got.IsGlobal(), tt.global)
			}

			if got.IsInheritedByChildren() != tt.inherited {
				t.Fatalf("IsInheritedByChildren() = %v, want %v", got.IsInheritedByChildren(), tt.inherited)
			}
		})
	}

	invalid := []string{"", "tree", "Global"}
	for _, raw := range invalid {
		raw := raw

		t.Run("invalid-"+raw, func(t *testing.T) {
			t.Parallel()

			_, err := NewOptionScope(raw)
			if err == nil {
				t.Fatalf("NewOptionScope(%q) returned nil error", raw)
			}

			if !errors.Is(err, ErrInvalidOptionPolicy) {
				t.Fatalf("NewOptionScope(%q) error = %v, want ErrInvalidOptionPolicy", raw, err)
			}
		})
	}
}
