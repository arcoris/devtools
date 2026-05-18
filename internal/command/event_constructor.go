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
	"strings"
	"time"
)

// NewEvent validates spec and returns an Event.
func NewEvent(spec EventSpec) (Event, error) {
	var eventID EventID
	var hasID bool

	if strings.TrimSpace(spec.ID) != "" {
		id, err := NewEventID(spec.ID)
		if err != nil {
			return Event{}, err
		}

		eventID = id
		hasID = true
	}

	severity := spec.Severity.OrDefault()
	occurredAt := spec.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}

	visibility := spec.Visibility.OrDefault()

	var result *Result
	if spec.Result != nil {
		value := *spec.Result
		result = &value
	}

	event := Event{
		id:         eventID,
		hasID:      hasID,
		kind:       spec.Kind,
		severity:   severity,
		occurredAt: occurredAt,
		commandID:  spec.CommandID,
		message:    normalizeEventBlock(spec.Message),
		fields:     cloneEventStringMap(spec.Fields),
		artifacts:  cloneEventArtifacts(spec.Artifacts),
		result:     result,
		labels:     cloneEventStrings(spec.Labels),
		metadata:   spec.Metadata,
		visibility: visibility,
	}

	if err := event.Validate(); err != nil {
		return Event{}, err
	}

	return event, nil
}

// MustEvent validates spec and returns an Event.
//
// MustEvent panics on invalid input. It is intended for tests and controlled
// static wiring.
func MustEvent(spec EventSpec) Event {
	event, err := NewEvent(spec)
	if err != nil {
		panic(err)
	}

	return event
}

// NewSimpleEvent returns an event with a kind and optional message.
func NewSimpleEvent(kind EventKind, message string) (Event, error) {
	return NewEvent(EventSpec{
		Kind:    kind,
		Message: message,
	})
}

// MustSimpleEvent returns an event with a kind and optional message and panics
// on invalid input.
func MustSimpleEvent(kind EventKind, message string) Event {
	return MustEvent(EventSpec{
		Kind:    kind,
		Message: message,
	})
}
