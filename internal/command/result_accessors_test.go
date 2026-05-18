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

// TestResultArtifactIDs verifies artifact ID ordering helpers.
func TestResultArtifactIDs(t *testing.T) {
	t.Parallel()

	result := MustResult(ResultSpec{
		Artifacts: []Artifact{
			resultTestArtifact("z.report"),
			resultTestArtifact("a.report"),
			resultTestArtifact("m.report"),
		},
	})

	ids := result.SortedArtifactIDs()
	if got, want := ids[0], MustArtifactID("a.report"); got != want {
		t.Fatalf("SortedArtifactIDs()[0] = %q, want %q", got, want)
	}
}

// TestResultCopySemantics verifies detached slices and maps.
func TestResultCopySemantics(t *testing.T) {
	t.Parallel()

	artifacts := []Artifact{resultTestArtifact("bench.report")}
	warnings := []ResultWarning{
		MustResultWarning(ResultWarningSpec{
			Kind:    "partial",
			Message: "Partial result.",
		}),
	}
	fields := map[string]string{"mode": "ci"}

	result := MustResult(ResultSpec{
		Artifacts: artifacts,
		Warnings:  warnings,
		Fields:    fields,
	})

	artifacts[0] = resultTestArtifact("changed")
	warnings[0] = MustResultWarning(ResultWarningSpec{
		Kind:    "changed",
		Message: "Changed.",
	})
	fields["mode"] = "changed"

	if result.HasArtifact(MustArtifactID("changed")) {
		t.Fatalf("result changed through input artifact slice")
	}

	if got, want := result.Warnings()[0].Kind(), "partial"; got != want {
		t.Fatalf("warning changed through input slice: got %q, want %q", got, want)
	}

	if got, want := resultTestMustField(t, result, "mode"), "ci"; got != want {
		t.Fatalf("field changed through input map: got %q, want %q", got, want)
	}

	outArtifacts := result.Artifacts()
	outArtifacts[0] = resultTestArtifact("changed")

	if result.HasArtifact(MustArtifactID("changed")) {
		t.Fatalf("result changed through output artifact slice")
	}

	outFields := result.Fields()
	outFields["mode"] = "changed"

	if got, want := resultTestMustField(t, result, "mode"), "ci"; got != want {
		t.Fatalf("field changed through output map: got %q, want %q", got, want)
	}
}

// TestResultFieldKeys verifies deterministic field key ordering.
func TestResultFieldKeys(t *testing.T) {
	t.Parallel()

	result := MustResult(ResultSpec{
		Fields: map[string]string{
			"z.key": "z",
			"a.key": "a",
			"m.key": "m",
		},
	})

	resultTestAssertStrings(t, result.FieldKeys(), []string{"a.key", "m.key", "z.key"})
}

// TestResultWithArtifactReplacesByID verifies artifact upsert behavior.
func TestResultWithArtifactReplacesByID(t *testing.T) {
	t.Parallel()

	result := OKResult("ok").
		MustWithArtifact(resultTestArtifact("bench.report")).
		MustWithArtifact(MustArtifact(ArtifactSpec{
			ID:          "bench.report",
			Kind:        ArtifactKindReport,
			Location:    "reports/new.md",
			Description: "New report.",
		}))

	if got, want := len(result.Artifacts()), 1; got != want {
		t.Fatalf("len(Artifacts()) = %d, want %d", got, want)
	}

	artifact, ok := result.Artifact(MustArtifactID("bench.report"))
	if !ok {
		t.Fatalf("Artifact(bench.report) ok = false, want true")
	}

	if got, want := artifact.Location(), "reports/new.md"; got != want {
		t.Fatalf("artifact Location() = %q, want %q", got, want)
	}
}
