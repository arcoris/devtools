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
	"strings"
	"testing"
)

func TestNewInvocationAcceptsValidInvocation(t *testing.T) {
	t.Parallel()

	invocation, err := NewInvocation(InvocationSpec{
		Arguments:    []string{"suite"},
		RawArguments: []string{"bench", "run", "--suite", "stable"},
		WorkingDir:   "/repo",
		Env: map[string]string{
			"GOOS":   "linux",
			"GOARCH": "amd64",
		},
		Fields: map[string]string{
			"source":  "test",
			"ci.mode": "smoke",
		},
	})
	if err != nil {
		t.Fatalf("NewInvocation() returned unexpected error: %v", err)
	}

	if got, want := invocation.ArgumentCount(), 1; got != want {
		t.Fatalf("ArgumentCount() = %d, want %d", got, want)
	}

	if got, want := invocation.RawArgumentCount(), 4; got != want {
		t.Fatalf("RawArgumentCount() = %d, want %d", got, want)
	}

	if got, want := invocation.WorkingDir(), "/repo"; got != want {
		t.Fatalf("WorkingDir() = %q, want %q", got, want)
	}

	if got, want := invocationTestEnv(t, invocation, "GOOS"), "linux"; got != want {
		t.Fatalf("EnvValue(GOOS) = %q, want %q", got, want)
	}

	if got, want := invocationTestField(t, invocation, "source"), "test"; got != want {
		t.Fatalf("Field(source) = %q, want %q", got, want)
	}
}

func TestEmptyInvocation(t *testing.T) {
	t.Parallel()

	invocation := EmptyInvocation()

	if !invocation.IsZero() {
		t.Fatalf("EmptyInvocation().IsZero() = false, want true")
	}

	if err := invocation.Validate(); err != nil {
		t.Fatalf("EmptyInvocation().Validate() returned unexpected error: %v", err)
	}
}

func TestNewInvocationRejectsInvalidInvocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec InvocationSpec
	}{
		{name: "invalid argument", spec: InvocationSpec{Arguments: []string{"bad\x00arg"}}},
		{name: "too long argument", spec: InvocationSpec{Arguments: []string{strings.Repeat("x", maxInvocationArgumentLength+1)}}},
		{name: "invalid raw argument", spec: InvocationSpec{RawArguments: []string{"bad\x00arg"}}},
		{name: "too long raw argument", spec: InvocationSpec{RawArguments: []string{strings.Repeat("x", maxInvocationRawArgumentLength+1)}}},
		{name: "blank working directory", spec: InvocationSpec{WorkingDir: "   "}},
		{name: "invalid working directory", spec: InvocationSpec{WorkingDir: "bad\x00dir"}},
		{name: "too long working directory", spec: InvocationSpec{WorkingDir: strings.Repeat("x", maxInvocationWorkingDirLength+1)}},
		{name: "invalid env name", spec: InvocationSpec{Env: map[string]string{"goos": "linux"}}},
		{name: "invalid env value", spec: InvocationSpec{Env: map[string]string{"GOOS": "bad\x00value"}}},
		{name: "invalid field key", spec: InvocationSpec{Fields: map[string]string{"Invalid": "value"}}},
		{name: "invalid field value", spec: InvocationSpec{Fields: map[string]string{"source": "bad\x00value"}}},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewInvocation(test.spec)
			if err == nil {
				t.Fatalf("NewInvocation() returned nil error")
			}

			if !errors.Is(err, ErrInvalidInvocation) {
				t.Fatalf("NewInvocation() error = %v, want ErrInvalidInvocation", err)
			}
		})
	}
}

func TestMustInvocationPanicsForInvalidInvocation(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustInvocation(InvocationSpec{Arguments: []string{"bad\x00arg"}})
	})
}
