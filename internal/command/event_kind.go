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

// EventKind is an open validated command lifecycle event kind.
type EventKind string

const (
	EventKindCommandStarted   EventKind = "command.started"
	EventKindCommandCompleted EventKind = "command.completed"
	EventKindBindingStarted   EventKind = "binding.started"
	EventKindBindingCompleted EventKind = "binding.completed"
	EventKindActionStarted    EventKind = "action.started"
	EventKindActionCompleted  EventKind = "action.completed"
	EventKindArtifactProduced EventKind = "artifact.produced"
	EventKindResultProduced   EventKind = "result.produced"
	EventKindDiagnostic       EventKind = "diagnostic"
	EventKindWarning          EventKind = "warning"
)

// NewEventKind validates raw and returns it as an EventKind.
func NewEventKind(raw string) (EventKind, error) {
	kind := EventKind(strings.TrimSpace(raw))
	if err := kind.Validate(); err != nil {
		return "", err
	}

	return kind, nil
}

// MustEventKind validates raw and returns it as an EventKind.
//
// MustEventKind panics on invalid input.
func MustEventKind(raw string) EventKind {
	kind, err := NewEventKind(raw)
	if err != nil {
		panic(err)
	}

	return kind
}

// String returns the canonical event kind string.
func (kind EventKind) String() string {
	return string(kind)
}

// IsZero reports whether kind has not been set.
func (kind EventKind) IsZero() bool {
	return kind == ""
}

// Validate verifies event kind structural rules.
func (kind EventKind) Validate() error {
	raw := string(kind)
	if raw == "" {
		return ErrEmptyEventKind
	}

	if len(raw) > maxEventKindLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidEventKind,
			len(raw),
			maxEventKindLength,
		)
	}

	if err := validateEventKey("event kind", raw, ErrInvalidEventKind); err != nil {
		return err
	}

	return nil
}
