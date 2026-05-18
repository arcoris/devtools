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

// EventID is an optional stable machine-facing lifecycle event identifier.
type EventID string

// NewEventID validates raw and returns it as an EventID.
func NewEventID(raw string) (EventID, error) {
	id := EventID(strings.TrimSpace(raw))
	if err := id.Validate(); err != nil {
		return "", err
	}

	return id, nil
}

// MustEventID validates raw and returns it as an EventID.
//
// MustEventID panics on invalid input. It is intended for static event
// declarations and tests.
func MustEventID(raw string) EventID {
	id, err := NewEventID(raw)
	if err != nil {
		panic(err)
	}

	return id
}

// String returns the canonical event ID string.
func (id EventID) String() string {
	return string(id)
}

// IsZero reports whether ID is absent.
func (id EventID) IsZero() bool {
	return id == ""
}

// IsValid reports whether ID satisfies the event ID grammar.
func (id EventID) IsValid() bool {
	return id.Validate() == nil
}

// Validate verifies event ID structural rules.
func (id EventID) Validate() error {
	raw := string(id)
	if raw == "" {
		return ErrEmptyEventID
	}

	if len(raw) > maxEventIDLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidEventID,
			len(raw),
			maxEventIDLength,
		)
	}

	if err := validateEventIDKey(raw); err != nil {
		return err
	}

	return nil
}

// validateEventIDSegment validates one event ID segment.
func validateEventIDSegment(segment string) error {
	if isASCIIDigitString(segment) {
		return nil
	}

	return validateCommandNameSegment(segment)
}

// isASCIIDigitString reports whether raw is non-empty and contains only ASCII
// digits.
func isASCIIDigitString(raw string) bool {
	if raw == "" {
		return false
	}

	for _, r := range raw {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}
