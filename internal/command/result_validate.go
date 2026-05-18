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

// Validate verifies result structural rules.
func (result Result) Validate() error {
	if err := result.status.Validate(); err != nil {
		return err
	}

	if result.message != "" {
		if err := validateResultBlock("message", result.message, maxResultMessageLength); err != nil {
			return err
		}
	}

	if !result.startedAt.IsZero() && !result.finishedAt.IsZero() && result.finishedAt.Before(result.startedAt) {
		return fmt.Errorf("%w: finished timestamp must not be before started timestamp", ErrInvalidResult)
	}

	if result.hasExit {
		if result.exitCode < 0 || result.exitCode > maxResultExitCode {
			return fmt.Errorf(
				"%w: exit code %d must be in range [0, %d]",
				ErrInvalidResult,
				result.exitCode,
				maxResultExitCode,
			)
		}
	}

	if err := validateResultArtifacts(result.artifacts); err != nil {
		return err
	}

	if err := validateResultWarnings(result.warnings); err != nil {
		return err
	}

	if err := validateResultFields(result.fields); err != nil {
		return err
	}

	if err := result.metadata.Validate(); err != nil {
		return fmt.Errorf("%w: invalid metadata: %w", ErrInvalidResult, err)
	}

	if err := result.visibility.Validate(); err != nil {
		return fmt.Errorf("%w: invalid visibility: %w", ErrInvalidResult, err)
	}

	return nil
}

// validateResultArtifacts validates artifact declarations and duplicate IDs.
func validateResultArtifacts(artifacts []Artifact) error {
	seen := make(map[ArtifactID]struct{}, len(artifacts))

	for index, artifact := range artifacts {
		if err := artifact.Validate(); err != nil {
			return fmt.Errorf("%w: artifact %d: %w", ErrInvalidResult, index, err)
		}

		if _, exists := seen[artifact.ID()]; exists {
			return fmt.Errorf("%w: duplicate artifact id %q", ErrInvalidResult, artifact.ID())
		}

		seen[artifact.ID()] = struct{}{}
	}

	return nil
}

// validateResultWarnings validates warnings.
func validateResultWarnings(warnings []ResultWarning) error {
	for index, warning := range warnings {
		if err := warning.Validate(); err != nil {
			return fmt.Errorf("%w: warning %d: %w", ErrInvalidResult, index, err)
		}
	}

	return nil
}

// validateResultFields validates result metadata fields.
func validateResultFields(fields map[string]string) error {
	for key, value := range fields {
		if err := validateResultKey("field key", key); err != nil {
			return err
		}

		if len(key) > maxResultFieldKeyLength {
			return fmt.Errorf(
				"%w: field key %q length %d exceeds maximum length %d",
				ErrInvalidResult,
				key,
				len(key),
				maxResultFieldKeyLength,
			)
		}

		if err := validateResultBlock("field "+key, value, maxResultFieldValueLength); err != nil {
			return err
		}
	}

	return nil
}

// validateResultKey validates a dot-separated machine-facing result key.
func validateResultKey(field string, raw string) error {
	if raw == "" {
		return fmt.Errorf("%w: %s must not be empty", ErrInvalidResult, field)
	}

	if strings.HasPrefix(raw, ".") {
		return fmt.Errorf("%w: %s must not start with %q", ErrInvalidResult, field, ".")
	}

	if strings.HasSuffix(raw, ".") {
		return fmt.Errorf("%w: %s must not end with %q", ErrInvalidResult, field, ".")
	}

	segments := strings.Split(raw, ".")
	for index, segment := range segments {
		if err := validateCommandNameSegment(segment); err != nil {
			return fmt.Errorf("%w: %s segment %d: %w", ErrInvalidResult, field, index, err)
		}
	}

	return nil
}

// validateResultBlock validates compact possibly multi-line result text.
func validateResultBlock(field string, raw string, maxLength int) error {
	if raw == "" {
		return nil
	}

	if raw != normalizeResultBlock(raw) {
		return fmt.Errorf("%w: %s is not canonical", ErrInvalidResult, field)
	}

	if err := textvalidate.ValidateCompactText(raw, maxLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", ErrInvalidResult, field, err)
	}

	return nil
}

// normalizeResultBlock returns canonical result text.
//
// It trims leading and trailing whitespace from the whole block and each line,
// normalizes CRLF/CR to LF, and preserves line boundaries.
func normalizeResultBlock(raw string) string {
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
