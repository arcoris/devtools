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
	"testing"
)

func TestContextActionRequest(t *testing.T) {
	t.Parallel()

	commandContext := MustContext(ContextSpec{
		Node: contextTestCommandNode(),
		Invocation: MustInvocation(InvocationSpec{
			Arguments: []string{"one"},
			Fields: map[string]string{
				"source": "invocation",
				"mode":   "smoke",
			},
		}),
		Fields: map[string]string{
			"source": "context",
			"trace":  "enabled",
		},
	})

	request, err := commandContext.ActionRequest()
	if err != nil {
		t.Fatalf("ActionRequest() returned unexpected error: %v", err)
	}

	if got, want := request.Node.ID(), MustID("bench.run"); got != want {
		t.Fatalf("ActionRequest Node ID = %q, want %q", got, want)
	}

	if got, want := request.Args()[0], "one"; got != want {
		t.Fatalf("ActionRequest arg = %q, want %q", got, want)
	}

	if got, want := mustActionRequestField(t, request, "source"), "context"; got != want {
		t.Fatalf("source field = %q, want %q", got, want)
	}

	if got, want := mustActionRequestField(t, request, "mode"), "smoke"; got != want {
		t.Fatalf("mode field = %q, want %q", got, want)
	}

	if got, want := mustActionRequestField(t, request, "trace"), "enabled"; got != want {
		t.Fatalf("trace field = %q, want %q", got, want)
	}
}

func TestContextActionRequestRejectsFamily(t *testing.T) {
	t.Parallel()

	commandContext := MustContext(ContextSpec{
		Node: contextTestFamilyNode(),
	})

	_, err := commandContext.ActionRequest()
	if err == nil {
		t.Fatalf("ActionRequest() returned nil error")
	}

	if !errors.Is(err, ErrInvalidActionRequest) {
		t.Fatalf("ActionRequest() error = %v, want ErrInvalidActionRequest", err)
	}
}

func TestContextMustActionRequestPanicsForInvalidContextNode(t *testing.T) {
	t.Parallel()

	commandContext := MustContext(ContextSpec{
		Node: contextTestFamilyNode(),
	})

	assertPanics(t, func() {
		_ = commandContext.MustActionRequest()
	})
}

func mustActionRequestField(t *testing.T, request ActionRequest, key string) string {
	t.Helper()

	value, ok := request.Field(key)
	if !ok {
		t.Fatalf("request field %q not found", key)
	}

	return value
}
