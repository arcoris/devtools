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
	"time"
)

func eventTestArtifact(id string) Artifact {
	return MustArtifact(ArtifactSpec{
		ID:       id,
		Kind:     ArtifactKindReport,
		Location: id + ".md",
	})
}

func eventTestMustField(t *testing.T, event Event, key string) string {
	t.Helper()

	value, ok := event.Field(key)
	if !ok {
		t.Fatalf("field %q not found", key)
	}

	return value
}

func eventTestAssertStrings(t *testing.T, got []string, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("slice length = %d, want %d; got %v, want %v", len(got), len(want), got, want)
	}

	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("slice[%d] = %q, want %q; got %v, want %v", index, got[index], want[index], got, want)
		}
	}
}

// TestNewEventAcceptsValidEvent verifies full lifecycle event construction.
func TestNewEventAcceptsValidEvent(t *testing.T) {
	t.Parallel()

	occurredAt := time.Date(2026, 5, 18, 10, 0, 0, 0, time.UTC)
	result := OKResult("completed")

	event, err := NewEvent(EventSpec{
		ID:         "command.completed.001",
		Kind:       EventKindCommandCompleted,
		Severity:   EventSeverityInfo,
		OccurredAt: occurredAt,
		CommandID:  MustID("bench.run"),
		Message:    "Command completed.",
		Fields: map[string]string{
			"duration": "1s",
		},
		Artifacts: []Artifact{
			eventTestArtifact("bench.report"),
		},
		Result: &result,
		Labels: []string{"bench", "completed"},
		Metadata: MustMetadata(MetadataSpec{
			Owner: "devtools",
		}),
		Visibility: VisibilityPublic,
	})
	if err != nil {
		t.Fatalf("NewEvent() returned unexpected error: %v", err)
	}

	if id, ok := event.ID(); !ok || id != MustEventID("command.completed.001") {
		t.Fatalf("ID() = %q, %v; want command.completed.001, true", id, ok)
	}

	if got, want := event.Kind(), EventKindCommandCompleted; got != want {
		t.Fatalf("Kind() = %q, want %q", got, want)
	}

	if got, want := event.Severity(), EventSeverityInfo; got != want {
		t.Fatalf("Severity() = %q, want %q", got, want)
	}

	if got, want := event.OccurredAt(), occurredAt; !got.Equal(want) {
		t.Fatalf("OccurredAt() = %v, want %v", got, want)
	}

	if commandID, ok := event.CommandID(); !ok || commandID != MustID("bench.run") {
		t.Fatalf("CommandID() = %q, %v; want bench.run, true", commandID, ok)
	}

	if got, ok := event.Field("duration"); !ok || got != "1s" {
		t.Fatalf("Field(duration) = %q, %v; want 1s, true", got, ok)
	}

	if !event.HasArtifact(MustArtifactID("bench.report")) {
		t.Fatalf("HasArtifact(bench.report) = false, want true")
	}

	if !event.HasResult() {
		t.Fatalf("HasResult() = false, want true")
	}

	if !event.HasLabel("completed") {
		t.Fatalf("HasLabel(completed) = false, want true")
	}

	if got, want := event.Metadata().Owner(), "devtools"; got != want {
		t.Fatalf("Metadata().Owner() = %q, want %q", got, want)
	}
}

// TestNewEventDefaults verifies timestamp, severity, and visibility defaults.
func TestNewEventDefaults(t *testing.T) {
	t.Parallel()

	event, err := NewEvent(EventSpec{
		Kind: EventKindCommandStarted,
	})
	if err != nil {
		t.Fatalf("NewEvent() returned unexpected error: %v", err)
	}

	if event.HasID() {
		t.Fatalf("HasID() = true, want false")
	}

	if event.OccurredAt().IsZero() {
		t.Fatalf("OccurredAt() is zero")
	}

	if got, want := event.Severity(), EventSeverityInfo; got != want {
		t.Fatalf("Severity() = %q, want %q", got, want)
	}

	if got, want := event.Visibility(), VisibilityPublic; got != want {
		t.Fatalf("Visibility() = %q, want %q", got, want)
	}
}

// TestNewEventNormalizesMessage verifies event message normalization.
func TestNewEventNormalizesMessage(t *testing.T) {
	t.Parallel()

	event := MustEvent(EventSpec{
		Kind:    EventKindDiagnostic,
		Message: "  First line.  \r\n  Second line.  ",
	})

	if got, want := event.Message(), "First line.\nSecond line."; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}
}

// TestNewSimpleEvent verifies simple event constructor.
func TestNewSimpleEvent(t *testing.T) {
	t.Parallel()

	event, err := NewSimpleEvent(EventKindDiagnostic, "Diagnostic event.")
	if err != nil {
		t.Fatalf("NewSimpleEvent() returned unexpected error: %v", err)
	}

	if got, want := event.Kind(), EventKindDiagnostic; got != want {
		t.Fatalf("Kind() = %q, want %q", got, want)
	}

	if got, want := event.Message(), "Diagnostic event."; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}
}
