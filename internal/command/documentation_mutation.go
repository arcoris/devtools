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

// WithSummary returns a validated copy with Summary replaced.
func (documentation Documentation) WithSummary(summary string) (Documentation, error) {
	spec := documentation.spec()
	spec.Summary = summary

	return NewDocumentation(spec)
}

// MustWithSummary returns a validated copy with Summary replaced and panics on
// invalid input.
func (documentation Documentation) MustWithSummary(summary string) Documentation {
	next, err := documentation.WithSummary(summary)
	if err != nil {
		panic(err)
	}

	return next
}

// WithDescription returns a validated copy with Description replaced.
func (documentation Documentation) WithDescription(description string) (Documentation, error) {
	spec := documentation.spec()
	spec.Description = description

	return NewDocumentation(spec)
}

// MustWithDescription returns a validated copy with Description replaced and
// panics on invalid input.
func (documentation Documentation) MustWithDescription(description string) Documentation {
	next, err := documentation.WithDescription(description)
	if err != nil {
		panic(err)
	}

	return next
}

// WithUsage returns a validated copy with Usage replaced.
func (documentation Documentation) WithUsage(usage Usage) (Documentation, error) {
	spec := documentation.spec()
	spec.Usage = usage

	return NewDocumentation(spec)
}

// MustWithUsage returns a validated copy with Usage replaced and panics on
// invalid input.
func (documentation Documentation) MustWithUsage(usage Usage) Documentation {
	next, err := documentation.WithUsage(usage)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutUsage returns a validated copy without Usage.
func (documentation Documentation) WithoutUsage() Documentation {
	spec := documentation.spec()
	spec.Usage = Usage{}

	return MustDocumentation(spec)
}

// WithNotes returns a validated copy with Notes replaced.
func (documentation Documentation) WithNotes(notes []string) (Documentation, error) {
	spec := documentation.spec()
	spec.Notes = cloneStringSlice(notes)

	return NewDocumentation(spec)
}

// MustWithNotes returns a validated copy with Notes replaced and panics on
// invalid input.
func (documentation Documentation) MustWithNotes(notes []string) Documentation {
	next, err := documentation.WithNotes(notes)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutNotes returns a validated copy without notes.
func (documentation Documentation) WithoutNotes() Documentation {
	spec := documentation.spec()
	spec.Notes = nil

	return MustDocumentation(spec)
}

// WithNote returns a validated copy with one note appended.
func (documentation Documentation) WithNote(note string) (Documentation, error) {
	spec := documentation.spec()
	spec.Notes = append(spec.Notes, note)

	return NewDocumentation(spec)
}

// MustWithNote returns a validated copy with one note appended and panics on
// invalid input.
func (documentation Documentation) MustWithNote(note string) Documentation {
	next, err := documentation.WithNote(note)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutNote returns a validated copy without matching note text.
//
// The input is normalized before comparison. If the note is not present, the
// original documentation is returned as a detached validated copy.
func (documentation Documentation) WithoutNote(note string) Documentation {
	normalized := normalizeDocumentationBlock(note)

	spec := documentation.spec()
	nextNotes := make([]string, 0, len(spec.Notes))

	for _, existing := range spec.Notes {
		if existing == normalized {
			continue
		}

		nextNotes = append(nextNotes, existing)
	}

	spec.Notes = nextNotes

	return MustDocumentation(spec)
}

// WithReferences returns a validated copy with References replaced.
func (documentation Documentation) WithReferences(references []DocumentationReference) (Documentation, error) {
	spec := documentation.spec()
	spec.References = cloneDocumentationReferences(references)

	return NewDocumentation(spec)
}

// MustWithReferences returns a validated copy with References replaced and
// panics on invalid input.
func (documentation Documentation) MustWithReferences(references []DocumentationReference) Documentation {
	next, err := documentation.WithReferences(references)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutReferences returns a validated copy without references.
func (documentation Documentation) WithoutReferences() Documentation {
	spec := documentation.spec()
	spec.References = nil

	return MustDocumentation(spec)
}

// WithReference returns a validated copy with one reference appended or
// replaced by key.
func (documentation Documentation) WithReference(reference DocumentationReference) (Documentation, error) {
	spec := documentation.spec()

	replaced := false
	for index, existing := range spec.References {
		if existing.Key() == reference.Key() {
			spec.References[index] = reference
			replaced = true

			break
		}
	}

	if !replaced {
		spec.References = append(spec.References, reference)
	}

	return NewDocumentation(spec)
}

// MustWithReference returns a validated copy with one reference appended or
// replaced and panics on invalid input.
func (documentation Documentation) MustWithReference(reference DocumentationReference) Documentation {
	next, err := documentation.WithReference(reference)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutReference returns a validated copy without the reference key.
func (documentation Documentation) WithoutReference(key string) Documentation {
	spec := documentation.spec()
	nextReferences := make([]DocumentationReference, 0, len(spec.References))

	for _, reference := range spec.References {
		if reference.Key() == key {
			continue
		}

		nextReferences = append(nextReferences, reference)
	}

	spec.References = nextReferences

	return MustDocumentation(spec)
}
