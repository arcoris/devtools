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
	"context"
	"fmt"
	"time"
)

// emit builds and records one lifecycle event.
func (runtime Runtime) emit(ctx context.Context, kind EventKind, occurredAt time.Time, result *Result, cause error) error {
	if runtime.eventSink == nil || runtime.options.SuppressEvents() {
		return nil
	}

	severity := EventSeverityInfo
	if cause != nil {
		severity = EventSeverityError
	}

	fields := map[string]string{
		"runtime": runtime.name,
	}

	if cause != nil {
		fields["error"] = cause.Error()
	}

	event, err := NewEvent(EventSpec{
		Kind:       kind,
		Severity:   severity,
		OccurredAt: occurredAt,
		CommandID:  runtime.commandID,
		Message:    runtime.eventMessage(kind, cause),
		Fields:     fields,
		Result:     result,
		Metadata:   runtime.metadata,
		Visibility: runtime.visibility,
	})
	if err != nil {
		return fmt.Errorf("%w: build event %q: %w", ErrRuntimeExecution, kind, err)
	}

	if err := runtime.eventSink.RecordEvent(ctx, event); err != nil {
		return fmt.Errorf("%w: record event %q: %w", ErrRuntimeExecution, kind, err)
	}

	return nil
}

// eventMessage returns a compact lifecycle event message.
func (runtime Runtime) eventMessage(kind EventKind, cause error) string {
	if cause != nil {
		return fmt.Sprintf("%s failed.", kind)
	}

	return fmt.Sprintf("%s.", kind)
}

// RuntimeEventSink records lifecycle events.
type RuntimeEventSink interface {
	RecordEvent(ctx context.Context, event Event) error
}

// RuntimeEventSinkFunc adapts a function to RuntimeEventSink.
type RuntimeEventSinkFunc func(ctx context.Context, event Event) error

// RecordEvent records event through fn.
func (fn RuntimeEventSinkFunc) RecordEvent(ctx context.Context, event Event) error {
	if fn == nil {
		return fmt.Errorf("%w: nil runtime event sink function", ErrInvalidRuntime)
	}

	return fn(ctx, event)
}

// RuntimeEventCollector is an in-memory RuntimeEventSink useful for tests and
// embedding scenarios where the caller wants to inspect emitted events.
type RuntimeEventCollector struct {
	events []Event
}

// RecordEvent stores event in memory.
func (collector *RuntimeEventCollector) RecordEvent(ctx context.Context, event Event) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if collector == nil {
		return fmt.Errorf("%w: nil runtime event collector", ErrInvalidRuntime)
	}

	if err := event.Validate(); err != nil {
		return err
	}

	collector.events = append(collector.events, event)

	return nil
}

// Events returns detached collected events.
func (collector *RuntimeEventCollector) Events() []Event {
	if collector == nil {
		return nil
	}

	out := make([]Event, len(collector.events))
	copy(out, collector.events)

	return out
}

// Len returns the number of collected events.
func (collector *RuntimeEventCollector) Len() int {
	if collector == nil {
		return 0
	}

	return len(collector.events)
}

// Reset removes collected events.
func (collector *RuntimeEventCollector) Reset() {
	if collector == nil {
		return
	}

	collector.events = nil
}

const (
	RuntimeEventCommandStarted   EventKind = EventKindCommandStarted
	RuntimeEventCommandCompleted EventKind = EventKindCommandCompleted
	RuntimeEventBindingStarted   EventKind = EventKindBindingStarted
	RuntimeEventBindingCompleted EventKind = EventKindBindingCompleted
	RuntimeEventActionStarted    EventKind = EventKindActionStarted
	RuntimeEventActionCompleted  EventKind = EventKindActionCompleted
)
