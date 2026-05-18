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

func TestInvocationWithHelpers(t *testing.T) {
	t.Parallel()

	invocation := EmptyInvocation().
		MustWithArguments("one", "two").
		MustWithRawArguments("bench", "run").
		MustWithWorkingDir("/repo").
		MustWithEnvValue("GOOS", "linux").
		MustWithEnvValue("GOARCH", "amd64").
		MustWithField("source", "test").
		MustWithField("ci.mode", "smoke")

	if got, want := invocation.ArgumentCount(), 2; got != want {
		t.Fatalf("ArgumentCount() = %d, want %d", got, want)
	}

	if got, want := invocation.RawArgumentCount(), 2; got != want {
		t.Fatalf("RawArgumentCount() = %d, want %d", got, want)
	}

	if got, want := invocation.WorkingDir(), "/repo"; got != want {
		t.Fatalf("WorkingDir() = %q, want %q", got, want)
	}

	assertStringSlicesEqual(t, invocation.EnvNames(), []string{"GOARCH", "GOOS"})
	assertStringSlicesEqual(t, invocation.FieldKeys(), []string{"ci.mode", "source"})

	if invocation.WithoutArguments().HasArguments() {
		t.Fatalf("WithoutArguments() still has arguments")
	}

	if invocation.WithoutRawArguments().HasRawArguments() {
		t.Fatalf("WithoutRawArguments() still has raw arguments")
	}

	if invocation.WithoutWorkingDir().HasWorkingDir() {
		t.Fatalf("WithoutWorkingDir() still has working directory")
	}

	if invocation.WithoutEnvValue("GOARCH").HasEnvValue("GOARCH") {
		t.Fatalf("WithoutEnvValue() still has removed value")
	}

	if invocation.WithoutEnv().HasEnv() {
		t.Fatalf("WithoutEnv() still has env")
	}

	if invocation.WithoutField("source").HasField("source") {
		t.Fatalf("WithoutField() still has removed value")
	}

	if invocation.WithoutFields().HasFields() {
		t.Fatalf("WithoutFields() still has fields")
	}
}

func TestInvocationWithCollectionsCopiesInput(t *testing.T) {
	t.Parallel()

	env := map[string]string{"GOOS": "linux"}
	fields := map[string]string{"source": "test"}

	invocation := EmptyInvocation().
		MustWithEnv(env).
		MustWithFields(fields)

	env["GOOS"] = "changed"
	fields["source"] = "changed"

	if got, want := invocationTestEnv(t, invocation, "GOOS"), "linux"; got != want {
		t.Fatalf("env changed through WithEnv input: got %q, want %q", got, want)
	}

	if got, want := invocationTestField(t, invocation, "source"), "test"; got != want {
		t.Fatalf("field changed through WithFields input: got %q, want %q", got, want)
	}
}

func TestInvocationWithHelpersRejectInvalidValues(t *testing.T) {
	t.Parallel()

	invocation := EmptyInvocation()

	if _, err := invocation.WithArguments("bad\x00arg"); !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("WithArguments() error = %v, want ErrInvalidInvocation", err)
	}

	if _, err := invocation.WithRawArguments("bad\x00arg"); !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("WithRawArguments() error = %v, want ErrInvalidInvocation", err)
	}

	if _, err := invocation.WithWorkingDir("   "); !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("WithWorkingDir() error = %v, want ErrInvalidInvocation", err)
	}

	if _, err := invocation.WithEnv(map[string]string{"goos": "linux"}); !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("WithEnv() error = %v, want ErrInvalidInvocation", err)
	}

	if _, err := invocation.WithEnvValue("goos", "linux"); !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("WithEnvValue() error = %v, want ErrInvalidInvocation", err)
	}

	if _, err := invocation.WithFields(map[string]string{"Invalid": "value"}); !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("WithFields() error = %v, want ErrInvalidInvocation", err)
	}

	if _, err := invocation.WithField("Invalid", "value"); !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("WithField() error = %v, want ErrInvalidInvocation", err)
	}
}
