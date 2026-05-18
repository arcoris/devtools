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

// ActionWarning is a structured non-fatal warning produced by an action.
type ActionWarning struct {
	// Kind is a compact machine-facing warning kind.
	Kind string

	// Message is a required compact human-facing warning message.
	Message string
}

// NewActionWarning validates fields and returns an ActionWarning.
func NewActionWarning(kind string, message string) (ActionWarning, error) {
	warning := ActionWarning{
		Kind:    kind,
		Message: message,
	}

	if err := warning.Validate(); err != nil {
		return ActionWarning{}, err
	}

	return warning, nil
}

// MustActionWarning validates fields and returns an ActionWarning.
//
// MustActionWarning panics on invalid input. It is intended for tests and
// static command wiring.
func MustActionWarning(kind string, message string) ActionWarning {
	warning, err := NewActionWarning(kind, message)
	if err != nil {
		panic(err)
	}

	return warning
}

// IsZero reports whether warning has no fields set.
func (warning ActionWarning) IsZero() bool {
	return warning.Kind == "" && warning.Message == ""
}

// Validate verifies warning structural rules.
func (warning ActionWarning) Validate() error {
	if err := validateActionResultFieldKey("warning kind", warning.Kind); err != nil {
		return err
	}

	if strings.TrimSpace(warning.Message) == "" {
		return fmt.Errorf("%w: warning message must not be blank", ErrInvalidActionResult)
	}

	if err := validateActionResultText("warning message", warning.Message, maxActionMessageLength); err != nil {
		return err
	}

	return nil
}
