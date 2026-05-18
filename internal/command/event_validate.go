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

// Validate verifies event structural rules.
func (event Event) Validate() error {
	if event.hasID {
		if err := event.id.Validate(); err != nil {
			return err
		}
	}

	if err := event.kind.Validate(); err != nil {
		return err
	}

	if err := event.severity.Validate(); err != nil {
		return err
	}

	if event.occurredAt.IsZero() {
		return fmt.Errorf("%w: occurred timestamp must not be zero", ErrInvalidEvent)
	}

	if event.commandID != "" {
		if err := event.commandID.Validate(); err != nil {
			return fmt.Errorf("%w: invalid command id: %w", ErrInvalidEvent, err)
		}
	}

	if event.message != "" {
		if err := validateEventBlock("message", event.message, maxEventMessageLength); err != nil {
			return err
		}
	}

	if err := validateEventFields(event.fields); err != nil {
		return err
	}

	if err := validateEventArtifacts(event.artifacts); err != nil {
		return err
	}

	if event.result != nil {
		if err := event.result.Validate(); err != nil {
			return fmt.Errorf("%w: invalid result: %w", ErrInvalidEvent, err)
		}
	}

	if err := validateEventLabels(event.labels); err != nil {
		return err
	}

	if err := event.metadata.Validate(); err != nil {
		return fmt.Errorf("%w: invalid metadata: %w", ErrInvalidEvent, err)
	}

	if err := event.visibility.Validate(); err != nil {
		return fmt.Errorf("%w: invalid visibility: %w", ErrInvalidEvent, err)
	}

	return nil
}

// validateEventFields validates event fields.
func validateEventFields(fields map[string]string) error {
	for key, value := range fields {
		if err := validateEventFieldKey(key); err != nil {
			return err
		}

		if err := validateEventBlock("field "+key, value, maxEventFieldValueLength); err != nil {
			return err
		}
	}

	return nil
}

// validateEventFieldKey validates one field key.
func validateEventFieldKey(key string) error {
	if key == "" {
		return fmt.Errorf("%w: field key must not be empty", ErrInvalidEvent)
	}

	if len(key) > maxEventFieldKeyLength {
		return fmt.Errorf(
			"%w: field key length %d exceeds maximum length %d",
			ErrInvalidEvent,
			len(key),
			maxEventFieldKeyLength,
		)
	}

	return validateEventKey("field key", key, ErrInvalidEvent)
}

// validateEventArtifacts validates artifacts and duplicate artifact IDs.
func validateEventArtifacts(artifacts []Artifact) error {
	seen := make(map[ArtifactID]struct{}, len(artifacts))

	for index, artifact := range artifacts {
		if err := artifact.Validate(); err != nil {
			return fmt.Errorf("%w: artifact %d: %w", ErrInvalidEvent, index, err)
		}

		if _, exists := seen[artifact.ID()]; exists {
			return fmt.Errorf("%w: duplicate artifact id %q", ErrInvalidEvent, artifact.ID())
		}

		seen[artifact.ID()] = struct{}{}
	}

	return nil
}

// validateEventLabels validates event labels.
func validateEventLabels(labels []string) error {
	seen := make(map[string]struct{}, len(labels))

	for index, label := range labels {
		if label == "" {
			return fmt.Errorf("%w: label %d must not be empty", ErrInvalidEvent, index)
		}

		if len(label) > maxEventLabelLength {
			return fmt.Errorf(
				"%w: label %d length %d exceeds maximum length %d",
				ErrInvalidEvent,
				index,
				len(label),
				maxEventLabelLength,
			)
		}

		if err := validateEventKey(fmt.Sprintf("label %d", index), label, ErrInvalidEvent); err != nil {
			return err
		}

		if _, exists := seen[label]; exists {
			return fmt.Errorf("%w: duplicate label %q", ErrInvalidEvent, label)
		}

		seen[label] = struct{}{}
	}

	return nil
}

// validateEventIDKey validates a dot-separated event ID.
//
// Event IDs use the same grammar as event keys, but also allow numeric-only
// suffix segments for generated stable IDs such as "command.started.001".
func validateEventIDKey(raw string) error {
	return validateEventKeyWithSegmentValidator("event id", raw, ErrInvalidEventID, validateEventIDSegment)
}

// validateEventKey validates a dot-separated machine-facing event key.
func validateEventKey(field string, raw string, sentinel error) error {
	return validateEventKeyWithSegmentValidator(field, raw, sentinel, validateCommandNameSegment)
}

// validateEventKeyWithSegmentValidator validates common event key boundaries
// and delegates per-segment grammar to validator.
func validateEventKeyWithSegmentValidator(
	field string,
	raw string,
	sentinel error,
	validator func(string) error,
) error {
	if raw == "" {
		return fmt.Errorf("%w: %s must not be empty", sentinel, field)
	}

	if strings.HasPrefix(raw, ".") {
		return fmt.Errorf("%w: %s must not start with %q", sentinel, field, ".")
	}

	if strings.HasSuffix(raw, ".") {
		return fmt.Errorf("%w: %s must not end with %q", sentinel, field, ".")
	}

	segments := strings.Split(raw, ".")
	for index, segment := range segments {
		if err := validator(segment); err != nil {
			return fmt.Errorf("%w: %s segment %d: %w", sentinel, field, index, err)
		}
	}

	return nil
}

// validateEventBlock validates compact possibly multi-line event text.
func validateEventBlock(field string, raw string, maxLength int) error {
	if raw == "" {
		return nil
	}

	if raw != normalizeEventBlock(raw) {
		return fmt.Errorf("%w: %s is not canonical", ErrInvalidEvent, field)
	}

	if err := textvalidate.ValidateCompactText(raw, maxLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", ErrInvalidEvent, field, err)
	}

	return nil
}

// normalizeEventBlock returns canonical event text.
func normalizeEventBlock(raw string) string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	raw = strings.TrimSpace(raw)

	if raw == "" {
		return ""
	}

	lines := strings.Split(raw, "\n")
	for index, line := range lines {
		lines[index] = strings.TrimSpace(line)
	}

	return strings.Join(lines, "\n")
}
