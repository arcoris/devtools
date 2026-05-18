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

// TestNormalizeOptionSources verifies precedence-order normalization.
func TestNormalizeOptionSources(t *testing.T) {
	t.Parallel()

	got := normalizeOptionSources([]OptionSource{
		OptionSourceCommandLine,
		OptionSourceDefault,
		OptionSourceDefault,
		OptionSourceEnvironment,
	})

	want := []OptionSource{
		OptionSourceDefault,
		OptionSourceEnvironment,
		OptionSourceCommandLine,
	}

	if len(got) != len(want) {
		t.Fatalf("len(normalizeOptionSources()) = %d, want %d", len(got), len(want))
	}

	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("normalizeOptionSources()[%d] = %q, want %q", index, got[index], want[index])
		}
	}

	if normalizeOptionSources(nil) != nil {
		t.Fatalf("normalizeOptionSources(nil) must return nil")
	}
}

// TestValidateAllowedOptionSources verifies direct source-list validation.
func TestValidateAllowedOptionSources(t *testing.T) {
	t.Parallel()

	if err := validateAllowedOptionSources([]OptionSource{OptionSourceDefault}); err != nil {
		t.Fatalf("validateAllowedOptionSources(valid) returned unexpected error: %v", err)
	}

	tests := []struct {
		name    string
		sources []OptionSource
	}{
		{
			name:    "empty",
			sources: nil,
		},
		{
			name:    "invalid",
			sources: []OptionSource{OptionSource("file")},
		},
		{
			name:    "duplicate",
			sources: []OptionSource{OptionSourceDefault, OptionSourceDefault},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateAllowedOptionSources(tt.sources)
			if err == nil {
				t.Fatalf("validateAllowedOptionSources() returned nil error")
			}

			if !errors.Is(err, ErrInvalidOptionPolicy) {
				t.Fatalf("validateAllowedOptionSources() error = %v, want ErrInvalidOptionPolicy", err)
			}
		})
	}
}
