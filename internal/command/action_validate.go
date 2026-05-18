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

	"arcoris.dev/devtools/internal/textvalidate"
)

const (
	// maxActionArgumentLength is the maximum byte length of one positional
	// argument stored in ActionRequest.
	maxActionArgumentLength = 4096

	// maxActionMessageLength is the maximum byte length of compact action
	// messages and warnings.
	maxActionMessageLength = 4096

	// maxActionDescriptionLength is the maximum byte length of compact artifact
	// descriptions.
	maxActionDescriptionLength = 4096

	// maxActionArtifactPathLength is the maximum byte length of one artifact
	// path reference.
	maxActionArtifactPathLength = 8192

	// maxActionFieldKeyLength is the maximum byte length of one action metadata
	// key.
	maxActionFieldKeyLength = 255

	// maxActionFieldKeyDepth is the maximum number of dot-separated segments in
	// one action metadata key.
	maxActionFieldKeyDepth = 32

	// maxActionFieldValueLength is the maximum byte length of one action field
	// value.
	maxActionFieldValueLength = 4096
)

// validateActionRequestFields validates request metadata fields.
func validateActionRequestFields(fields map[string]string) error {
	return validateActionFields(ErrInvalidActionRequest, fields)
}

// validateActionResultFields validates result metadata fields.
func validateActionResultFields(fields map[string]string) error {
	return validateActionFields(ErrInvalidActionResult, fields)
}

// validateActionFields validates machine-facing action metadata fields.
func validateActionFields(domainErr error, fields map[string]string) error {
	for key, value := range fields {
		if err := validateActionFieldKey(domainErr, "field key", key); err != nil {
			return err
		}

		if err := validateActionText(domainErr, "field "+key, value, maxActionFieldValueLength); err != nil {
			return err
		}
	}

	return nil
}

// validateActionRequestFieldKey validates a request metadata key.
func validateActionRequestFieldKey(field string, raw string) error {
	return validateActionFieldKey(ErrInvalidActionRequest, field, raw)
}

// validateActionResultFieldKey validates a result metadata key.
func validateActionResultFieldKey(field string, raw string) error {
	return validateActionFieldKey(ErrInvalidActionResult, field, raw)
}

// validateActionFieldKey validates a dot-separated machine-facing field key.
func validateActionFieldKey(domainErr error, field string, raw string) error {
	if err := textvalidate.ValidateDottedKebabKey(raw, maxActionFieldKeyLength, maxActionFieldKeyDepth); err != nil {
		return fmt.Errorf("%w: invalid %s %q: %w", domainErr, field, raw, err)
	}

	return nil
}

// validateActionRequestText validates compact request text.
func validateActionRequestText(field string, raw string, maxLength int) error {
	return validateActionText(ErrInvalidActionRequest, field, raw, maxLength)
}

// validateActionResultText validates compact result text.
func validateActionResultText(field string, raw string, maxLength int) error {
	return validateActionText(ErrInvalidActionResult, field, raw, maxLength)
}

// validateActionText validates compact UTF-8 action text and applies the
// supplied domain sentinel error.
//
// Action text may contain whitespace and punctuation, but it must not contain
// disallowed control characters. Tab, newline, and carriage return are allowed
// because command output often contains multi-line diagnostic text.
func validateActionText(domainErr error, field string, raw string, maxLength int) error {
	if err := textvalidate.ValidateCompactText(raw, maxLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", domainErr, field, err)
	}

	return nil
}
