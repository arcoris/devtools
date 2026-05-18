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

// TestNewDocumentationAcceptsValidDocumentation verifies full documentation
// construction.
func TestNewDocumentationAcceptsValidDocumentation(t *testing.T) {
	t.Parallel()

	documentation, err := NewDocumentation(DocumentationSpec{
		Summary:     "Run configured benchmarks.",
		Description: "Runs configured benchmark suites.\nWrites raw output for later comparison.",
		Usage:       MustSimpleUsage("bench run [flags]"),
		Notes: []string{
			"Use stable suites for reproducible CI checks.",
			"Use smoke suites for fast local validation.",
		},
		References: []DocumentationReference{
			MustDocumentationReference(DocumentationReferenceSpec{
				Key:    "bench-compare",
				Kind:   DocumentationReferenceCommandPath,
				Label:  "Compare benchmark outputs",
				Target: "bench compare",
			}),
		},
	})
	if err != nil {
		t.Fatalf("NewDocumentation() returned unexpected error: %v", err)
	}

	if got, want := documentation.Summary(), "Run configured benchmarks."; got != want {
		t.Fatalf("Summary() = %q, want %q", got, want)
	}

	if got, want := documentation.Description(), "Runs configured benchmark suites.\nWrites raw output for later comparison."; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}

	usage, ok := documentation.Usage()
	if !ok || usage.String() != "bench run [flags]" {
		t.Fatalf("Usage() = %q, %t; want %q, true", usage, ok, "bench run [flags]")
	}

	if got, want := documentation.NoteCount(), 2; got != want {
		t.Fatalf("NoteCount() = %d, want %d", got, want)
	}

	reference, ok := documentation.Reference("bench-compare")
	if !ok || reference.Target() != "bench compare" {
		t.Fatalf("Reference(bench-compare) = %q, %t; want target %q, true", reference, ok, "bench compare")
	}
}

// TestEmptyDocumentation verifies zero documentation behavior.
func TestEmptyDocumentation(t *testing.T) {
	t.Parallel()

	documentation := EmptyDocumentation()

	if !documentation.IsZero() {
		t.Fatalf("EmptyDocumentation().IsZero() = false, want true")
	}

	if !documentation.IsValid() {
		t.Fatalf("EmptyDocumentation().IsValid() = false, want true")
	}

	if err := documentation.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}
}

// TestNewSummaryDocumentation verifies summary-only convenience construction.
func TestNewSummaryDocumentation(t *testing.T) {
	t.Parallel()

	documentation, err := NewSummaryDocumentation("Run checks.")
	if err != nil {
		t.Fatalf("NewSummaryDocumentation() returned unexpected error: %v", err)
	}

	if got, want := documentation.Summary(), "Run checks."; got != want {
		t.Fatalf("Summary() = %q, want %q", got, want)
	}

	if documentation.HasDescription() || documentation.HasUsage() || documentation.HasNotes() || documentation.HasReferences() {
		t.Fatalf("summary-only documentation has unexpected extra content")
	}
}

// TestNewDocumentationNormalizesText verifies summary, description, and notes
// normalization.
func TestNewDocumentationNormalizesText(t *testing.T) {
	t.Parallel()

	documentation := MustDocumentation(DocumentationSpec{
		Summary:     "  Run   benchmarks.  ",
		Description: "  First line.  \r\n  Second line.  ",
		Notes: []string{
			"  Note   text.  ",
			"   ",
			"\n",
		},
	})

	if got, want := documentation.Summary(), "Run benchmarks."; got != want {
		t.Fatalf("Summary() = %q, want %q", got, want)
	}

	if got, want := documentation.Description(), "First line.\nSecond line."; got != want {
		t.Fatalf("Description() = %q, want %q", got, want)
	}

	assertStringSlicesEqual(t, documentation.Notes(), []string{"Note   text."})
}

// TestNewDocumentationRejectsInvalidDocumentation verifies documentation
// validation and shared textvalidation error preservation.
func TestNewDocumentationRejectsInvalidDocumentation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		spec       DocumentationSpec
		want       error
		wantShared error
	}{
		{
			name: "too long summary",
			spec: DocumentationSpec{
				Summary: strings.Repeat("x", maxDocumentationSummaryLength+1),
			},
			want:       ErrInvalidDocumentation,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name: "summary with control rune",
			spec: DocumentationSpec{
				Summary: "bad\x00summary",
			},
			want:       ErrInvalidDocumentation,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name: "summary with newline",
			spec: DocumentationSpec{
				Summary: "bad\nsummary",
			},
			want:       ErrInvalidDocumentation,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name: "too long description",
			spec: DocumentationSpec{
				Description: strings.Repeat("x", maxDocumentationDescriptionLength+1),
			},
			want:       ErrInvalidDocumentation,
			wantShared: textvalidate.ErrInvalidCompactText,
		},
		{
			name: "description with control rune",
			spec: DocumentationSpec{
				Description: "bad\x00description",
			},
			want:       ErrInvalidDocumentation,
			wantShared: textvalidate.ErrInvalidCompactText,
		},
		{
			name: "invalid usage",
			spec: DocumentationSpec{
				Usage: Usage{syntax: UsageLine("bench\x00run")},
			},
			want:       ErrInvalidDocumentation,
			wantShared: ErrInvalidUsage,
		},
		{
			name: "duplicate notes",
			spec: DocumentationSpec{
				Notes: []string{"same", " same "},
			},
			want: ErrInvalidDocumentation,
		},
		{
			name: "invalid reference",
			spec: DocumentationSpec{
				References: []DocumentationReference{
					{},
				},
			},
			want:       ErrInvalidDocumentation,
			wantShared: ErrInvalidDocumentationReference,
		},
		{
			name: "duplicate reference key",
			spec: DocumentationSpec{
				References: []DocumentationReference{
					MustDocumentationReference(DocumentationReferenceSpec{
						Key:    "same",
						Kind:   DocumentationReferenceDocument,
						Label:  "One",
						Target: "one.md",
					}),
					MustDocumentationReference(DocumentationReferenceSpec{
						Key:    "same",
						Kind:   DocumentationReferenceDocument,
						Label:  "Two",
						Target: "two.md",
					}),
				},
			},
			want: ErrInvalidDocumentation,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewDocumentation(tt.spec)
			if err == nil {
				t.Fatalf("NewDocumentation() returned nil error")
			}

			if !errors.Is(err, tt.want) {
				t.Fatalf("NewDocumentation() error = %v, want %v", err, tt.want)
			}

			if tt.wantShared != nil && !errors.Is(err, tt.wantShared) {
				t.Fatalf("NewDocumentation() error = %v, want shared sentinel %v", err, tt.wantShared)
			}
		})
	}
}

// TestMustDocumentationPanicsForInvalidDocumentation verifies fail-fast
// construction.
func TestMustDocumentationPanicsForInvalidDocumentation(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustDocumentation(DocumentationSpec{
			Summary: "bad\x00summary",
		})
	})
}
