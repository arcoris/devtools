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

// TestNewDocumentationReferenceAcceptsValidReference verifies reference
// construction and normalization.
func TestNewDocumentationReferenceAcceptsValidReference(t *testing.T) {
	t.Parallel()

	reference, err := NewDocumentationReference(DocumentationReferenceSpec{
		Key:    " bench-run ",
		Kind:   DocumentationReferenceCommandPath,
		Label:  "  Run   benchmarks  ",
		Target: "  bench   run  ",
	})
	if err != nil {
		t.Fatalf("NewDocumentationReference() returned unexpected error: %v", err)
	}

	if got, want := reference.Key(), "bench-run"; got != want {
		t.Fatalf("Key() = %q, want %q", got, want)
	}

	if got, want := reference.Kind(), DocumentationReferenceCommandPath; got != want {
		t.Fatalf("Kind() = %q, want %q", got, want)
	}

	if got, want := reference.Label(), "Run benchmarks"; got != want {
		t.Fatalf("Label() = %q, want %q", got, want)
	}

	if got, want := reference.Target(), "bench run"; got != want {
		t.Fatalf("Target() = %q, want %q", got, want)
	}

	if !reference.IsValid() {
		t.Fatalf("IsValid() = false, want true")
	}
}

// TestNewDocumentationReferenceRejectsInvalidReference verifies reference
// validation and shared textvalidation error preservation.
func TestNewDocumentationReferenceRejectsInvalidReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		spec       DocumentationReferenceSpec
		wantShared error
	}{
		{
			name: "empty key",
			spec: DocumentationReferenceSpec{
				Kind:   DocumentationReferenceDocument,
				Label:  "Doc",
				Target: "doc.md",
			},
		},
		{
			name: "invalid key",
			spec: DocumentationReferenceSpec{
				Key:    "Bad_Key",
				Kind:   DocumentationReferenceDocument,
				Label:  "Doc",
				Target: "doc.md",
			},
			wantShared: textvalidate.ErrInvalidDottedKebabKey,
		},
		{
			name: "key newline",
			spec: DocumentationReferenceSpec{
				Key:    "doc\nkey",
				Kind:   DocumentationReferenceDocument,
				Label:  "Doc",
				Target: "doc.md",
			},
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name: "empty kind",
			spec: DocumentationReferenceSpec{
				Key:    "doc",
				Label:  "Doc",
				Target: "doc.md",
			},
		},
		{
			name: "unknown kind",
			spec: DocumentationReferenceSpec{
				Key:    "doc",
				Kind:   DocumentationReferenceKind("unknown"),
				Label:  "Doc",
				Target: "doc.md",
			},
		},
		{
			name: "blank label",
			spec: DocumentationReferenceSpec{
				Key:    "doc",
				Kind:   DocumentationReferenceDocument,
				Label:  "   ",
				Target: "doc.md",
			},
		},
		{
			name: "label newline",
			spec: DocumentationReferenceSpec{
				Key:    "doc",
				Kind:   DocumentationReferenceDocument,
				Label:  "Doc\nlabel",
				Target: "doc.md",
			},
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name: "blank target",
			spec: DocumentationReferenceSpec{
				Key:    "doc",
				Kind:   DocumentationReferenceDocument,
				Label:  "Doc",
				Target: "   ",
			},
		},
		{
			name: "too long target",
			spec: DocumentationReferenceSpec{
				Key:    "doc",
				Kind:   DocumentationReferenceDocument,
				Label:  "Doc",
				Target: strings.Repeat("x", maxDocumentationReferenceTargetLength+1),
			},
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name: "target control rune",
			spec: DocumentationReferenceSpec{
				Key:    "doc",
				Kind:   DocumentationReferenceDocument,
				Label:  "Doc",
				Target: "bad\x00target",
			},
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewDocumentationReference(tt.spec)
			if err == nil {
				t.Fatalf("NewDocumentationReference() returned nil error")
			}

			if !errors.Is(err, ErrInvalidDocumentationReference) {
				t.Fatalf("NewDocumentationReference() error = %v, want ErrInvalidDocumentationReference", err)
			}

			if tt.wantShared != nil && !errors.Is(err, tt.wantShared) {
				t.Fatalf("NewDocumentationReference() error = %v, want shared sentinel %v", err, tt.wantShared)
			}
		})
	}
}

// TestMustDocumentationReferencePanicsForInvalidReference verifies fail-fast
// construction.
func TestMustDocumentationReferencePanicsForInvalidReference(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustDocumentationReference(DocumentationReferenceSpec{})
	})
}
