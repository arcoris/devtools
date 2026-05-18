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

// TestArtifactValidationHelpers verifies lower-level validators.
func TestArtifactValidationHelpers(t *testing.T) {
	t.Parallel()

	if err := validateArtifactLocation("bench/reports/report.md"); err != nil {
		t.Fatalf("validateArtifactLocation(valid) returned unexpected error: %v", err)
	}

	if err := validateArtifactMediaType("application/json"); err != nil {
		t.Fatalf("validateArtifactMediaType(valid) returned unexpected error: %v", err)
	}

	if err := validateArtifactText("field", "value", maxArtifactDescriptionLength); err != nil {
		t.Fatalf("validateArtifactText(valid) returned unexpected error: %v", err)
	}

	if err := validateArtifactBlock("field", "line one\nline two", maxArtifactDescriptionLength); err != nil {
		t.Fatalf("validateArtifactBlock(valid) returned unexpected error: %v", err)
	}

	if err := validateArtifactText("field", "bad\nline", maxArtifactDescriptionLength); !errors.Is(err, textvalidate.ErrInvalidSingleLineText) {
		t.Fatalf("validateArtifactText(line break) error = %v, want ErrInvalidSingleLineText", err)
	}

	if err := validateArtifactBlock("field", "bad\x00block", maxArtifactDescriptionLength); !errors.Is(err, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("validateArtifactBlock(control) error = %v, want ErrInvalidCompactText", err)
	}

	if !isArtifactLowerHex('a') || !isArtifactLowerHex('9') {
		t.Fatalf("isArtifactLowerHex(valid) = false, want true")
	}

	if isArtifactLowerHex('A') || isArtifactLowerHex('g') {
		t.Fatalf("isArtifactLowerHex(invalid) = true, want false")
	}
}

// TestNewArtifactRejectsInvalidArtifact verifies artifact validation.
func TestNewArtifactRejectsInvalidArtifact(t *testing.T) {
	t.Parallel()

	negativeSize := int64(-1)

	tests := []struct {
		name string
		spec ArtifactSpec
		err  error
	}{
		{
			name: "empty id",
			spec: ArtifactSpec{
				Kind:     ArtifactKindReport,
				Location: "report.md",
			},
			err: ErrEmptyArtifactID,
		},
		{
			name: "invalid id",
			spec: ArtifactSpec{
				ID:       "Report",
				Kind:     ArtifactKindReport,
				Location: "report.md",
			},
			err: ErrInvalidArtifactID,
		},
		{
			name: "empty kind",
			spec: ArtifactSpec{
				ID:       "report",
				Location: "report.md",
			},
			err: ErrEmptyArtifactKind,
		},
		{
			name: "invalid kind",
			spec: ArtifactSpec{
				ID:       "report",
				Kind:     ArtifactKind("Report"),
				Location: "report.md",
			},
			err: ErrInvalidArtifactKind,
		},
		{
			name: "empty location",
			spec: ArtifactSpec{
				ID:   "report",
				Kind: ArtifactKindReport,
			},
			err: ErrEmptyArtifactLocation,
		},
		{
			name: "invalid location",
			spec: ArtifactSpec{
				ID:       "report",
				Kind:     ArtifactKindReport,
				Location: "bad\x00location",
			},
			err: ErrInvalidArtifactLocation,
		},
		{
			name: "invalid format",
			spec: ArtifactSpec{
				ID:       "report",
				Kind:     ArtifactKindReport,
				Location: "report.md",
				Format:   ArtifactFormat("Markdown"),
			},
			err: ErrInvalidArtifactFormat,
		},
		{
			name: "invalid media type",
			spec: ArtifactSpec{
				ID:        "report",
				Kind:      ArtifactKindReport,
				Location:  "report.md",
				MediaType: "markdown",
			},
			err: ErrInvalidArtifact,
		},
		{
			name: "negative size",
			spec: ArtifactSpec{
				ID:        "report",
				Kind:      ArtifactKindReport,
				Location:  "report.md",
				SizeBytes: &negativeSize,
			},
			err: ErrInvalidArtifact,
		},
		{
			name: "invalid digest",
			spec: ArtifactSpec{
				ID:       "report",
				Kind:     ArtifactKindReport,
				Location: "report.md",
				Digest: &ArtifactDigestSpec{
					Algorithm: ArtifactDigestSHA256,
					Value:     "bad",
				},
			},
			err: ErrInvalidArtifact,
		},
		{
			name: "invalid label",
			spec: ArtifactSpec{
				ID:       "report",
				Kind:     ArtifactKindReport,
				Location: "report.md",
				Labels:   []string{"BadLabel"},
			},
			err: ErrInvalidArtifact,
		},
		{
			name: "duplicate label",
			spec: ArtifactSpec{
				ID:       "report",
				Kind:     ArtifactKindReport,
				Location: "report.md",
				Labels:   []string{"bench", "bench"},
			},
			err: ErrInvalidArtifact,
		},
		{
			name: "invalid metadata",
			spec: ArtifactSpec{
				ID:       "report",
				Kind:     ArtifactKindReport,
				Location: "report.md",
				Metadata: Metadata{
					owner: "BadOwner",
				},
			},
			err: ErrInvalidArtifact,
		},
		{
			name: "invalid visibility",
			spec: ArtifactSpec{
				ID:         "report",
				Kind:       ArtifactKindReport,
				Location:   "report.md",
				Visibility: Visibility("private"),
			},
			err: ErrInvalidArtifact,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewArtifact(tt.spec)
			if err == nil {
				t.Fatalf("NewArtifact() returned nil error")
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewArtifact() error = %v, want %v", err, tt.err)
			}
		})
	}
}
