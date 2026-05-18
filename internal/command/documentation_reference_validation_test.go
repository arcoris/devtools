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

// TestDocumentationReferenceAccessors verifies reference helper behavior.
func TestDocumentationReferenceAccessors(t *testing.T) {
	t.Parallel()

	reference := testDocumentationReference("doc")

	if reference.IsZero() {
		t.Fatalf("IsZero() = true, want false")
	}

	if got, want := reference.String(), "docs/doc.md"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	spec := reference.Spec()
	spec.Key = "changed"
	if got, want := reference.Key(), "doc"; got != want {
		t.Fatalf("reference changed through Spec(): got %q, want %q", got, want)
	}
}

// TestDocumentationReferenceZeroAccessors verifies zero-value helper behavior.
func TestDocumentationReferenceZeroAccessors(t *testing.T) {
	t.Parallel()

	var reference DocumentationReference

	if !reference.IsZero() {
		t.Fatalf("zero reference IsZero() = false, want true")
	}

	if reference.IsValid() {
		t.Fatalf("zero reference IsValid() = true, want false")
	}
}

// TestDocumentationReferenceValidateRejectsInvalidStoredReference verifies
// validation of manually constructed values.
func TestDocumentationReferenceValidateRejectsInvalidStoredReference(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		value      DocumentationReference
		wantShared error
	}{
		{
			name: "invalid key",
			value: DocumentationReference{
				key:    "Bad_Key",
				kind:   DocumentationReferenceDocument,
				label:  "Doc",
				target: "doc.md",
			},
			wantShared: textvalidate.ErrInvalidDottedKebabKey,
		},
		{
			name: "invalid kind",
			value: DocumentationReference{
				key:    "doc",
				kind:   "unknown",
				label:  "Doc",
				target: "doc.md",
			},
		},
		{
			name: "label not canonical",
			value: DocumentationReference{
				key:    "doc",
				kind:   DocumentationReferenceDocument,
				label:  "Doc   Label",
				target: "doc.md",
			},
		},
		{
			name: "target newline",
			value: DocumentationReference{
				key:    "doc",
				kind:   DocumentationReferenceDocument,
				label:  "Doc",
				target: "doc\npath",
			},
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.value.Validate()
			if err == nil {
				t.Fatalf("Validate() returned nil error")
			}

			if !errors.Is(err, ErrInvalidDocumentationReference) {
				t.Fatalf("Validate() error = %v, want ErrInvalidDocumentationReference", err)
			}

			if tt.wantShared != nil && !errors.Is(err, tt.wantShared) {
				t.Fatalf("Validate() error = %v, want shared sentinel %v", err, tt.wantShared)
			}
		})
	}
}

// TestDocumentationReferenceKindValidation verifies reference kind validation.
func TestDocumentationReferenceKindValidation(t *testing.T) {
	t.Parallel()

	valid := []DocumentationReferenceKind{
		DocumentationReferenceCommandID,
		DocumentationReferenceCommandPath,
		DocumentationReferenceDocument,
		DocumentationReferenceURL,
	}

	for _, kind := range valid {
		kind := kind

		t.Run(kind.String(), func(t *testing.T) {
			t.Parallel()

			if !kind.IsKnown() {
				t.Fatalf("IsKnown() = false, want true")
			}

			if err := kind.Validate(); err != nil {
				t.Fatalf("Validate() returned unexpected error: %v", err)
			}
		})
	}

	invalid := []DocumentationReferenceKind{
		"",
		"unknown",
	}

	for _, kind := range invalid {
		kind := kind

		t.Run("invalid-"+kind.String(), func(t *testing.T) {
			t.Parallel()

			if kind.IsKnown() {
				t.Fatalf("IsKnown() = true, want false")
			}

			err := kind.Validate()
			if err == nil {
				t.Fatalf("Validate() returned nil error")
			}

			if !errors.Is(err, ErrInvalidDocumentationReference) {
				t.Fatalf("Validate() error = %v, want ErrInvalidDocumentationReference", err)
			}
		})
	}
}

// TestCloneDocumentationReferences verifies detached copy behavior.
func TestCloneDocumentationReferences(t *testing.T) {
	t.Parallel()

	original := []DocumentationReference{testDocumentationReference("doc")}

	cloned := cloneDocumentationReferences(original)
	cloned[0] = testDocumentationReference("other")

	if got, want := original[0].Key(), "doc"; got != want {
		t.Fatalf("original changed through clone: got %q, want %q", got, want)
	}

	if cloneDocumentationReferences(nil) != nil {
		t.Fatalf("cloneDocumentationReferences(nil) must return nil")
	}
}
