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

// TestEventCloneHelpers verifies detached clone helper behavior.
func TestEventCloneHelpers(t *testing.T) {
	t.Parallel()

	stringsValue := []string{"a"}
	clonedStrings := cloneEventStrings(stringsValue)
	clonedStrings[0] = "b"

	if got, want := stringsValue[0], "a"; got != want {
		t.Fatalf("cloneEventStrings mutated source: got %q, want %q", got, want)
	}

	fields := map[string]string{"mode": "ci"}
	clonedFields := cloneEventStringMap(fields)
	clonedFields["mode"] = "changed"

	if got, want := fields["mode"], "ci"; got != want {
		t.Fatalf("cloneEventStringMap mutated source: got %q, want %q", got, want)
	}

	artifacts := []Artifact{eventTestArtifact("bench.report")}
	clonedArtifacts := cloneEventArtifacts(artifacts)
	clonedArtifacts[0] = eventTestArtifact("changed")

	if artifacts[0].ID() != MustArtifactID("bench.report") {
		t.Fatalf("cloneEventArtifacts mutated source")
	}

	if cloneEventStrings(nil) != nil {
		t.Fatalf("cloneEventStrings(nil) must return nil")
	}

	if cloneEventStringMap(nil) != nil {
		t.Fatalf("cloneEventStringMap(nil) must return nil")
	}

	if cloneEventArtifacts(nil) != nil {
		t.Fatalf("cloneEventArtifacts(nil) must return nil")
	}
}

// TestEventNormalizationHelpers verifies text normalization.
func TestEventNormalizationHelpers(t *testing.T) {
	t.Parallel()

	if got, want := normalizeEventBlock("  One.  \r\n  Two.  "), "One.\nTwo."; got != want {
		t.Fatalf("normalizeEventBlock() = %q, want %q", got, want)
	}
}
