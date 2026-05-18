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
	"sort"
)

// ActionResult is the structured non-error result of action execution.
//
// ActionResult is intentionally small and adapter-neutral. It can be converted
// into plain text, JSON, GitHub summaries, reports, or test assertions by
// higher layers.
//
// A zero ActionResult normalizes to an OK result.
type ActionResult struct {
	// Status is the high-level action status.
	Status ActionStatus

	// Message is an optional compact human-facing summary.
	Message string

	// Data is optional adapter-neutral structured data.
	//
	// The command kernel does not inspect Data. Serialization layers are
	// responsible for supporting or rejecting the concrete value.
	Data any

	// Artifacts contains optional artifact references produced by the action.
	Artifacts []ActionArtifact

	// Warnings contains non-fatal warnings produced by the action.
	Warnings []ActionWarning

	// Fields contains optional machine-facing metadata for the result.
	Fields map[string]string
}

// NewActionResult normalizes result, validates it, and returns the normalized
// value.
func NewActionResult(result ActionResult) (ActionResult, error) {
	normalized := result.Normalize()
	if err := normalized.Validate(); err != nil {
		return ActionResult{}, err
	}

	return normalized, nil
}

// MustActionResult normalizes result and panics on invalid input.
func MustActionResult(result ActionResult) ActionResult {
	normalized, err := NewActionResult(result)
	if err != nil {
		panic(err)
	}

	return normalized
}

// Normalize returns a copy with default values applied and mutable collections
// detached.
func (result ActionResult) Normalize() ActionResult {
	return ActionResult{
		Status:    result.Status.OrDefault(),
		Message:   result.Message,
		Data:      result.Data,
		Artifacts: cloneActionArtifacts(result.Artifacts),
		Warnings:  cloneActionWarnings(result.Warnings),
		Fields:    cloneActionStringMap(result.Fields),
	}
}

// Validate verifies action result structural rules.
func (result ActionResult) Validate() error {
	normalized := result.Normalize()

	if err := normalized.Status.Validate(); err != nil {
		return err
	}

	if normalized.Message != "" {
		if err := validateActionResultText("message", normalized.Message, maxActionMessageLength); err != nil {
			return err
		}
	}

	for index, artifact := range normalized.Artifacts {
		if err := artifact.Validate(); err != nil {
			return fmt.Errorf("%w: artifact %d: %w", ErrInvalidActionResult, index, err)
		}
	}

	for index, warning := range normalized.Warnings {
		if err := warning.Validate(); err != nil {
			return fmt.Errorf("%w: warning %d: %w", ErrInvalidActionResult, index, err)
		}
	}

	return validateActionResultFields(normalized.Fields)
}

// IsOK reports whether the normalized result status is OK.
func (result ActionResult) IsOK() bool {
	return result.Status.OrDefault() == ActionStatusOK
}

// IsSkipped reports whether the normalized result status is skipped.
func (result ActionResult) IsSkipped() bool {
	return result.Status.OrDefault() == ActionStatusSkipped
}

// IsFailed reports whether the normalized result status is failed.
func (result ActionResult) IsFailed() bool {
	return result.Status.OrDefault() == ActionStatusFailed
}

// HasData reports whether result contains adapter-neutral structured data.
func (result ActionResult) HasData() bool {
	return result.Data != nil
}

// ArtifactRefs returns a detached copy of artifact references.
func (result ActionResult) ArtifactRefs() []ActionArtifact {
	return cloneActionArtifacts(result.Artifacts)
}

// HasArtifacts reports whether result contains artifact references.
func (result ActionResult) HasArtifacts() bool {
	return len(result.Artifacts) > 0
}

// ArtifactCount returns the number of artifact references.
func (result ActionResult) ArtifactCount() int {
	return len(result.Artifacts)
}

// WarningRefs returns a detached copy of warnings.
func (result ActionResult) WarningRefs() []ActionWarning {
	return cloneActionWarnings(result.Warnings)
}

// HasWarnings reports whether result contains warnings.
func (result ActionResult) HasWarnings() bool {
	return len(result.Warnings) > 0
}

// WarningCount returns the number of warnings.
func (result ActionResult) WarningCount() int {
	return len(result.Warnings)
}

// Field returns a result metadata field and whether it exists.
func (result ActionResult) Field(key string) (string, bool) {
	value, ok := result.Fields[key]

	return value, ok
}

// HasField reports whether a result metadata field exists.
func (result ActionResult) HasField(key string) bool {
	_, ok := result.Field(key)

	return ok
}

// HasFields reports whether result contains metadata fields.
func (result ActionResult) HasFields() bool {
	return len(result.Fields) > 0
}

// FieldCount returns the number of result metadata fields.
func (result ActionResult) FieldCount() int {
	return len(result.Fields)
}

// FieldMap returns a detached copy of result metadata fields.
func (result ActionResult) FieldMap() map[string]string {
	return cloneActionStringMap(result.Fields)
}

// FieldKeys returns result metadata keys in deterministic lexical order.
func (result ActionResult) FieldKeys() []string {
	keys := make([]string, 0, len(result.Fields))
	for key := range result.Fields {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
