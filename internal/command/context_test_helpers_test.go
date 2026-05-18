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

func contextTestCommandNode() Node {
	return MustCommandNode(
		MustID("bench.run"),
		MustPath("bench", "run"),
		"run",
	)
}

func contextTestOtherCommandNode() Node {
	return MustCommandNode(
		MustID("check"),
		MustPath("check"),
		"check",
	)
}

func contextTestFamilyNode() Node {
	return MustFamilyNode(
		MustID("bench"),
		MustPath("bench"),
		"bench",
		contextTestCommandNode(),
	)
}

func contextTestField(t *testing.T, commandContext Context, key string) string {
	t.Helper()

	value, ok := commandContext.Field(key)
	if !ok {
		t.Fatalf("context field %q not found", key)
	}

	return value
}

func invocationTestEnv(t *testing.T, invocation Invocation, name string) string {
	t.Helper()

	value, ok := invocation.EnvValue(name)
	if !ok {
		t.Fatalf("env %q not found", name)
	}

	return value
}

func invocationTestField(t *testing.T, invocation Invocation, key string) string {
	t.Helper()

	value, ok := invocation.Field(key)
	if !ok {
		t.Fatalf("invocation field %q not found", key)
	}

	return value
}
