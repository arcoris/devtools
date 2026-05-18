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

func TestNewContextAcceptsValidContext(t *testing.T) {
	t.Parallel()

	startedAt := time.Date(2026, 5, 18, 10, 0, 0, 0, time.UTC)

	commandContext, err := NewContext(ContextSpec{
		Context: context.Background(),
		Node:    contextTestCommandNode(),
		Invocation: MustInvocation(InvocationSpec{
			Arguments:    []string{"suite"},
			RawArguments: []string{"bench", "run", "--suite", "stable"},
			WorkingDir:   "/repo",
			Env:          map[string]string{"GOOS": "linux"},
			Fields:       map[string]string{"invocation.source": "adapter"},
		}),
		StartedAt: startedAt,
		Fields:    map[string]string{"source": "test"},
	})
	if err != nil {
		t.Fatalf("NewContext() returned unexpected error: %v", err)
	}

	if got, want := commandContext.ID(), MustID("bench.run"); got != want {
		t.Fatalf("ID() = %q, want %q", got, want)
	}

	if got, want := commandContext.Path(), MustPath("bench", "run"); !got.Equal(want) {
		t.Fatalf("Path() = %q, want %q", got, want)
	}

	if got, want := commandContext.StartedAt(), startedAt; !got.Equal(want) {
		t.Fatalf("StartedAt() = %v, want %v", got, want)
	}

	if got, want := contextTestField(t, commandContext, "source"), "test"; got != want {
		t.Fatalf("Field(source) = %q, want %q", got, want)
	}
}

func TestNewContextDefaults(t *testing.T) {
	t.Parallel()

	commandContext, err := NewContext(ContextSpec{
		Node: contextTestCommandNode(),
	})
	if err != nil {
		t.Fatalf("NewContext() returned unexpected error: %v", err)
	}

	if commandContext.Base() == nil {
		t.Fatalf("Base() is nil")
	}

	if commandContext.StartedAt().IsZero() {
		t.Fatalf("StartedAt() is zero")
	}
}

func TestBackgroundContext(t *testing.T) {
	t.Parallel()

	commandContext := BackgroundContext(contextTestCommandNode())

	if commandContext.Base() == nil {
		t.Fatalf("Base() is nil")
	}

	if got, want := commandContext.ID(), MustID("bench.run"); got != want {
		t.Fatalf("ID() = %q, want %q", got, want)
	}
}

func TestNewContextRejectsInvalidContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec ContextSpec
	}{
		{
			name: "missing node",
			spec: ContextSpec{},
		},
		{
			name: "invalid invocation",
			spec: ContextSpec{
				Node: contextTestCommandNode(),
				Invocation: Invocation{
					arguments: []string{"bad\x00arg"},
				},
			},
		},
		{
			name: "invalid field key",
			spec: ContextSpec{
				Node:   contextTestCommandNode(),
				Fields: map[string]string{"Invalid": "value"},
			},
		},
		{
			name: "invalid field value",
			spec: ContextSpec{
				Node:   contextTestCommandNode(),
				Fields: map[string]string{"source": "bad\x00value"},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewContext(test.spec)
			if err == nil {
				t.Fatalf("NewContext() returned nil error")
			}

			if !errors.Is(err, ErrInvalidContext) {
				t.Fatalf("NewContext() error = %v, want ErrInvalidContext", err)
			}
		})
	}
}

func TestMustContextPanicsForInvalidContext(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustContext(ContextSpec{})
	})
}
