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

import "fmt"

// Validate verifies documentation structural rules.
//
// Empty documentation is valid. Higher-level registry or release policy may
// require public commands to have Summary, Usage, examples, or references.
func (documentation Documentation) Validate() error {
	if documentation.summary != "" {
		if err := validateDocumentationTextLine("summary", documentation.summary, maxDocumentationSummaryLength); err != nil {
			return err
		}
	}

	if documentation.description != "" {
		if err := validateDocumentationBlock("description", documentation.description, maxDocumentationDescriptionLength); err != nil {
			return err
		}
	}

	if !documentation.usage.IsZero() {
		if err := documentation.usage.Validate(); err != nil {
			return fmt.Errorf("%w: invalid usage: %w", ErrInvalidDocumentation, err)
		}
	}

	if err := validateDocumentationNotes(documentation.notes); err != nil {
		return err
	}

	if err := validateDocumentationReferences(documentation.references); err != nil {
		return err
	}

	return nil
}

// validateDocumentationSpecInput checks raw constructor input before
// normalization can hide unsupported single-line or control characters.
func validateDocumentationSpecInput(spec DocumentationSpec) error {
	if spec.Summary != "" {
		if err := validateDocumentationTextLineInput("summary", spec.Summary); err != nil {
			return err
		}
	}

	if spec.Description != "" {
		if err := validateDocumentationBlockInput("description", spec.Description); err != nil {
			return err
		}
	}

	for index, note := range spec.Notes {
		if normalizeDocumentationBlock(note) == "" {
			continue
		}

		if err := validateDocumentationBlockInput(fmt.Sprintf("note %d", index), note); err != nil {
			return err
		}
	}

	return nil
}

// validateDocumentationNotes validates note text and duplicate note content.
func validateDocumentationNotes(notes []string) error {
	seen := make(map[string]struct{}, len(notes))

	for index, note := range notes {
		if err := validateDocumentationBlock(fmt.Sprintf("note %d", index), note, maxDocumentationNoteLength); err != nil {
			return err
		}

		if _, exists := seen[note]; exists {
			return fmt.Errorf("%w: duplicate note %q", ErrInvalidDocumentation, note)
		}

		seen[note] = struct{}{}
	}

	return nil
}

// validateDocumentationReferences validates references and duplicate reference
// keys.
func validateDocumentationReferences(references []DocumentationReference) error {
	seen := make(map[string]struct{}, len(references))

	for index, reference := range references {
		if err := reference.Validate(); err != nil {
			return fmt.Errorf("%w: reference %d: %w", ErrInvalidDocumentation, index, err)
		}

		if _, exists := seen[reference.Key()]; exists {
			return fmt.Errorf("%w: duplicate reference key %q", ErrInvalidDocumentation, reference.Key())
		}

		seen[reference.Key()] = struct{}{}
	}

	return nil
}
