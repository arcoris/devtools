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
	// maxContextFieldKeyLength is the maximum byte length of one context
	// metadata field key.
	maxContextFieldKeyLength = 255

	// maxContextFieldKeyDepth is the maximum number of dot-separated segments
	// in one context metadata field key.
	maxContextFieldKeyDepth = 32

	// maxContextFieldValueLength is the maximum byte length of one command
	// context metadata field value.
	maxContextFieldValueLength = 4096
)

// Validate verifies command context structural rules.
func (commandContext Context) Validate() error {
	if commandContext.base == nil {
		return fmt.Errorf("%w: base context is not set", ErrInvalidContext)
	}

	if commandContext.node.Kind() == "" {
		return fmt.Errorf("%w: node is not set", ErrInvalidContext)
	}

	if err := commandContext.node.Validate(); err != nil {
		return fmt.Errorf("%w: invalid node: %w", ErrInvalidContext, err)
	}

	if commandContext.startedAt.IsZero() {
		return fmt.Errorf("%w: started timestamp is not set", ErrInvalidContext)
	}

	if err := commandContext.invocation.Validate(); err != nil {
		return fmt.Errorf("%w: invalid invocation: %w", ErrInvalidContext, err)
	}

	return validateContextFields(commandContext.fields)
}

// validateContextFields validates context metadata fields.
func validateContextFields(fields map[string]string) error {
	for key, value := range fields {
		if err := validateContextFieldKey("field key", key); err != nil {
			return err
		}

		if err := validateContextText("field "+key, value, maxContextFieldValueLength); err != nil {
			return err
		}
	}

	return nil
}

// validateContextFieldKey validates a dot-separated context metadata key.
func validateContextFieldKey(field string, raw string) error {
	if err := textvalidate.ValidateDottedKebabKey(raw, maxContextFieldKeyLength, maxContextFieldKeyDepth); err != nil {
		return fmt.Errorf("%w: invalid %s %q: %w", ErrInvalidContext, field, raw, err)
	}

	return nil
}

// validateContextText validates compact UTF-8 command context text.
func validateContextText(field string, raw string, maxLength int) error {
	if err := textvalidate.ValidateCompactText(raw, maxLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", ErrInvalidContext, field, err)
	}

	return nil
}
