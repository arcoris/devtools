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
	"errors"
	"time"
)

const (
	// maxEventIDLength is the maximum byte length of an event ID.
	maxEventIDLength = 255

	// maxEventKindLength is the maximum byte length of an event kind.
	maxEventKindLength = 255

	// maxEventMessageLength is the maximum byte length of a compact event
	// message.
	maxEventMessageLength = 4096

	// maxEventFieldKeyLength is the maximum byte length of one event field key.
	maxEventFieldKeyLength = 255

	// maxEventFieldValueLength is the maximum byte length of one event field
	// value.
	maxEventFieldValueLength = 4096

	// maxEventLabelLength is the maximum byte length of one event label.
	maxEventLabelLength = 255
)

var (
	// ErrEmptyEventID reports that an event ID was not provided.
	ErrEmptyEventID = errors.New("command event id is empty")

	// ErrInvalidEventID reports that an event ID violates the event ID grammar.
	ErrInvalidEventID = errors.New("command event id is invalid")

	// ErrEmptyEventKind reports that an event kind was not provided.
	ErrEmptyEventKind = errors.New("command event kind is empty")

	// ErrInvalidEventKind reports that an event kind violates the event kind
	// grammar.
	ErrInvalidEventKind = errors.New("command event kind is invalid")

	// ErrInvalidEventSeverity reports that an event severity is not supported.
	ErrInvalidEventSeverity = errors.New("command event severity is invalid")

	// ErrInvalidEvent reports that a command lifecycle event is malformed.
	ErrInvalidEvent = errors.New("command event is invalid")
)

// EventSpec describes a command lifecycle event before validation.
//
// EventSpec is a construction DTO. Event stores detached copies of mutable
// input state, so callers cannot mutate constructed events through shared
// slices, maps, or pointers.
//
// Event is an append-only lifecycle observation. It is not a logger, parser,
// executor, metrics exporter, or tracing backend. Logging, OpenTelemetry,
// JSONL output, terminal rendering, and persistence adapters can translate
// Event values into their own formats.
type EventSpec struct {
	// ID is an optional stable event identifier.
	//
	// Empty means the event sink may generate an ID or may store the event
	// without a stable ID. When present, ID must be a compact dot-separated key.
	ID string

	// Kind is the required event kind.
	//
	// Kind is open but validated. Built-in constants cover the common command
	// lifecycle phases, while projects may define more specific kinds.
	Kind EventKind

	// Severity is the event severity.
	//
	// Zero defaults to EventSeverityInfo.
	Severity EventSeverity

	// OccurredAt is the event timestamp.
	//
	// Zero defaults to time.Now().UTC().
	OccurredAt time.Time

	// CommandID optionally identifies the command associated with this event.
	//
	// A zero ID means the event is not attached to one command, or the command
	// identity is unavailable at this lifecycle stage.
	CommandID ID

	// Message is an optional compact human-facing event message.
	Message string

	// Fields contains optional machine-facing event fields.
	//
	// Field keys use a compact dot-separated key grammar. Values are compact
	// UTF-8 text.
	Fields map[string]string

	// Artifacts contains artifact references associated with this event.
	Artifacts []Artifact

	// Result optionally attaches a final or intermediate command result.
	//
	// Most events should not carry a Result. Result-bearing events are useful
	// for command.completed, action.completed, or result.produced phases.
	Result *Result

	// Labels contains optional machine-facing event labels.
	Labels []string

	// Metadata contains optional machine-facing event metadata.
	Metadata Metadata

	// Visibility controls default exposure in reports, logs, docs, and
	// discovery.
	//
	// A zero visibility defaults to public.
	Visibility Visibility
}

// Event is a validated framework-neutral command lifecycle event.
//
// Event is immutable-style:
//
//   - constructors normalize defaults and copy mutable input state;
//   - accessors return detached copies;
//   - With* methods return validated copies;
//   - callers cannot mutate internal state through returned values.
type Event struct {
	id         EventID
	hasID      bool
	kind       EventKind
	severity   EventSeverity
	occurredAt time.Time
	commandID  ID
	message    string
	fields     map[string]string
	artifacts  []Artifact
	result     *Result
	labels     []string
	metadata   Metadata
	visibility Visibility
}
