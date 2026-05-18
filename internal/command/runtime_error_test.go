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

// TestRuntimeExecuteBindingFailure verifies binding errors become failed results.
func TestRuntimeExecuteBindingFailure(t *testing.T) {
	t.Parallel()

	runtime := MustRuntime(RuntimeSpec{
		Binding: runtimeTestBinding(),
		Handler: runtimeTestOKHandler(),
		Clock:   FixedRuntimeClock{Time: runtimeTestTime()},
	})

	result, err := runtime.Execute(context.Background(), RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "xml"),
		},
		PositionalValues: []string{"stable"},
	})
	if err == nil {
		t.Fatalf("Execute() returned nil error")
	}

	if !errors.Is(err, ErrRuntimeExecution) {
		t.Fatalf("Execute() error = %v, want ErrRuntimeExecution", err)
	}

	if !result.IsFailed() {
		t.Fatalf("result IsFailed() = false, want true")
	}
}

// TestRuntimeExecuteCanceledContext verifies context cancellation behavior.
func TestRuntimeExecuteCanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	collector := &RuntimeEventCollector{}

	runtime := MustRuntime(RuntimeSpec{
		Binding:   runtimeTestBinding(),
		Clock:     FixedRuntimeClock{Time: runtimeTestTime()},
		EventSink: collector,
		Handler:   runtimeTestOKHandler(),
	})

	result, err := runtime.Execute(ctx, RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})
	if err == nil {
		t.Fatalf("Execute() returned nil error")
	}

	if !errors.Is(err, ErrRuntimeCanceled) {
		t.Fatalf("Execute() error = %v, want ErrRuntimeCanceled", err)
	}

	if !result.IsCanceled() {
		t.Fatalf("result IsCanceled() = false, want true")
	}

	if got, want := result.RecommendedExitCode(), 130; got != want {
		t.Fatalf("RecommendedExitCode() = %d, want %d", got, want)
	}

	assertEventKinds(t, collector.Events(), []EventKind{RuntimeEventCommandCompleted})
}

// TestRuntimeExecuteHandlerError verifies handler errors become failed results.
func TestRuntimeExecuteHandlerError(t *testing.T) {
	t.Parallel()

	handlerErr := errors.New("handler failed")

	runtime := MustRuntime(RuntimeSpec{
		Binding: runtimeTestBinding(),
		Clock:   FixedRuntimeClock{Time: runtimeTestTime()},
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			return FailedResult("handler failed"), handlerErr
		}),
	})

	result, err := runtime.Execute(context.Background(), RuntimeExecutionSpec{
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

	if !result.IsFailed() {
		t.Fatalf("result IsFailed() = false, want true")
	}
}

func TestRuntimeExecuteClassifiesCancellationDuringHandler(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	collector := &RuntimeEventCollector{}

	runtime := MustRuntime(RuntimeSpec{
		Binding:   runtimeTestBinding(),
		Clock:     FixedRuntimeClock{Time: runtimeTestTime()},
		EventSink: collector,
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			cancel()

			return OKResult("nominally ok"), nil
		}),
	})

	result, err := runtime.Execute(ctx, RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})
	if err == nil {
		t.Fatalf("Execute() returned nil error")
	}

	if !errors.Is(err, ErrRuntimeCanceled) {
		t.Fatalf("Execute() error = %v, want ErrRuntimeCanceled", err)
	}

	if errors.Is(err, ErrRuntimeExecution) {
		t.Fatalf("Execute() error = %v, should not wrap ErrRuntimeExecution", err)
	}

	if !result.IsCanceled() {
		t.Fatalf("result IsCanceled() = false, want true")
	}

	if !result.HasTiming() {
		t.Fatalf("result HasTiming() = false, want true")
	}

	assertEventKinds(t, collector.Events(), []EventKind{
		RuntimeEventCommandStarted,
		RuntimeEventBindingStarted,
		RuntimeEventBindingCompleted,
		RuntimeEventActionStarted,
		RuntimeEventActionCompleted,
		RuntimeEventCommandCompleted,
	})

	events := collector.Events()
	if got, want := events[4].Severity(), EventSeverityError; got != want {
		t.Fatalf("action.completed Severity() = %q, want %q", got, want)
	}

	if got, want := events[5].Severity(), EventSeverityError; got != want {
		t.Fatalf("command.completed Severity() = %q, want %q", got, want)
	}

	if events[4].HasResult() {
		t.Fatalf("action.completed carries full result payload")
	}

	completedResult, ok := events[5].Result()
	if !ok {
		t.Fatalf("command.completed result missing")
	}

	if !completedResult.IsCanceled() {
		t.Fatalf("command.completed result IsCanceled() = false, want true")
	}
}

func TestRuntimeExecuteClassifiesHandlerContextErrorAsCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	runtime := MustRuntime(RuntimeSpec{
		Binding: runtimeTestBinding(),
		Clock:   FixedRuntimeClock{Time: runtimeTestTime()},
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			cancel()

			return Result{}, ctx.Err()
		}),
	})

	result, err := runtime.Execute(ctx, RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})
	if err == nil {
		t.Fatalf("Execute() returned nil error")
	}

	if !errors.Is(err, ErrRuntimeCanceled) {
		t.Fatalf("Execute() error = %v, want ErrRuntimeCanceled", err)
	}

	if !result.IsCanceled() {
		t.Fatalf("result IsCanceled() = false, want true")
	}
}
