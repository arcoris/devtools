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
	"testing"
)

// TestRuntimeExecuteSuccess verifies successful runtime execution.
func TestRuntimeExecuteSuccess(t *testing.T) {
	t.Parallel()

	collector := &RuntimeEventCollector{}
	startedAt := runtimeTestTime()

	runtime := MustRuntime(RuntimeSpec{
		Name:      "test-runtime",
		CommandID: MustID("bench.run"),
		Binding:   runtimeTestBinding(),
		Clock:     FixedRuntimeClock{Time: startedAt},
		EventSink: collector,
		Handler: RuntimeHandlerFunc(func(ctx context.Context, request RuntimeRequest) (Result, error) {
			format, ok := request.Option(MustOptionName("format"))
			if !ok {
				t.Fatalf("format option not found")
			}

			if got, want := format.MustValue(), "json"; got != want {
				t.Fatalf("format = %q, want %q", got, want)
			}

			suite, ok := request.Argument(MustArgumentName("suite"))
			if !ok {
				t.Fatalf("suite argument not found")
			}

			if got, want := suite.MustValue(), "stable"; got != want {
				t.Fatalf("suite = %q, want %q", got, want)
			}

			return OKResult("completed"), nil
		}),
	})

	result, err := runtime.Execute(context.Background(), RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	if !result.IsOK() {
		t.Fatalf("result IsOK() = false, want true")
	}

	if got, want := result.Message(), "completed"; got != want {
		t.Fatalf("result Message() = %q, want %q", got, want)
	}

	if _, ok := result.Duration(); !ok {
		t.Fatalf("result Duration() ok = false, want true")
	}

	if got, want := collector.Len(), 6; got != want {
		t.Fatalf("collector Len() = %d, want %d", got, want)
	}
}
