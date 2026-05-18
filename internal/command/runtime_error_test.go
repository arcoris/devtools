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
	"fmt"
	"testing"
	"time"
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

func TestRuntimeExecuteDeadlineExceededBeforeStart(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

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

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Execute() error = %v, want context.DeadlineExceeded", err)
	}

	if !result.IsCanceled() {
		t.Fatalf("result IsCanceled() = false, want true")
	}

	events := collector.Events()
	assertEventKinds(t, events, []EventKind{RuntimeEventCommandCompleted})

	if got, ok := events[0].Field("result.status"); !ok || got != ResultStatusCanceled.String() {
		t.Fatalf("command.completed result.status = %q, %v; want canceled", got, ok)
	}

	if _, ok := events[0].Result(); !ok {
		t.Fatalf("command.completed result missing")
	}
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

func TestRuntimeExecuteClassifiesWrappedContextCanceledAsCanceled(t *testing.T) {
	t.Parallel()

	runtime := MustRuntime(RuntimeSpec{
		Binding: runtimeTestBinding(),
		Clock:   FixedRuntimeClock{Time: runtimeTestTime()},
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			return Result{}, fmt.Errorf("handler stopped: %w", context.Canceled)
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

	if !errors.Is(err, ErrRuntimeCanceled) {
		t.Fatalf("Execute() error = %v, want ErrRuntimeCanceled", err)
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Execute() error = %v, want context.Canceled", err)
	}

	if errors.Is(err, ErrRuntimeExecution) {
		t.Fatalf("Execute() error = %v, should not wrap ErrRuntimeExecution", err)
	}

	if !result.IsCanceled() {
		t.Fatalf("result IsCanceled() = false, want true")
	}
}

func TestRuntimeExecuteClassifiesHandlerDeadlineExceededAsCanceled(t *testing.T) {
	t.Parallel()

	runtime := MustRuntime(RuntimeSpec{
		Binding: runtimeTestBinding(),
		Clock:   FixedRuntimeClock{Time: runtimeTestTime()},
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			return Result{}, context.DeadlineExceeded
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

	if !errors.Is(err, ErrRuntimeCanceled) {
		t.Fatalf("Execute() error = %v, want ErrRuntimeCanceled", err)
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Execute() error = %v, want context.DeadlineExceeded", err)
	}

	if !result.IsCanceled() {
		t.Fatalf("result IsCanceled() = false, want true")
	}
}

func TestRuntimeExecuteCanceledPartialResultPreservesArtifactsAndWarnings(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	collector := &RuntimeEventCollector{}

	runtime := MustRuntime(RuntimeSpec{
		Binding:   runtimeTestBinding(),
		Clock:     FixedRuntimeClock{Time: runtimeTestTime()},
		EventSink: collector,
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			cancel()

			return MustResult(ResultSpec{
				Status:  ResultStatusOK,
				Message: "partial",
				Artifacts: []Artifact{
					resultTestArtifact("bench.report"),
				},
				Warnings: []ResultWarning{
					MustResultWarning(ResultWarningSpec{
						Kind:    "partial",
						Message: "partial output is available",
					}),
				},
			}), nil
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

	if got, want := len(result.Artifacts()), 1; got != want {
		t.Fatalf("len(Artifacts()) = %d, want %d", got, want)
	}

	if got, want := result.WarningCount(), 1; got != want {
		t.Fatalf("WarningCount() = %d, want %d", got, want)
	}

	events := collector.Events()
	completedResult, ok := events[len(events)-1].Result()
	if !ok {
		t.Fatalf("command.completed result missing")
	}

	if got, want := len(completedResult.Artifacts()), 1; got != want {
		t.Fatalf("command.completed artifact count = %d, want %d", got, want)
	}

	if got, want := completedResult.WarningCount(), 1; got != want {
		t.Fatalf("command.completed warning count = %d, want %d", got, want)
	}
}
