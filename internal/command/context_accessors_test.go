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
	"time"
)

func TestContextDelegatesBaseContext(t *testing.T) {
	t.Parallel()

	type contextKey string

	base := context.WithValue(context.Background(), contextKey("key"), "value")
	deadline := time.Now().Add(time.Hour)
	base, cancel := context.WithDeadline(base, deadline)
	defer cancel()

	commandContext := MustContext(ContextSpec{
		Context: base,
		Node:    contextTestCommandNode(),
	})

	gotDeadline, ok := commandContext.Deadline()
	if !ok {
		t.Fatalf("Deadline() ok = false, want true")
	}

	if !gotDeadline.Equal(deadline) {
		t.Fatalf("Deadline() = %v, want %v", gotDeadline, deadline)
	}

	if got := commandContext.Value(contextKey("key")); got != "value" {
		t.Fatalf("Value() = %v, want value", got)
	}

	cancel()

	if !errors.Is(commandContext.Err(), context.Canceled) {
		t.Fatalf("Err() = %v, want context.Canceled", commandContext.Err())
	}
}

func TestContextAccessorsAndOrdering(t *testing.T) {
	t.Parallel()

	commandContext := MustContext(ContextSpec{
		Node: contextTestCommandNode(),
		Fields: map[string]string{
			"z.key": "z",
			"a.key": "a",
			"m.key": "m",
		},
	})

	if commandContext.IsZero() {
		t.Fatalf("IsZero() = true, want false")
	}

	if !commandContext.HasFields() {
		t.Fatalf("HasFields() = false, want true")
	}

	if got, want := commandContext.FieldCount(), 3; got != want {
		t.Fatalf("FieldCount() = %d, want %d", got, want)
	}

	if !commandContext.HasField("a.key") {
		t.Fatalf("HasField(a.key) = false, want true")
	}

	assertStringSlicesEqual(t, commandContext.FieldKeys(), []string{"a.key", "m.key", "z.key"})
}

func TestContextCopySemantics(t *testing.T) {
	t.Parallel()

	fields := map[string]string{"source": "test"}
	invocation := MustInvocation(InvocationSpec{
		Arguments: []string{"one"},
		Env:       map[string]string{"GOOS": "linux"},
	})

	commandContext := MustContext(ContextSpec{
		Node:       contextTestCommandNode(),
		Invocation: invocation,
		Fields:     fields,
	})

	fields["source"] = "changed"

	if got, want := contextTestField(t, commandContext, "source"), "test"; got != want {
		t.Fatalf("field changed through input map: got %q, want %q", got, want)
	}

	outFields := commandContext.Fields()
	outFields["source"] = "changed"

	if got, want := contextTestField(t, commandContext, "source"), "test"; got != want {
		t.Fatalf("field changed through output map: got %q, want %q", got, want)
	}

	outSpec := commandContext.Spec()
	outSpec.Fields["source"] = "changed"
	outSpec.Invocation = outSpec.Invocation.MustWithArguments("changed")

	if got, want := contextTestField(t, commandContext, "source"), "test"; got != want {
		t.Fatalf("field changed through Spec: got %q, want %q", got, want)
	}

	if got, want := commandContext.Invocation().Arguments()[0], "one"; got != want {
		t.Fatalf("invocation changed through Spec: got %q, want %q", got, want)
	}
}
