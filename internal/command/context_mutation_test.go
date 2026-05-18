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

func TestContextWithHelpers(t *testing.T) {
	t.Parallel()

	startedAt := time.Date(2026, 5, 18, 11, 0, 0, 0, time.UTC)
	base := context.WithValue(context.Background(), "key", "value")

	commandContext := BackgroundContext(contextTestCommandNode()).
		MustWithBase(base).
		MustWithNode(contextTestOtherCommandNode()).
		MustWithStartedAt(startedAt).
		MustWithInvocation(MustInvocation(InvocationSpec{Arguments: []string{"one"}})).
		MustWithField("source", "test").
		MustWithFields(map[string]string{"mode": "smoke", "source": "test"})

	if got, want := commandContext.ID(), MustID("check"); got != want {
		t.Fatalf("WithNode() ID = %q, want %q", got, want)
	}

	if got, want := commandContext.StartedAt(), startedAt; !got.Equal(want) {
		t.Fatalf("StartedAt() = %v, want %v", got, want)
	}

	if got := commandContext.Value("key"); got != "value" {
		t.Fatalf("Value() = %v, want value", got)
	}

	if commandContext.WithoutField("mode").HasField("mode") {
		t.Fatalf("WithoutField() still has removed field")
	}

	if commandContext.WithoutFields().HasFields() {
		t.Fatalf("WithoutFields() still has fields")
	}
}

func TestContextWithFieldsCopiesInput(t *testing.T) {
	t.Parallel()

	fields := map[string]string{"source": "test"}
	commandContext := BackgroundContext(contextTestCommandNode()).MustWithFields(fields)

	fields["source"] = "changed"

	if got, want := contextTestField(t, commandContext, "source"), "test"; got != want {
		t.Fatalf("field changed through WithFields input: got %q, want %q", got, want)
	}
}

func TestContextWithHelpersRejectInvalidValues(t *testing.T) {
	t.Parallel()

	commandContext := BackgroundContext(contextTestCommandNode())

	if _, err := commandContext.WithNode(Node{}); !errors.Is(err, ErrInvalidContext) {
		t.Fatalf("WithNode() error = %v, want ErrInvalidContext", err)
	}

	if _, err := commandContext.WithInvocation(Invocation{arguments: []string{"bad\x00arg"}}); !errors.Is(err, ErrInvalidContext) {
		t.Fatalf("WithInvocation() error = %v, want ErrInvalidContext", err)
	}

	if _, err := commandContext.WithField("Invalid", "value"); !errors.Is(err, ErrInvalidContext) {
		t.Fatalf("WithField(key) error = %v, want ErrInvalidContext", err)
	}

	if _, err := commandContext.WithFields(map[string]string{"source": "bad\x00value"}); !errors.Is(err, ErrInvalidContext) {
		t.Fatalf("WithFields(value) error = %v, want ErrInvalidContext", err)
	}
}
