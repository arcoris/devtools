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

// TestDocumentationAccessors verifies read-only helper behavior.
func TestDocumentationAccessors(t *testing.T) {
	t.Parallel()

	documentation := MustDocumentation(DocumentationSpec{
		Summary:     "Run checks.",
		Description: "Runs configured checks.",
		Usage:       MustSimpleUsage("check [flags]"),
		Notes:       []string{"First note.", "Second note."},
		References: []DocumentationReference{
			testDocumentationReference("z-doc"),
			testDocumentationReference("a-doc"),
		},
	})

	if !documentation.HasSummary() || !documentation.HasDescription() || !documentation.HasUsage() {
		t.Fatalf("expected summary, description, and usage to be present")
	}

	note, ok := documentation.Note(1)
	if !ok || note != "Second note." {
		t.Fatalf("Note(1) = %q, %t; want %q, true", note, ok, "Second note.")
	}

	if _, ok := documentation.Note(-1); ok {
		t.Fatalf("Note(-1) ok = true, want false")
	}

	if _, ok := documentation.Note(99); ok {
		t.Fatalf("Note(99) ok = true, want false")
	}

	reference, ok := documentation.ReferenceAt(0)
	if !ok || reference.Key() != "z-doc" {
		t.Fatalf("ReferenceAt(0) = %q, %t; want key %q, true", reference.Key(), ok, "z-doc")
	}

	if _, ok := documentation.ReferenceAt(-1); ok {
		t.Fatalf("ReferenceAt(-1) ok = true, want false")
	}

	if _, ok := documentation.ReferenceAt(99); ok {
		t.Fatalf("ReferenceAt(99) ok = true, want false")
	}

	if !documentation.HasReference("a-doc") {
		t.Fatalf("HasReference(a-doc) = false, want true")
	}

	if documentation.HasReference("missing") {
		t.Fatalf("HasReference(missing) = true, want false")
	}

	assertStringSlicesEqual(t, documentation.ReferenceKeys(), []string{"a-doc", "z-doc"})
}

// TestDocumentationCopySemantics verifies detached slices from constructors,
// accessors, and Spec.
func TestDocumentationCopySemantics(t *testing.T) {
	t.Parallel()

	notes := []string{"note"}
	references := []DocumentationReference{testDocumentationReference("doc")}

	documentation := MustDocumentation(DocumentationSpec{
		Notes:      notes,
		References: references,
	})

	notes[0] = "changed"
	references[0] = testDocumentationReference("changed")

	if got, want := documentation.Notes()[0], "note"; got != want {
		t.Fatalf("note changed through input slice: got %q, want %q", got, want)
	}

	if documentation.HasReference("changed") {
		t.Fatalf("reference changed through input slice")
	}

	outNotes := documentation.Notes()
	outNotes[0] = "changed"
	if got, want := documentation.Notes()[0], "note"; got != want {
		t.Fatalf("note changed through output slice: got %q, want %q", got, want)
	}

	outReferences := documentation.References()
	outReferences[0] = testDocumentationReference("changed")
	if documentation.HasReference("changed") {
		t.Fatalf("reference changed through output slice")
	}

	spec := documentation.Spec()
	spec.Notes[0] = "changed"
	spec.References[0] = testDocumentationReference("changed")
	if got, want := documentation.Notes()[0], "note"; got != want {
		t.Fatalf("note changed through Spec(): got %q, want %q", got, want)
	}

	if documentation.HasReference("changed") {
		t.Fatalf("reference changed through Spec()")
	}
}

// TestDocumentationZeroAccessors verifies zero-value helper behavior.
func TestDocumentationZeroAccessors(t *testing.T) {
	t.Parallel()

	var documentation Documentation

	if !documentation.IsZero() {
		t.Fatalf("zero Documentation IsZero() = false, want true")
	}

	if !documentation.IsValid() {
		t.Fatalf("zero Documentation IsValid() = false, want true")
	}

	if documentation.HasSummary() || documentation.HasDescription() || documentation.HasUsage() || documentation.HasNotes() || documentation.HasReferences() {
		t.Fatalf("zero Documentation reports content")
	}

	if notes := documentation.Notes(); notes != nil {
		t.Fatalf("zero Notes() = %v, want nil", notes)
	}

	if references := documentation.References(); references != nil {
		t.Fatalf("zero References() = %v, want nil", references)
	}

	if _, ok := documentation.Usage(); ok {
		t.Fatalf("zero Usage() ok = true, want false")
	}
}
