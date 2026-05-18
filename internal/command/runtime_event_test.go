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
	"errors"
	"sync"
	"testing"
)

// TestRuntimeEventCollector verifies in-memory event collection.
func TestRuntimeEventCollector(t *testing.T) {
	t.Parallel()

	collector := &RuntimeEventCollector{}
	event := MustSimpleEvent(EventKindDiagnostic, "diagnostic")

	if err := collector.RecordEvent(context.Background(), event); err != nil {
		t.Fatalf("RecordEvent() returned unexpected error: %v", err)
	}

	if got, want := collector.Len(), 1; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}

	events := collector.Events()
	events[0] = MustSimpleEvent(EventKindWarning, "changed")

	if got, want := collector.Events()[0].Kind(), EventKindDiagnostic; got != want {
		t.Fatalf("Events() returned mutable state: got %q, want %q", got, want)
	}

	collector.Reset()

	if got, want := collector.Len(), 0; got != want {
		t.Fatalf("Len() after Reset() = %d, want %d", got, want)
	}
}

// TestRuntimeEventCollectorNil verifies nil collector behavior is deterministic.
func TestRuntimeEventCollectorNil(t *testing.T) {
	t.Parallel()

	var collector *RuntimeEventCollector
	event := MustSimpleEvent(EventKindDiagnostic, "diagnostic")

	if err := collector.RecordEvent(context.Background(), event); err == nil {
		t.Fatalf("RecordEvent(nil collector) returned nil error")
	} else if !errors.Is(err, ErrInvalidRuntime) {
		t.Fatalf("RecordEvent(nil collector) error = %v, want ErrInvalidRuntime", err)
	}

	if got := collector.Events(); got != nil {
		t.Fatalf("Events(nil collector) = %#v, want nil", got)
	}

	if got := collector.Len(); got != 0 {
		t.Fatalf("Len(nil collector) = %d, want 0", got)
	}

	collector.Reset()
}

// TestRuntimeEventCollectorConcurrentRecord verifies the collector can be used
// by concurrent runtime observers without racing.
func TestRuntimeEventCollectorConcurrentRecord(t *testing.T) {
	t.Parallel()

	collector := &RuntimeEventCollector{}
	event := MustSimpleEvent(EventKindDiagnostic, "diagnostic")

	const goroutines = 16
	const perGoroutine = 32

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < perGoroutine; j++ {
				if err := collector.RecordEvent(context.Background(), event); err != nil {
					t.Errorf("RecordEvent() returned unexpected error: %v", err)
					return
				}
			}
		}()
	}

	wg.Wait()

	if got, want := collector.Len(), goroutines*perGoroutine; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}

	collector.Reset()

	if got, want := collector.Len(), 0; got != want {
		t.Fatalf("Len() after Reset() = %d, want %d", got, want)
	}
}

func TestRuntimeEventCollectorRejectsCanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	collector := &RuntimeEventCollector{}

	err := collector.RecordEvent(ctx, MustSimpleEvent(EventKindDiagnostic, "diagnostic"))
	if err == nil {
		t.Fatalf("RecordEvent() returned nil error")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("RecordEvent() error = %v, want context.Canceled", err)
	}

	if got, want := collector.Len(), 0; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}
}

// TestRuntimeEventSinkFuncNil verifies nil event sink function behavior.
func TestRuntimeEventSinkFuncNil(t *testing.T) {
	t.Parallel()

	var fn RuntimeEventSinkFunc

	err := fn.RecordEvent(context.Background(), MustSimpleEvent(EventKindDiagnostic, "diagnostic"))
	if err == nil {
		t.Fatalf("RecordEvent() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRuntime) {
		t.Fatalf("RecordEvent() error = %v, want ErrInvalidRuntime", err)
	}
}
