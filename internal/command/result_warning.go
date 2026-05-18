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
)

// ResultWarningSpec describes a non-fatal result warning before validation.
type ResultWarningSpec struct {
	// Kind is a compact machine-facing warning kind.
	//
	// Examples: "partial", "unstable", "missing-tool", "deprecated",
	// "best-effort".
	Kind string

	// Message is a compact human-facing warning message.
	Message string

	// Hints contains optional remediation hints.
	Hints []string
}

// ResultWarning is a validated non-fatal warning collected during command
// execution.
type ResultWarning struct {
	kind    string
	message string
	hints   []string
}

// NewResultWarning validates spec and returns a ResultWarning.
func NewResultWarning(spec ResultWarningSpec) (ResultWarning, error) {
	warning := ResultWarning{
		kind:    strings.TrimSpace(spec.Kind),
		message: normalizeResultBlock(spec.Message),
		hints:   normalizeResultHints(spec.Hints),
	}

	if err := warning.Validate(); err != nil {
		return ResultWarning{}, err
	}

	return warning, nil
}

// MustResultWarning validates spec and returns a ResultWarning.
//
// MustResultWarning panics on invalid input. It is intended for tests and
// controlled static wiring.
func MustResultWarning(spec ResultWarningSpec) ResultWarning {
	warning, err := NewResultWarning(spec)
	if err != nil {
		panic(err)
	}

	return warning
}

// Kind returns the machine-facing warning kind.
func (warning ResultWarning) Kind() string {
	return warning.kind
}

// Message returns the human-facing warning message.
func (warning ResultWarning) Message() string {
	return warning.message
}

// Hints returns detached remediation hints.
func (warning ResultWarning) Hints() []string {
	return cloneResultStrings(warning.hints)
}

// HasHints reports whether warning hints are present.
func (warning ResultWarning) HasHints() bool {
	return len(warning.hints) > 0
}

// Validate verifies warning structural rules.
func (warning ResultWarning) Validate() error {
	if err := validateResultKey("warning kind", warning.kind); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidResultWarning, err)
	}

	if warning.message == "" {
		return fmt.Errorf("%w: message must not be empty", ErrInvalidResultWarning)
	}

	if err := validateResultBlock("warning message", warning.message, maxResultWarningMessageLength); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidResultWarning, err)
	}

	seen := make(map[string]struct{}, len(warning.hints))
	for index, hint := range warning.hints {
		if err := validateResultBlock(fmt.Sprintf("hint %d", index), hint, maxResultHintLength); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidResultWarning, err)
		}

		if _, exists := seen[hint]; exists {
			return fmt.Errorf("%w: duplicate hint %q", ErrInvalidResultWarning, hint)
		}

		seen[hint] = struct{}{}
	}

	return nil
}

// WithHint returns a validated copy with one hint appended.
func (warning ResultWarning) WithHint(hint string) (ResultWarning, error) {
	spec := warning.spec()
	spec.Hints = append(spec.Hints, hint)

	return NewResultWarning(spec)
}

// MustWithHint returns a validated copy with one hint appended and panics on
// invalid input.
func (warning ResultWarning) MustWithHint(hint string) ResultWarning {
	next, err := warning.WithHint(hint)
	if err != nil {
		panic(err)
	}

	return next
}

// spec returns a detached construction spec.
func (warning ResultWarning) spec() ResultWarningSpec {
	return ResultWarningSpec{
		Kind:    warning.kind,
		Message: warning.message,
		Hints:   cloneResultStrings(warning.hints),
	}
}
