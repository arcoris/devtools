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
	"testing"
)

// TestRuntimeEventSinkError verifies event sink errors stop execution.
func TestRuntimeEventSinkError(t *testing.T) {
	t.Parallel()

	sinkErr := errors.New("sink failed")

	runtime := MustRuntime(RuntimeSpec{
		Binding: runtimeTestBinding(),
		Clock:   FixedRuntimeClock{Time: runtimeTestTime()},
		EventSink: RuntimeEventSinkFunc(func(ctx context.Context, event Event) error {
			return sinkErr
		}),
		Handler: runtimeTestOKHandler(),
	})

	_, err := runtime.Execute(context.Background(), RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})
	if err == nil {
		t.Fatalf("Execute() returned nil error")
	}

	if !errors.Is(err, ErrRuntimeExecution) {
		t.Fatalf("Execute() error = %v, want ErrRuntimeExecution", err)
	}
}

// TestRuntimeSuppressEvents verifies event suppression option.
func TestRuntimeSuppressEvents(t *testing.T) {
	t.Parallel()

	collector := &RuntimeEventCollector{}

	runtime := MustRuntime(RuntimeSpec{
		Binding:   runtimeTestBinding(),
		Clock:     FixedRuntimeClock{Time: runtimeTestTime()},
		EventSink: collector,
		Options:   DefaultRuntimeOptions().MustWithSuppressEvents(true),
		Handler:   runtimeTestOKHandler(),
	})

	_, err := runtime.Execute(context.Background(), RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	if got, want := collector.Len(), 0; got != want {
		t.Fatalf("collector Len() = %d, want %d", got, want)
	}
}

// TestRuntimeEventResultPayloadPolicy verifies full Result payloads are only
// attached to final lifecycle events.
func TestRuntimeEventResultPayloadPolicy(t *testing.T) {
	t.Parallel()

	collector := &RuntimeEventCollector{}

	runtime := MustRuntime(RuntimeSpec{
		Binding:   runtimeTestBinding(),
		Clock:     FixedRuntimeClock{Time: runtimeTestTime()},
		EventSink: collector,
		Handler:   runtimeTestOKHandler(),
	})

	_, err := runtime.Execute(context.Background(), RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	events := collector.Events()
	if got, want := events[4].Kind(), RuntimeEventActionCompleted; got != want {
		t.Fatalf("event 4 kind = %q, want %q", got, want)
	}

	if events[4].HasResult() {
		t.Fatalf("action.completed carries full result payload")
	}

	if got, ok := events[4].Field("result.status"); !ok || got != ResultStatusOK.String() {
		t.Fatalf("action.completed result.status = %q, %v; want ok", got, ok)
	}

	if got, want := events[5].Kind(), RuntimeEventCommandCompleted; got != want {
		t.Fatalf("event 5 kind = %q, want %q", got, want)
	}

	if got, ok := events[5].Field("result.status"); !ok || got != ResultStatusOK.String() {
		t.Fatalf("command.completed result.status = %q, %v; want ok", got, ok)
	}

	result, ok := events[5].Result()
	if !ok {
		t.Fatalf("command.completed result missing")
	}

	if !result.IsOK() {
		t.Fatalf("command.completed result IsOK() = false, want true")
	}
}
