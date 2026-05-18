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

import "sort"

// Spec returns a detached construction spec.
func (documentation Documentation) Spec() DocumentationSpec {
	return documentation.spec()
}

// IsZero reports whether documentation has no content.
func (documentation Documentation) IsZero() bool {
	return documentation.summary == "" &&
		documentation.description == "" &&
		documentation.usage.IsZero() &&
		len(documentation.notes) == 0 &&
		len(documentation.references) == 0
}

// IsValid reports whether documentation satisfies structural rules.
func (documentation Documentation) IsValid() bool {
	return documentation.Validate() == nil
}

// Summary returns the short one-line documentation summary.
func (documentation Documentation) Summary() string {
	return documentation.summary
}

// HasSummary reports whether Summary is set.
func (documentation Documentation) HasSummary() bool {
	return documentation.summary != ""
}

// Description returns the long documentation description.
func (documentation Documentation) Description() string {
	return documentation.description
}

// HasDescription reports whether Description is set.
func (documentation Documentation) HasDescription() bool {
	return documentation.description != ""
}

// Usage returns the structured usage declaration and whether it is set.
func (documentation Documentation) Usage() (Usage, bool) {
	if documentation.usage.IsZero() {
		return Usage{}, false
	}

	return documentation.usage, true
}

// HasUsage reports whether Usage is set.
func (documentation Documentation) HasUsage() bool {
	return !documentation.usage.IsZero()
}

// Notes returns a detached copy of documentation notes.
func (documentation Documentation) Notes() []string {
	return cloneStringSlice(documentation.notes)
}

// Note returns the documentation note at index.
//
// The second return value is false when index is out of range. Note never
// panics.
func (documentation Documentation) Note(index int) (string, bool) {
	if index < 0 || index >= len(documentation.notes) {
		return "", false
	}

	return documentation.notes[index], true
}

// HasNotes reports whether notes are present.
func (documentation Documentation) HasNotes() bool {
	return len(documentation.notes) > 0
}

// NoteCount returns the number of documentation notes.
func (documentation Documentation) NoteCount() int {
	return len(documentation.notes)
}

// References returns a detached copy of documentation references.
func (documentation Documentation) References() []DocumentationReference {
	return cloneDocumentationReferences(documentation.references)
}

// ReferenceAt returns the documentation reference at index.
//
// The second return value is false when index is out of range. ReferenceAt
// never panics.
func (documentation Documentation) ReferenceAt(index int) (DocumentationReference, bool) {
	if index < 0 || index >= len(documentation.references) {
		return DocumentationReference{}, false
	}

	return documentation.references[index], true
}

// HasReferences reports whether references are present.
func (documentation Documentation) HasReferences() bool {
	return len(documentation.references) > 0
}

// ReferenceCount returns the number of documentation references.
func (documentation Documentation) ReferenceCount() int {
	return len(documentation.references)
}

// Reference returns a documentation reference by key.
func (documentation Documentation) Reference(key string) (DocumentationReference, bool) {
	for _, reference := range documentation.references {
		if reference.Key() == key {
			return reference, true
		}
	}

	return DocumentationReference{}, false
}

// HasReference reports whether a documentation reference with key exists.
func (documentation Documentation) HasReference(key string) bool {
	_, ok := documentation.Reference(key)

	return ok
}

// ReferenceKeys returns reference keys in deterministic lexical order.
func (documentation Documentation) ReferenceKeys() []string {
	keys := make([]string, 0, len(documentation.references))
	for _, reference := range documentation.references {
		keys = append(keys, reference.Key())
	}

	sort.Strings(keys)

	return keys
}

// spec returns a detached construction spec.
func (documentation Documentation) spec() DocumentationSpec {
	return DocumentationSpec{
		Summary:     documentation.summary,
		Description: documentation.description,
		Usage:       documentation.usage,
		Notes:       cloneStringSlice(documentation.notes),
		References:  cloneDocumentationReferences(documentation.references),
	}
}
