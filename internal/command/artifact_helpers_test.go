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

// TestArtifactNormalizationHelpers verifies text normalization.
func TestArtifactNormalizationHelpers(t *testing.T) {
	t.Parallel()

	if got, want := normalizeArtifactText("  a   b  "), "a b"; got != want {
		t.Fatalf("normalizeArtifactText() = %q, want %q", got, want)
	}

	if got, want := normalizeArtifactBlock("  one  \r\n  two  "), "one\ntwo"; got != want {
		t.Fatalf("normalizeArtifactBlock() = %q, want %q", got, want)
	}
}

// TestCloneArtifactStrings verifies detached string copy behavior.
func TestCloneArtifactStrings(t *testing.T) {
	t.Parallel()

	input := []string{"a"}
	output := cloneArtifactStrings(input)
	output[0] = "b"

	if got, want := input[0], "a"; got != want {
		t.Fatalf("cloneArtifactStrings mutated source: got %q, want %q", got, want)
	}

	if cloneArtifactStrings(nil) != nil {
		t.Fatalf("cloneArtifactStrings(nil) must return nil")
	}
}
