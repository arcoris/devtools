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

// TestDocumentationValidateRejectsInvalidStoredDocumentation verifies
// validation of manually constructed values.
func TestDocumentationValidateRejectsInvalidStoredDocumentation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		value      Documentation
		want       error
		wantShared error
	}{
		{
			name: "summary not canonical",
			value: Documentation{
				summary: "Run   checks.",
			},
			want: ErrInvalidDocumentation,
		},
		{
			name: "summary control",
			value: Documentation{
				summary: "Run\x00checks.",
			},
			want:       ErrInvalidDocumentation,
			wantShared: textvalidate.ErrInvalidSingleLineText,
		},
		{
			name: "description not canonical",
			value: Documentation{
				description: " First line.",
			},
			want: ErrInvalidDocumentation,
		},
		{
			name: "description control",
			value: Documentation{
				description: "bad\x00description",
			},
			want:       ErrInvalidDocumentation,
			wantShared: textvalidate.ErrInvalidCompactText,
		},
		{
			name: "invalid usage",
			value: Documentation{
				usage: Usage{syntax: UsageLine("bench\x00run")},
			},
			want:       ErrInvalidDocumentation,
			wantShared: ErrInvalidUsage,
		},
		{
			name: "note control",
			value: Documentation{
				notes: []string{"bad\x00note"},
			},
			want:       ErrInvalidDocumentation,
			wantShared: textvalidate.ErrInvalidCompactText,
		},
		{
			name: "duplicate notes",
			value: Documentation{
				notes: []string{"same", "same"},
			},
			want: ErrInvalidDocumentation,
		},
		{
			name: "invalid reference",
			value: Documentation{
				references: []DocumentationReference{{}},
			},
			want:       ErrInvalidDocumentation,
			wantShared: ErrInvalidDocumentationReference,
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

			if !errors.Is(err, tt.want) {
				t.Fatalf("Validate() error = %v, want %v", err, tt.want)
			}

			if tt.wantShared != nil && !errors.Is(err, tt.wantShared) {
				t.Fatalf("Validate() error = %v, want shared sentinel %v", err, tt.wantShared)
			}
		})
	}
}

// TestDocumentationTextValidators verifies documentation text validation
// helpers.
func TestDocumentationTextValidators(t *testing.T) {
	t.Parallel()

	if err := validateDocumentationTextLine("summary", "Run checks.", maxDocumentationSummaryLength); err != nil {
		t.Fatalf("validateDocumentationTextLine() returned unexpected error: %v", err)
	}

	err := validateDocumentationTextLine("summary", "Run\nchecks.", maxDocumentationSummaryLength)
	if !errors.Is(err, ErrInvalidDocumentation) || !errors.Is(err, textvalidate.ErrInvalidSingleLineText) {
		t.Fatalf("validateDocumentationTextLine(newline) error = %v, want domain and shared sentinels", err)
	}

	if err := validateDocumentationBlock("description", "Line one.\nLine two.", maxDocumentationDescriptionLength); err != nil {
		t.Fatalf("validateDocumentationBlock() returned unexpected error: %v", err)
	}

	err = validateDocumentationBlock("description", "bad\x00text", maxDocumentationDescriptionLength)
	if !errors.Is(err, ErrInvalidDocumentation) || !errors.Is(err, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("validateDocumentationBlock(control) error = %v, want domain and shared sentinels", err)
	}

	invalidUTF8 := string([]byte{0xff, 0xfe})
	err = validateDocumentationBlock("description", invalidUTF8, maxDocumentationDescriptionLength)
	if !errors.Is(err, ErrInvalidDocumentation) || !errors.Is(err, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("validateDocumentationBlock(invalid UTF-8) error = %v, want domain and shared sentinels", err)
	}

	if err := validateDocumentationReferenceTextLine("reference label", "Docs", maxDocumentationReferenceLabelLength); err != nil {
		t.Fatalf("validateDocumentationReferenceTextLine() returned unexpected error: %v", err)
	}

	err = validateDocumentationReferenceTextLine("reference label", " ", maxDocumentationReferenceLabelLength)
	if !errors.Is(err, ErrInvalidDocumentationReference) {
		t.Fatalf("validateDocumentationReferenceTextLine(blank) error = %v, want ErrInvalidDocumentationReference", err)
	}
}
