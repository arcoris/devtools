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

func TestRuntimeHandlerFromActionAdaptsToCanonicalResult(t *testing.T) {
	t.Parallel()

	node := mustTestCommandNode(t, "bench.run")

	runtime := MustRuntime(RuntimeSpec{
		Name:      "bench-runtime",
		CommandID: node.ID(),
		Binding:   runtimeTestBinding(),
		Clock:     FixedRuntimeClock{Time: runtimeTestTime()},
		Handler: RuntimeHandlerFromAction(ActionFunc(func(ctx context.Context, request ActionRequest) (ActionResult, error) {
			if got, want := request.Node.ID(), node.ID(); got != want {
				t.Fatalf("ActionRequest node ID = %q, want %q", got, want)
			}

			if got, ok := request.Argument(0); !ok || got != "stable" {
				t.Fatalf("ActionRequest argument = %q, %v; want stable", got, ok)
			}

			if got, ok := request.Field("option.format"); !ok || got != "json" {
				t.Fatalf("ActionRequest option.format = %q, %v; want json", got, ok)
			}

			return ActionResult{
				Status:  ActionStatusOK,
				Message: "action completed",
				Data:    map[string]string{"summary": "present"},
				Artifacts: []ActionArtifact{
					MustActionArtifact("report", "bench/reports/result.txt", "Benchmark report"),
				},
				Warnings: []ActionWarning{
					MustActionWarning("partial", "partial data"),
				},
				Fields: map[string]string{
					"action.summary": "ok",
				},
			}, nil
		}), node),
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

	if got, want := result.Message(), "action completed"; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}

	if got, want := len(result.Artifacts()), 1; got != want {
		t.Fatalf("artifact count = %d, want %d", got, want)
	}

	if got, want := result.WarningCount(), 1; got != want {
		t.Fatalf("WarningCount() = %d, want %d", got, want)
	}

	if got, ok := result.Field("action.data.present"); !ok || got != "true" {
		t.Fatalf("action.data.present = %q, %v; want true", got, ok)
	}
}

func TestResultStatusFromActionStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		action ActionStatus
		want   ResultStatus
	}{
		{name: "zero", want: ResultStatusOK},
		{name: "ok", action: ActionStatusOK, want: ResultStatusOK},
		{name: "skipped", action: ActionStatusSkipped, want: ResultStatusSkipped},
		{name: "failed", action: ActionStatusFailed, want: ResultStatusFailed},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := ResultStatusFromActionStatus(test.action); got != test.want {
				t.Fatalf("ResultStatusFromActionStatus() = %q, want %q", got, test.want)
			}
		})
	}
}
