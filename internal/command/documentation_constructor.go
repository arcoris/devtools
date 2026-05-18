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

// EmptyDocumentation returns a valid empty documentation value.
func EmptyDocumentation() Documentation {
	return Documentation{}
}

// NewDocumentation validates spec and returns Documentation.
func NewDocumentation(spec DocumentationSpec) (Documentation, error) {
	if err := validateDocumentationSpecInput(spec); err != nil {
		return Documentation{}, err
	}

	documentation := Documentation{
		summary:     normalizeDocumentationSingleLine(spec.Summary),
		description: normalizeDocumentationBlock(spec.Description),
		usage:       spec.Usage,
		notes:       normalizeDocumentationNotes(spec.Notes),
		references:  cloneDocumentationReferences(spec.References),
	}

	if err := documentation.Validate(); err != nil {
		return Documentation{}, err
	}

	return documentation, nil
}

// MustDocumentation validates spec and returns Documentation.
//
// MustDocumentation panics on invalid input. It is intended for static command
// definitions and tests where invalid documentation is a programmer error.
func MustDocumentation(spec DocumentationSpec) Documentation {
	documentation, err := NewDocumentation(spec)
	if err != nil {
		panic(err)
	}

	return documentation
}

// NewSummaryDocumentation creates documentation with only a summary.
func NewSummaryDocumentation(summary string) (Documentation, error) {
	return NewDocumentation(DocumentationSpec{
		Summary: summary,
	})
}

// MustSummaryDocumentation creates documentation with only a summary and panics
// on invalid input.
func MustSummaryDocumentation(summary string) Documentation {
	return MustDocumentation(DocumentationSpec{
		Summary: summary,
	})
}
