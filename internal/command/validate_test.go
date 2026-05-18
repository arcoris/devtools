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

	"arcoris.dev/devtools/internal/textvalidate"
)

// TestValidateCommandNameSegmentAcceptsValidSegments verifies that the command
// adapter accepts valid shared kebab-case segments.
func TestValidateCommandNameSegmentAcceptsValidSegments(t *testing.T) {
	t.Parallel()

	tests := []string{
		"a",
		"abc",
		"a1",
		"a-1",
		"bench",
		"run",
		"release-notes",
		"generate2",
		"config-validate",
	}

	for _, segment := range tests {
		segment := segment

		t.Run(segment, func(t *testing.T) {
			t.Parallel()

			if err := validateCommandNameSegment(segment); err != nil {
				t.Fatalf("validateCommandNameSegment(%q) returned unexpected error: %v", segment, err)
			}
		})
	}
}

// TestValidateCommandNameSegmentMapsSharedErrors verifies that command callers
// see command-specific sentinel errors while the underlying shared validation
// error remains available through errors.Is.
func TestValidateCommandNameSegmentMapsSharedErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		segment         string
		wantCommandErr  error
		wantSharedError error
	}{
		{
			name:            "empty",
			segment:         "",
			wantCommandErr:  ErrEmptyCommandNameSegment,
			wantSharedError: textvalidate.ErrEmptyKebabSegment,
		},
		{
			name:            "invalid",
			segment:         "Bench",
			wantCommandErr:  ErrInvalidCommandNameSegment,
			wantSharedError: textvalidate.ErrInvalidKebabSegment,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateCommandNameSegment(tt.segment)
			if err == nil {
				t.Fatalf("validateCommandNameSegment(%q) returned nil error", tt.segment)
			}

			if !errors.Is(err, tt.wantCommandErr) {
				t.Fatalf(
					"validateCommandNameSegment(%q) error = %v, want command sentinel %v",
					tt.segment,
					err,
					tt.wantCommandErr,
				)
			}

			if !errors.Is(err, tt.wantSharedError) {
				t.Fatalf(
					"validateCommandNameSegment(%q) error = %v, want shared sentinel %v",
					tt.segment,
					err,
					tt.wantSharedError,
				)
			}
		})
	}
}
