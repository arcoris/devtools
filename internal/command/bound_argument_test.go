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

// TestBoundArgumentHelpers verifies bound argument accessors.
func TestBoundArgumentHelpers(t *testing.T) {
	t.Parallel()

	argument := bindingTestSuiteArgument()
	bound := BoundArgument{
		argument: argument,
		values:   []string{"stable"},
	}

	if got, want := bound.Name(), MustArgumentName("suite"); got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}

	if got, want := bound.Kind(), OptionKindEnum; got != want {
		t.Fatalf("Kind() = %q, want %q", got, want)
	}

	if got, want := bound.String(), "stable"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	if got, want := bound.MustValue(), "stable"; got != want {
		t.Fatalf("MustValue() = %q, want %q", got, want)
	}

	values := bound.Values()
	values[0] = "changed"

	if got, want := bound.Values()[0], "stable"; got != want {
		t.Fatalf("Values() returned mutable state: got %q, want %q", got, want)
	}
}
