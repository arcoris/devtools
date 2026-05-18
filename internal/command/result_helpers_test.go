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

// TestResultCloneHelpers verifies detached helper behavior.
func TestResultCloneHelpers(t *testing.T) {
	t.Parallel()

	stringsValue := []string{"a"}
	clonedStrings := cloneResultStrings(stringsValue)
	clonedStrings[0] = "b"

	if got, want := stringsValue[0], "a"; got != want {
		t.Fatalf("cloneResultStrings mutated source: got %q, want %q", got, want)
	}

	fields := map[string]string{"mode": "ci"}
	clonedFields := cloneResultStringMap(fields)
	clonedFields["mode"] = "changed"

	if got, want := fields["mode"], "ci"; got != want {
		t.Fatalf("cloneResultStringMap mutated source: got %q, want %q", got, want)
	}

	if cloneResultStrings(nil) != nil {
		t.Fatalf("cloneResultStrings(nil) must return nil")
	}

	if cloneResultStringMap(nil) != nil {
		t.Fatalf("cloneResultStringMap(nil) must return nil")
	}
}

// TestResultNormalizationHelpers verifies normalization helpers.
func TestResultNormalizationHelpers(t *testing.T) {
	t.Parallel()

	if got, want := normalizeResultBlock("  One.  \r\n  Two.  "), "One.\nTwo."; got != want {
		t.Fatalf("normalizeResultBlock() = %q, want %q", got, want)
	}

	hints := normalizeResultHints([]string{"  One.  ", " ", "  Two.  "})
	resultTestAssertStrings(t, hints, []string{"One.", "Two."})
}
