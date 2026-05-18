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

// WithID returns a validated copy with ID set.
func (event Event) WithID(id string) (Event, error) {
	spec := event.spec()
	spec.ID = id

	return NewEvent(spec)
}

// MustWithID returns a validated copy with ID set and panics on invalid input.
func (event Event) MustWithID(id string) Event {
	next, err := event.WithID(id)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutID returns a validated copy without stable ID.
func (event Event) WithoutID() Event {
	spec := event.spec()
	spec.ID = ""

	return MustEvent(spec)
}

// WithKind returns a validated copy with Kind replaced.
func (event Event) WithKind(kind EventKind) (Event, error) {
	spec := event.spec()
	spec.Kind = kind

	return NewEvent(spec)
}

// MustWithKind returns a validated copy with Kind replaced and panics on
// invalid input.
func (event Event) MustWithKind(kind EventKind) Event {
	next, err := event.WithKind(kind)
	if err != nil {
		panic(err)
	}

	return next
}

// WithSeverity returns a validated copy with Severity replaced.
func (event Event) WithSeverity(severity EventSeverity) (Event, error) {
	spec := event.spec()
	spec.Severity = severity

	return NewEvent(spec)
}

// MustWithSeverity returns a validated copy with Severity replaced and panics on
// invalid input.
func (event Event) MustWithSeverity(severity EventSeverity) Event {
	next, err := event.WithSeverity(severity)
	if err != nil {
		panic(err)
	}

	return next
}

// WithMessage returns a validated copy with Message replaced.
func (event Event) WithMessage(message string) (Event, error) {
	spec := event.spec()
	spec.Message = message

	return NewEvent(spec)
}

// MustWithMessage returns a validated copy with Message replaced and panics on
// invalid input.
func (event Event) MustWithMessage(message string) Event {
	next, err := event.WithMessage(message)
	if err != nil {
		panic(err)
	}

	return next
}
