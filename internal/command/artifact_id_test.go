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

// TestArtifactIDValidation verifies ArtifactID value object behavior.
func TestArtifactIDValidation(t *testing.T) {
	t.Parallel()

	id, err := NewArtifactID("bench.report")
	if err != nil {
		t.Fatalf("NewArtifactID() returned unexpected error: %v", err)
	}

	if got, want := id.String(), "bench.report"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	invalid := []struct {
		raw string
		err error
	}{
		{raw: "", err: ErrEmptyArtifactID},
		{raw: "Bench.Report", err: ErrInvalidArtifactID},
		{raw: ".bench", err: ErrInvalidArtifactID},
		{raw: "bench.", err: ErrInvalidArtifactID},
		{raw: "bench..report", err: ErrInvalidArtifactID},
		{raw: strings.Repeat("x", maxArtifactIDLength+1), err: ErrInvalidArtifactID},
	}

	for _, tt := range invalid {
		tt := tt

		t.Run(tt.raw, func(t *testing.T) {
			t.Parallel()

			_, err := NewArtifactID(tt.raw)
			if err == nil {
				t.Fatalf("NewArtifactID(%q) returned nil error", tt.raw)
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewArtifactID(%q) error = %v, want %v", tt.raw, err, tt.err)
			}
		})
	}
}
