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

func TestCommandLifecycleIntegrationSuccess(t *testing.T) {
	t.Parallel()

	binding := optionResolverTestBinding()
	collector := &RuntimeEventCollector{}

	values, err := ResolveOptionValues(OptionResolutionSpec{
		Binding: binding,
		EnvironmentValues: []OptionValue{
			MustScalarOptionValue("output", OptionKindString, OptionSourceEnvironment, "bench.out"),
		},
		CommandLineValues: []OptionValue{
			MustScalarOptionValue("fmt", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
	})
	if err != nil {
		t.Fatalf("ResolveOptionValues() returned unexpected error: %v", err)
	}

	node := MustNode(NodeSpec{
		Kind:    NodeCommand,
		ID:      MustID("bench.run"),
		Path:    MustPath("bench", "run"),
		Use:     "run",
		Binding: binding,
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			format, ok := request.Option(MustOptionName("format"))
			if !ok || format.Name() != MustOptionName("format") || format.MustValue() != "json" {
				t.Fatalf("format option = %v, %v; want canonical json", format, ok)
			}

			packages, ok := request.Option(MustOptionName("package"))
			if !ok {
				t.Fatalf("package default option missing")
			}

			assertOptionValueStrings(t, packages.Values(), []string{"./..."})

			suite, ok := request.Argument(MustArgumentName("suite"))
			if !ok || suite.MustValue() != "smoke" {
				t.Fatalf("suite argument = %v, %v; want smoke", suite, ok)
			}

			return MustResult(ResultSpec{
				Status:  ResultStatusOK,
				Message: "benchmark complete",
				Artifacts: []Artifact{
					MustArtifact(ArtifactSpec{
						ID:          "benchmark.report",
						Kind:        ArtifactKindReport,
						Location:    "bench/reports/benchmark.md",
						Format:      ArtifactFormatMarkdown,
						Description: "Benchmark report",
					}),
				},
				Warnings: []ResultWarning{
					MustResultWarning(ResultWarningSpec{
						Kind:    "partial",
						Message: "some benchmarks were skipped",
					}),
				},
				Fields: map[string]string{
					"suite": "smoke",
				},
			}), nil
		}),
		Metadata:   MustMetadata(MetadataSpec{Owner: "devtools"}),
		Visibility: VisibilityHidden,
	})

	runtime := MustRuntimeFromNode(RuntimeFromNodeSpec{
		Node:      node,
		Clock:     FixedRuntimeClock{Time: runtimeTestTime()},
		EventSink: collector,
	})

	result, err := runtime.Execute(context.Background(), RuntimeExecutionSpec{
		OptionValues:     values,
		PositionalValues: []string{"smoke"},
	})
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	if !result.IsOK() {
		t.Fatalf("result IsOK() = false, want true")
	}

	if got, want := len(result.Artifacts()), 1; got != want {
		t.Fatalf("artifact count = %d, want %d", got, want)
	}

	if got, want := result.WarningCount(), 1; got != want {
		t.Fatalf("WarningCount() = %d, want %d", got, want)
	}

	assertEventKinds(t, collector.Events(), []EventKind{
		RuntimeEventCommandStarted,
		RuntimeEventBindingStarted,
		RuntimeEventBindingCompleted,
		RuntimeEventActionStarted,
		RuntimeEventActionCompleted,
		RuntimeEventCommandCompleted,
	})
}

func TestCommandLifecycleIntegrationRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		values   []OptionValue
		position []string
	}{
		{
			name:     "missing required option",
			values:   nil,
			position: []string{"smoke"},
		},
		{
			name: "invalid positional",
			values: []OptionValue{
				MustScalarOptionValue("output", OptionKindString, OptionSourceCommandLine, "bench.out"),
			},
			position: []string{"nightly"},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runtime := MustRuntime(RuntimeSpec{
				Binding: optionResolverTestBinding(),
				Clock:   FixedRuntimeClock{Time: runtimeTestTime()},
				Handler: runtimeTestOKHandler(),
			})

			result, err := runtime.Execute(context.Background(), RuntimeExecutionSpec{
				OptionValues:     test.values,
				PositionalValues: test.position,
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
		})
	}
}

func TestCommandLifecycleIntegrationHandlerErrorAndPanic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		handler RuntimeHandler
		wantErr error
	}{
		{
			name: "handler error",
			handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
				return FailedResult("handler failed"), errors.New("handler failed")
			}),
			wantErr: ErrRuntimeExecution,
		},
		{
			name: "panic recovery",
			handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
				panic("boom")
			}),
			wantErr: ErrRuntimePanic,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runtime := MustRuntime(RuntimeSpec{
				Binding: runtimeTestBinding(),
				Clock:   FixedRuntimeClock{Time: runtimeTestTime()},
				Options: DefaultRuntimeOptions(),
				Handler: test.handler,
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

			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, test.wantErr)
			}

			if !result.IsFailed() {
				t.Fatalf("result IsFailed() = false, want true")
			}
		})
	}
}

func assertEventKinds(t *testing.T, events []Event, want []EventKind) {
	t.Helper()

	if len(events) != len(want) {
		t.Fatalf("event count = %d, want %d", len(events), len(want))
	}

	for index, event := range events {
		if event.Kind() != want[index] {
			t.Fatalf("event %d kind = %q, want %q", index, event.Kind(), want[index])
		}
	}
}
