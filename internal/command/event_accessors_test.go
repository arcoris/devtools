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

// TestEventArtifactHelpers verifies artifact ID helpers.
func TestEventArtifactHelpers(t *testing.T) {
	t.Parallel()

	event := MustEvent(EventSpec{
		Kind: EventKindArtifactProduced,
		Artifacts: []Artifact{
			eventTestArtifact("z.report"),
			eventTestArtifact("a.report"),
			eventTestArtifact("m.report"),
		},
	})

	ids := event.SortedArtifactIDs()
	if got, want := ids[0], MustArtifactID("a.report"); got != want {
		t.Fatalf("SortedArtifactIDs()[0] = %q, want %q", got, want)
	}
}

// TestEventCopySemantics verifies detached slices, maps, and result pointers.
func TestEventCopySemantics(t *testing.T) {
	t.Parallel()

	fields := map[string]string{"mode": "ci"}
	artifacts := []Artifact{eventTestArtifact("bench.report")}
	labels := []string{"bench"}
	result := OKResult("ok")

	event := MustEvent(EventSpec{
		Kind:      EventKindCommandCompleted,
		Fields:    fields,
		Artifacts: artifacts,
		Result:    &result,
		Labels:    labels,
	})

	fields["mode"] = "changed"
	artifacts[0] = eventTestArtifact("changed")
	labels[0] = "changed"
	result = FailedResult("failed")

	if got, want := eventTestMustField(t, event, "mode"), "ci"; got != want {
		t.Fatalf("field changed through input map: got %q, want %q", got, want)
	}

	if event.HasArtifact(MustArtifactID("changed")) {
		t.Fatalf("event changed through input artifact slice")
	}

	if event.HasLabel("changed") {
		t.Fatalf("event changed through input label slice")
	}

	gotResult, ok := event.Result()
	if !ok {
		t.Fatalf("Result() ok = false, want true")
	}

	if !gotResult.IsOK() {
		t.Fatalf("result changed through input pointer")
	}

	outFields := event.Fields()
	outFields["mode"] = "changed"

	if got, want := eventTestMustField(t, event, "mode"), "ci"; got != want {
		t.Fatalf("field changed through output map: got %q, want %q", got, want)
	}

	outArtifacts := event.Artifacts()
	outArtifacts[0] = eventTestArtifact("changed")

	if event.HasArtifact(MustArtifactID("changed")) {
		t.Fatalf("event changed through output artifact slice")
	}

	outLabels := event.Labels()
	outLabels[0] = "changed"

	if event.HasLabel("changed") {
		t.Fatalf("event changed through output label slice")
	}
}

// TestEventOrderingHelpers verifies field and label ordering.
func TestEventOrderingHelpers(t *testing.T) {
	t.Parallel()

	event := MustEvent(EventSpec{
		Kind: EventKindDiagnostic,
		Fields: map[string]string{
			"z.key": "z",
			"a.key": "a",
			"m.key": "m",
		},
		Labels: []string{"z", "a", "m"},
	})

	eventTestAssertStrings(t, event.FieldKeys(), []string{"a.key", "m.key", "z.key"})
	eventTestAssertStrings(t, event.SortedLabels(), []string{"a", "m", "z"})
}
