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
	"testing"
)

// TestKnownOptionSourcesReturnsDetachedStableOrder verifies stable precedence order.
func TestKnownOptionSourcesReturnsDetachedStableOrder(t *testing.T) {
	t.Parallel()

	got := KnownOptionSources()

	want := []OptionSource{
		OptionSourceDefault,
		OptionSourceInherited,
		OptionSourceConfig,
		OptionSourceEnvironment,
		OptionSourceRuntime,
		OptionSourceInteractive,
		OptionSourceCommandLine,
	}

	if len(got) != len(want) {
		t.Fatalf("KnownOptionSources length = %d, want %d", len(got), len(want))
	}

	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("KnownOptionSources()[%d] = %q, want %q", index, got[index], want[index])
		}
	}

	got[0] = OptionSource("changed")

	again := KnownOptionSources()
	if again[0] != OptionSourceDefault {
		t.Fatalf("KnownOptionSources returned mutable state: got %q, want %q", again[0], OptionSourceDefault)
	}
}

// TestOptionSourceString verifies canonical string rendering.
func TestOptionSourceString(t *testing.T) {
	t.Parallel()

	if got, want := OptionSourceCommandLine.String(), "command-line"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
