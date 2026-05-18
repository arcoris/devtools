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

// TestEventWithHelpers verifies immutable-style event updates.
func TestEventWithHelpers(t *testing.T) {
	t.Parallel()

	event := MustSimpleEvent(EventKindCommandStarted, "").
		MustWithID("command.started.001").
		MustWithKind(EventKindCommandCompleted).
		MustWithSeverity(EventSeverityWarning).
		MustWithMessage("Completed with warnings.").
		MustWithField("mode", "ci").
		MustWithArtifact(eventTestArtifact("bench.report")).
		MustWithResult(OKResult("ok")).
		MustWithLabel("bench").
		MustWithMetadata(MustMetadata(MetadataSpec{Owner: "devtools"})).
		MustWithVisibility(VisibilityHidden)

	if !event.HasID() {
		t.Fatalf("HasID() = false, want true")
	}

	if got, want := event.Kind(), EventKindCommandCompleted; got != want {
		t.Fatalf("Kind() = %q, want %q", got, want)
	}

	if !event.IsWarning() {
		t.Fatalf("IsWarning() = false, want true")
	}

	if got, want := event.Message(), "Completed with warnings."; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}

	if !event.HasField("mode") {
		t.Fatalf("HasField(mode) = false, want true")
	}

	if !event.HasArtifact(MustArtifactID("bench.report")) {
		t.Fatalf("HasArtifact(bench.report) = false, want true")
	}

	if !event.HasResult() {
		t.Fatalf("HasResult() = false, want true")
	}

	if !event.HasLabel("bench") {
		t.Fatalf("HasLabel(bench) = false, want true")
	}

	if got, want := event.Metadata().Owner(), "devtools"; got != want {
		t.Fatalf("Metadata().Owner() = %q, want %q", got, want)
	}

	if !event.Visibility().IsHidden() {
		t.Fatalf("Visibility().IsHidden() = false, want true")
	}

	withoutID := event.WithoutID()
	if withoutID.HasID() {
		t.Fatalf("WithoutID() still has ID")
	}

	withoutField := event.WithoutField("mode")
	if withoutField.HasField("mode") {
		t.Fatalf("WithoutField() still has field")
	}

	withoutArtifact := event.WithoutArtifact(MustArtifactID("bench.report"))
	if withoutArtifact.HasArtifact(MustArtifactID("bench.report")) {
		t.Fatalf("WithoutArtifact() still has artifact")
	}

	withoutResult := event.WithoutResult()
	if withoutResult.HasResult() {
		t.Fatalf("WithoutResult() still has result")
	}

	withoutLabel := event.WithoutLabel("bench")
	if withoutLabel.HasLabel("bench") {
		t.Fatalf("WithoutLabel() still has label")
	}
}
