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
	"fmt"
	"strings"

	"arcoris.dev/devtools/internal/textvalidate"
)

// Validate verifies documentation reference structural rules.
func (reference DocumentationReference) Validate() error {
	if err := validateDocumentationReferenceKey(reference.key); err != nil {
		return err
	}

	if err := reference.kind.Validate(); err != nil {
		return err
	}

	if err := validateDocumentationReferenceTextLine("reference label", reference.label, maxDocumentationReferenceLabelLength); err != nil {
		return err
	}

	if err := validateDocumentationReferenceTextLine("reference target", reference.target, maxDocumentationReferenceTargetLength); err != nil {
		return err
	}

	return nil
}

// IsKnown reports whether kind is one of the supported reference kinds.
func (kind DocumentationReferenceKind) IsKnown() bool {
	switch kind {
	case DocumentationReferenceCommandID,
		DocumentationReferenceCommandPath,
		DocumentationReferenceDocument,
		DocumentationReferenceURL:
		return true
	default:
		return false
	}
}

// Validate verifies that kind is a supported non-zero reference kind.
func (kind DocumentationReferenceKind) Validate() error {
	if kind == "" {
		return fmt.Errorf("%w: reference kind is empty", ErrInvalidDocumentationReference)
	}

	if kind.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported reference kind %q", ErrInvalidDocumentationReference, kind)
}

// String returns the canonical reference kind string.
func (kind DocumentationReferenceKind) String() string {
	return string(kind)
}

// validateDocumentationReferenceSpecInput checks raw constructor input before
// normalization can hide unsupported single-line characters.
func validateDocumentationReferenceSpecInput(spec DocumentationReferenceSpec) error {
	if err := validateDocumentationReferenceKeyInput(spec.Key); err != nil {
		return err
	}

	if err := validateDocumentationReferenceTextLineInput("reference label", spec.Label); err != nil {
		return err
	}

	if err := validateDocumentationReferenceTextLineInput("reference target", spec.Target); err != nil {
		return err
	}

	return nil
}

// validateDocumentationReferenceKey validates one canonical reference key.
func validateDocumentationReferenceKey(raw string) error {
	if err := textvalidate.ValidateDottedKebabKey(
		raw,
		maxDocumentationReferenceKeyLength,
		maxDocumentationReferenceKeyDepth,
	); err != nil {
		return fmt.Errorf("%w: invalid reference key %q: %w", ErrInvalidDocumentationReference, raw, err)
	}

	return nil
}

// validateDocumentationReferenceKeyInput validates raw reference-key input
// before trim normalization.
func validateDocumentationReferenceKeyInput(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return fmt.Errorf("%w: reference key must not be blank", ErrInvalidDocumentationReference)
	}

	if err := textvalidate.ValidateSingleLineText(raw, len(raw)); err != nil {
		return fmt.Errorf("%w: invalid reference key %q: %w", ErrInvalidDocumentationReference, raw, err)
	}

	return nil
}
