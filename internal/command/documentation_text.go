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

// validateDocumentationTextLine validates canonical one-line documentation
// text and wraps reusable text validation errors in documentation diagnostics.
func validateDocumentationTextLine(field string, raw string, maxLength int) error {
	return validateDocumentationSingleLine(ErrInvalidDocumentation, field, raw, maxLength)
}

// validateDocumentationReferenceTextLine validates canonical one-line reference
// text and wraps reusable text validation errors in reference diagnostics.
func validateDocumentationReferenceTextLine(field string, raw string, maxLength int) error {
	return validateDocumentationSingleLine(ErrInvalidDocumentationReference, field, raw, maxLength)
}

// validateDocumentationTextLineInput validates raw one-line documentation
// constructor input before normalization.
func validateDocumentationTextLineInput(field string, raw string) error {
	return validateDocumentationSingleLineInput(ErrInvalidDocumentation, field, raw)
}

// validateDocumentationReferenceTextLineInput validates raw one-line reference
// constructor input before normalization.
func validateDocumentationReferenceTextLineInput(field string, raw string) error {
	return validateDocumentationSingleLineInput(ErrInvalidDocumentationReference, field, raw)
}

// validateDocumentationSingleLine validates canonical one-line text for a
// command documentation domain.
func validateDocumentationSingleLine(domainErr error, field string, raw string, maxLength int) error {
	if strings.TrimSpace(raw) == "" {
		return fmt.Errorf("%w: %s must not be blank", domainErr, field)
	}

	if err := textvalidate.ValidateSingleLineText(raw, maxLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", domainErr, field, err)
	}

	if raw != normalizeDocumentationSingleLine(raw) {
		return fmt.Errorf("%w: %s is not canonical", domainErr, field)
	}

	return nil
}

// validateDocumentationSingleLineInput validates raw one-line text before it is
// normalized. Length is validated after normalization by Documentation.Validate
// or DocumentationReference.Validate.
func validateDocumentationSingleLineInput(domainErr error, field string, raw string) error {
	if strings.TrimSpace(raw) == "" {
		return fmt.Errorf("%w: %s must not be blank", domainErr, field)
	}

	if err := textvalidate.ValidateSingleLineText(raw, len(raw)); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", domainErr, field, err)
	}

	return nil
}

// validateDocumentationBlock validates canonical possibly multi-line
// documentation text.
func validateDocumentationBlock(field string, raw string, maxLength int) error {
	if strings.TrimSpace(raw) == "" {
		return fmt.Errorf("%w: %s must not be blank", ErrInvalidDocumentation, field)
	}

	if err := textvalidate.ValidateCompactText(raw, maxLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", ErrInvalidDocumentation, field, err)
	}

	if raw != normalizeDocumentationBlock(raw) {
		return fmt.Errorf("%w: %s is not canonical", ErrInvalidDocumentation, field)
	}

	return nil
}

// validateDocumentationBlockInput validates raw block text before it is
// normalized. Length is validated after normalization by Documentation.Validate.
func validateDocumentationBlockInput(field string, raw string) error {
	if strings.TrimSpace(raw) == "" {
		return fmt.Errorf("%w: %s must not be blank", ErrInvalidDocumentation, field)
	}

	if err := textvalidate.ValidateCompactText(raw, len(raw)); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", ErrInvalidDocumentation, field, err)
	}

	return nil
}
