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

import "testing"

// TestDocumentationNormalizationHelpers verifies normalization helpers.
func TestDocumentationNormalizationHelpers(t *testing.T) {
	t.Parallel()

	if got, want := normalizeDocumentationSingleLine("  Run   checks. "), "Run checks."; got != want {
		t.Fatalf("normalizeDocumentationSingleLine() = %q, want %q", got, want)
	}

	if got, want := normalizeDocumentationBlock("  One.  \r\n  Two.  "), "One.\nTwo."; got != want {
		t.Fatalf("normalizeDocumentationBlock() = %q, want %q", got, want)
	}

	if got, want := normalizeDocumentationReferenceKey(" doc "), "doc"; got != want {
		t.Fatalf("normalizeDocumentationReferenceKey() = %q, want %q", got, want)
	}

	notes := normalizeDocumentationNotes([]string{"  One.  ", " ", "  Two.  "})
	assertStringSlicesEqual(t, notes, []string{"One.", "Two."})

	if normalizeDocumentationNotes(nil) != nil {
		t.Fatalf("normalizeDocumentationNotes(nil) must return nil")
	}
}
