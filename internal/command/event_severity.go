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

// EventSeverity describes event importance.
type EventSeverity string

const (
	EventSeverityTrace   EventSeverity = "trace"
	EventSeverityDebug   EventSeverity = "debug"
	EventSeverityInfo    EventSeverity = "info"
	EventSeverityWarning EventSeverity = "warning"
	EventSeverityError   EventSeverity = "error"
)

// NewEventSeverity validates raw and returns it as an EventSeverity.
func NewEventSeverity(raw string) (EventSeverity, error) {
	severity := EventSeverity(strings.TrimSpace(raw))
	if err := severity.Validate(); err != nil {
		return "", err
	}

	return severity, nil
}

// MustEventSeverity validates raw and returns it as an EventSeverity.
//
// MustEventSeverity panics on invalid input.
func MustEventSeverity(raw string) EventSeverity {
	severity, err := NewEventSeverity(raw)
	if err != nil {
		panic(err)
	}

	return severity
}

// OrDefault returns EventSeverityInfo when severity is zero.
func (severity EventSeverity) OrDefault() EventSeverity {
	if severity == "" {
		return EventSeverityInfo
	}

	return severity
}

// String returns the canonical severity string.
func (severity EventSeverity) String() string {
	return string(severity)
}

// IsZero reports whether severity has not been set.
func (severity EventSeverity) IsZero() bool {
	return severity == ""
}

// IsKnown reports whether severity is supported.
func (severity EventSeverity) IsKnown() bool {
	switch severity {
	case EventSeverityTrace,
		EventSeverityDebug,
		EventSeverityInfo,
		EventSeverityWarning,
		EventSeverityError:
		return true
	default:
		return false
	}
}

// Validate verifies event severity structural rules.
func (severity EventSeverity) Validate() error {
	if severity == "" {
		return fmt.Errorf("%w: severity is empty", ErrInvalidEventSeverity)
	}

	if severity.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported severity %q", ErrInvalidEventSeverity, severity)
}
