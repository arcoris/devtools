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

// Validate verifies metadata structural rules.
func (metadata Metadata) Validate() error {
	if metadata.owner != "" {
		if err := validateMetadataKey("owner", metadata.owner); err != nil {
			return err
		}
	}

	if metadata.area != "" {
		if err := validateMetadataKey("area", metadata.area); err != nil {
			return err
		}
	}

	if metadata.since != "" {
		if err := validateMetadataText("since", metadata.since, maxMetadataTextLength); err != nil {
			return err
		}
	}

	if metadata.deprecation != nil {
		if err := metadata.deprecation.Validate(); err != nil {
			return fmt.Errorf("%w: invalid deprecation: %w", ErrInvalidMetadata, err)
		}
	}

	return validateAnnotations(metadata.annotations)
}

// validateAnnotations validates all annotation keys and values.
func validateAnnotations(annotations map[string]string) error {
	for key, value := range annotations {
		if err := validateMetadataKey("annotation key", key); err != nil {
			return err
		}

		if err := validateMetadataText("annotation "+key, value, maxAnnotationValueLength); err != nil {
			return err
		}
	}

	return nil
}

// validateMetadataKey validates a stable metadata key and wraps reusable
// textvalidation errors with command metadata diagnostics.
func validateMetadataKey(field string, raw string) error {
	if err := textvalidate.ValidateDottedKebabKey(raw, maxMetadataKeyLength, maxMetadataKeyDepth); err != nil {
		return fmt.Errorf("%w: invalid %s %q: %w", ErrInvalidMetadata, field, raw, err)
	}

	return nil
}

// validateMetadataText validates compact metadata text and wraps reusable
// textvalidation errors with command metadata diagnostics.
func validateMetadataText(field string, raw string, maxLength int) error {
	return validateCommandCompactText(ErrInvalidMetadata, field, raw, maxLength)
}

// validateDeprecationText validates compact deprecation text and wraps reusable
// textvalidation errors with command deprecation diagnostics.
func validateDeprecationText(field string, raw string, maxLength int) error {
	return validateCommandCompactText(ErrInvalidDeprecation, field, raw, maxLength)
}

// validateCommandCompactText validates compact command-model text and applies
// the supplied domain sentinel error.
func validateCommandCompactText(domainErr error, field string, raw string, maxLength int) error {
	if err := textvalidate.ValidateCompactText(raw, maxLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", domainErr, field, err)
	}

	return nil
}
