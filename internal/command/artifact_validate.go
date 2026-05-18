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

// Validate verifies artifact structural rules.
func (artifact Artifact) Validate() error {
	if err := artifact.id.Validate(); err != nil {
		return err
	}

	if err := artifact.kind.Validate(); err != nil {
		return err
	}

	if err := validateArtifactLocation(artifact.location); err != nil {
		return err
	}

	if !artifact.format.IsZero() {
		if err := artifact.format.Validate(); err != nil {
			return err
		}
	}

	if artifact.mediaType != "" {
		if err := validateArtifactMediaType(artifact.mediaType); err != nil {
			return err
		}
	}

	if artifact.description != "" {
		if err := validateArtifactBlock("description", artifact.description, maxArtifactDescriptionLength); err != nil {
			return err
		}
	}

	if artifact.hasSize && artifact.sizeBytes < 0 {
		return fmt.Errorf("%w: size must be non-negative", ErrInvalidArtifact)
	}

	if artifact.digest != nil {
		if err := artifact.digest.Validate(); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidArtifact, err)
		}
	}

	if err := validateArtifactLabels(artifact.labels); err != nil {
		return err
	}

	if err := artifact.metadata.Validate(); err != nil {
		return fmt.Errorf("%w: invalid metadata: %w", ErrInvalidArtifact, err)
	}

	if err := artifact.visibility.Validate(); err != nil {
		return fmt.Errorf("%w: invalid visibility: %w", ErrInvalidArtifact, err)
	}

	return nil
}

// validateArtifactLocation validates artifact location text.
func validateArtifactLocation(raw string) error {
	if raw == "" {
		return ErrEmptyArtifactLocation
	}

	if strings.TrimSpace(raw) == "" {
		return ErrEmptyArtifactLocation
	}

	if raw != normalizeArtifactText(raw) {
		return fmt.Errorf("%w: location is not canonical", ErrInvalidArtifactLocation)
	}

	if err := validateArtifactBlock("location", raw, maxArtifactLocationLength); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArtifactLocation, err)
	}

	return nil
}

// validateArtifactMediaType validates a compact media type.
//
// This is intentionally lightweight. It rejects whitespace and requires a
// "type/subtype" shape. It does not implement the full RFC grammar.
func validateArtifactMediaType(raw string) error {
	if raw == "" {
		return nil
	}

	if raw != normalizeArtifactText(raw) {
		return fmt.Errorf("%w: media type is not canonical", ErrInvalidArtifact)
	}

	if err := validateArtifactText("media type", raw, maxArtifactMediaTypeLength); err != nil {
		return err
	}

	parts := strings.Split(raw, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("%w: media type must have type/subtype shape", ErrInvalidArtifact)
	}

	if strings.ContainsAny(raw, " \t\r\n") {
		return fmt.Errorf("%w: media type must not contain whitespace", ErrInvalidArtifact)
	}

	return nil
}

// validateArtifactLabels validates label grammar and duplicate labels.
func validateArtifactLabels(labels []string) error {
	seen := make(map[string]struct{}, len(labels))

	for index, label := range labels {
		if label == "" {
			return fmt.Errorf("%w: label %d must not be empty", ErrInvalidArtifact, index)
		}

		if len(label) > maxArtifactLabelLength {
			return fmt.Errorf(
				"%w: label %d length %d exceeds maximum length %d",
				ErrInvalidArtifact,
				index,
				len(label),
				maxArtifactLabelLength,
			)
		}

		if err := validateArtifactKey(label, ErrInvalidArtifact, fmt.Sprintf("label %d", index)); err != nil {
			return err
		}

		if _, exists := seen[label]; exists {
			return fmt.Errorf("%w: duplicate label %q", ErrInvalidArtifact, label)
		}

		seen[label] = struct{}{}
	}

	return nil
}

// validateArtifactKey validates a dot-separated machine-facing artifact key.
func validateArtifactKey(raw string, sentinel error, field string) error {
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
		if err := validateCommandNameSegment(segment); err != nil {
			return fmt.Errorf("%w: %s segment %d: %w", sentinel, field, index, err)
		}
	}

	return nil
}

// validateArtifactText validates compact single-line artifact text.
func validateArtifactText(field string, raw string, maxLength int) error {
	if err := textvalidate.ValidateSingleLineText(raw, maxLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", ErrInvalidArtifact, field, err)
	}

	return nil
}

// validateArtifactBlock validates possibly multi-line artifact text.
func validateArtifactBlock(field string, raw string, maxLength int) error {
	if err := textvalidate.ValidateCompactText(raw, maxLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", ErrInvalidArtifact, field, err)
	}

	return nil
}

// normalizeArtifactText returns canonical single-line artifact text.
func normalizeArtifactText(raw string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(raw)), " ")
}

// normalizeArtifactBlock returns canonical block artifact text.
func normalizeArtifactBlock(raw string) string {
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
