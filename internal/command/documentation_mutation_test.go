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

// TestDocumentationMutationHelpers verifies immutable-style update helpers.
func TestDocumentationMutationHelpers(t *testing.T) {
	t.Parallel()

	original := EmptyDocumentation()
	documentation := original.
		MustWithSummary("Run checks.").
		MustWithDescription("Runs configured checks.").
		MustWithUsage(MustSimpleUsage("check [flags]")).
		MustWithNotes([]string{"  First note.  ", " "}).
		MustWithNote("Second note.").
		MustWithReferences([]DocumentationReference{testDocumentationReference("doc")}).
		MustWithReference(testDocumentationReference("extra-doc"))

	if !original.IsZero() {
		t.Fatalf("original documentation changed")
	}

	assertStringSlicesEqual(t, documentation.Notes(), []string{"First note.", "Second note."})
	assertStringSlicesEqual(t, documentation.ReferenceKeys(), []string{"doc", "extra-doc"})

	if documentation.WithoutUsage().HasUsage() {
		t.Fatalf("WithoutUsage() still has usage")
	}

	withoutNote := documentation.WithoutNote("  First note.  ")
	assertStringSlicesEqual(t, withoutNote.Notes(), []string{"Second note."})

	if documentation.WithoutNotes().HasNotes() {
		t.Fatalf("WithoutNotes() still has notes")
	}

	if documentation.WithoutReference("doc").HasReference("doc") {
		t.Fatalf("WithoutReference() still has doc reference")
	}

	if documentation.WithoutReferences().HasReferences() {
		t.Fatalf("WithoutReferences() still has references")
	}
}

// TestDocumentationWithReferenceReplacesByKey verifies reference upsert
// behavior.
func TestDocumentationWithReferenceReplacesByKey(t *testing.T) {
	t.Parallel()

	documentation := EmptyDocumentation().
		MustWithReference(MustDocumentationReference(DocumentationReferenceSpec{
			Key:    "doc",
			Kind:   DocumentationReferenceDocument,
			Label:  "Old",
			Target: "old.md",
		})).
		MustWithReference(MustDocumentationReference(DocumentationReferenceSpec{
			Key:    "doc",
			Kind:   DocumentationReferenceDocument,
			Label:  "New",
			Target: "new.md",
		}))

	if got, want := documentation.ReferenceCount(), 1; got != want {
		t.Fatalf("ReferenceCount() = %d, want %d", got, want)
	}

	reference, ok := documentation.Reference("doc")
	if !ok {
		t.Fatalf("Reference(doc) ok = false, want true")
	}

	if got, want := reference.Target(), "new.md"; got != want {
		t.Fatalf("Reference target = %q, want %q", got, want)
	}
}

// TestDocumentationMutationRejectsInvalidUpdates verifies update helper
// validation.
func TestDocumentationMutationRejectsInvalidUpdates(t *testing.T) {
	t.Parallel()

	documentation := MustDocumentation(DocumentationSpec{
		Summary: "Run checks.",
		Notes:   []string{"note"},
	})

	tests := []struct {
		name       string
		fn         func() (Documentation, error)
		want       error
		wantShared error
	}{
		{
			name: "invalid summary",
			fn: func() (Documentation, error) {
				return documentation.WithSummary("bad\nsummary")
			},
			want:       ErrInvalidDocumentation,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name: "invalid description",
			fn: func() (Documentation, error) {
				return documentation.WithDescription("bad\x00description")
			},
			want:       ErrInvalidDocumentation,
			wantShared: textvalidate.ErrInvalidCompactText,
		},
		{
			name: "invalid usage",
			fn: func() (Documentation, error) {
				return documentation.WithUsage(Usage{syntax: UsageLine("bench\x00run")})
			},
			want:       ErrInvalidDocumentation,
			wantShared: ErrInvalidUsage,
		},
		{
			name: "duplicate note",
			fn: func() (Documentation, error) {
				return documentation.WithNote(" note ")
			},
			want: ErrInvalidDocumentation,
		},
		{
			name: "duplicate references",
			fn: func() (Documentation, error) {
				return documentation.WithReferences([]DocumentationReference{
					testDocumentationReference("doc"),
					testDocumentationReference("doc"),
				})
			},
			want: ErrInvalidDocumentation,
		},
		{
			name: "invalid reference",
			fn: func() (Documentation, error) {
				return documentation.WithReference(DocumentationReference{})
			},
			want:       ErrInvalidDocumentation,
			wantShared: ErrInvalidDocumentationReference,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tt.fn()
			if err == nil {
				t.Fatalf("%s returned nil error", tt.name)
			}

			if !errors.Is(err, tt.want) {
				t.Fatalf("%s error = %v, want %v", tt.name, err, tt.want)
			}

			if tt.wantShared != nil && !errors.Is(err, tt.wantShared) {
				t.Fatalf("%s error = %v, want shared sentinel %v", tt.name, err, tt.wantShared)
			}
		})
	}
}

// TestDocumentationMustMutationHelpersPanic verifies fail-fast update helpers.
func TestDocumentationMustMutationHelpersPanic(t *testing.T) {
	t.Parallel()

	documentation := EmptyDocumentation()

	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "summary",
			fn: func() {
				_ = documentation.MustWithSummary("bad\nsummary")
			},
		},
		{
			name: "description",
			fn: func() {
				_ = documentation.MustWithDescription("bad\x00description")
			},
		},
		{
			name: "usage",
			fn: func() {
				_ = documentation.MustWithUsage(Usage{syntax: UsageLine("bad\x00usage")})
			},
		},
		{
			name: "notes",
			fn: func() {
				_ = documentation.MustWithNotes([]string{"same", "same"})
			},
		},
		{
			name: "note",
			fn: func() {
				base := documentation.MustWithNote("same")
				_ = base.MustWithNote("same")
			},
		},
		{
			name: "references",
			fn: func() {
				_ = documentation.MustWithReferences([]DocumentationReference{{}})
			},
		},
		{
			name: "reference",
			fn: func() {
				_ = documentation.MustWithReference(DocumentationReference{})
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assertPanics(t, tt.fn)
		})
	}
}
