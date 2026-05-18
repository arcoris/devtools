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

import "testing"

func TestInvocationIndexedAccess(t *testing.T) {
	t.Parallel()

	invocation := MustInvocation(InvocationSpec{
		Arguments:    []string{"one"},
		RawArguments: []string{"raw"},
	})

	if got, ok := invocation.Argument(0); !ok || got != "one" {
		t.Fatalf("Argument(0) = %q, %v; want one, true", got, ok)
	}

	if got, ok := invocation.Argument(1); ok || got != "" {
		t.Fatalf("Argument(1) = %q, %v; want empty, false", got, ok)
	}

	if got, ok := invocation.RawArgument(0); !ok || got != "raw" {
		t.Fatalf("RawArgument(0) = %q, %v; want raw, true", got, ok)
	}

	if got, ok := invocation.RawArgument(1); ok || got != "" {
		t.Fatalf("RawArgument(1) = %q, %v; want empty, false", got, ok)
	}
}

func TestInvocationAccessorsAndOrdering(t *testing.T) {
	t.Parallel()

	invocation := MustInvocation(InvocationSpec{
		Arguments:    []string{"one"},
		RawArguments: []string{"raw"},
		WorkingDir:   "/repo",
		Env: map[string]string{
			"Z_KEY": "z",
			"A_KEY": "a",
			"M_KEY": "m",
		},
		Fields: map[string]string{
			"z.key": "z",
			"a.key": "a",
			"m.key": "m",
		},
	})

	if !invocation.HasArguments() || invocation.ArgumentCount() != 1 {
		t.Fatalf("argument accessors returned unexpected values")
	}

	if !invocation.HasRawArguments() || invocation.RawArgumentCount() != 1 {
		t.Fatalf("raw argument accessors returned unexpected values")
	}

	if !invocation.HasWorkingDir() {
		t.Fatalf("HasWorkingDir() = false, want true")
	}

	if !invocation.HasEnv() || invocation.EnvCount() != 3 {
		t.Fatalf("env accessors returned unexpected values")
	}

	if !invocation.HasEnvValue("A_KEY") {
		t.Fatalf("HasEnvValue(A_KEY) = false, want true")
	}

	if !invocation.HasFields() || invocation.FieldCount() != 3 {
		t.Fatalf("field accessors returned unexpected values")
	}

	if !invocation.HasField("a.key") {
		t.Fatalf("HasField(a.key) = false, want true")
	}

	assertStringSlicesEqual(t, invocation.EnvNames(), []string{"A_KEY", "M_KEY", "Z_KEY"})
	assertStringSlicesEqual(t, invocation.FieldKeys(), []string{"a.key", "m.key", "z.key"})
}

func TestInvocationCopySemantics(t *testing.T) {
	t.Parallel()

	args := []string{"one"}
	rawArgs := []string{"raw"}
	env := map[string]string{"GOOS": "linux"}
	fields := map[string]string{"source": "test"}

	invocation := MustInvocation(InvocationSpec{
		Arguments:    args,
		RawArguments: rawArgs,
		Env:          env,
		Fields:       fields,
	})

	args[0] = "changed"
	rawArgs[0] = "changed"
	env["GOOS"] = "changed"
	fields["source"] = "changed"

	if got, want := invocation.Arguments()[0], "one"; got != want {
		t.Fatalf("argument changed through input slice: got %q, want %q", got, want)
	}

	if got, want := invocation.RawArguments()[0], "raw"; got != want {
		t.Fatalf("raw argument changed through input slice: got %q, want %q", got, want)
	}

	if got, want := invocationTestEnv(t, invocation, "GOOS"), "linux"; got != want {
		t.Fatalf("env changed through input map: got %q, want %q", got, want)
	}

	if got, want := invocationTestField(t, invocation, "source"), "test"; got != want {
		t.Fatalf("field changed through input map: got %q, want %q", got, want)
	}

	outSpec := invocation.Spec()
	outSpec.Arguments[0] = "changed"
	outSpec.RawArguments[0] = "changed"
	outSpec.Env["GOOS"] = "changed"
	outSpec.Fields["source"] = "changed"

	if got, want := invocation.Arguments()[0], "one"; got != want {
		t.Fatalf("argument changed through Spec: got %q, want %q", got, want)
	}

	if got, want := invocation.RawArguments()[0], "raw"; got != want {
		t.Fatalf("raw argument changed through Spec: got %q, want %q", got, want)
	}

	if got, want := invocationTestEnv(t, invocation, "GOOS"), "linux"; got != want {
		t.Fatalf("env changed through Spec: got %q, want %q", got, want)
	}

	if got, want := invocationTestField(t, invocation, "source"), "test"; got != want {
		t.Fatalf("field changed through Spec: got %q, want %q", got, want)
	}
}
