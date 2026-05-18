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

// TestResultWithHelpers verifies immutable-style result updates.
func TestResultWithHelpers(t *testing.T) {
	t.Parallel()

	startedAt := time.Date(2026, 5, 18, 10, 0, 0, 0, time.UTC)
	finishedAt := startedAt.Add(time.Second)

	result := OKResult("ok").
		MustWithStatus(ResultStatusSkipped).
		MustWithMessage("Skipped.").
		MustWithTiming(startedAt, finishedAt).
		MustWithExitCode(0).
		MustWithArtifact(resultTestArtifact("bench.report")).
		MustWithWarning(MustResultWarning(ResultWarningSpec{
			Kind:    "partial",
			Message: "Partial.",
		})).
		MustWithField("mode", "ci").
		MustWithMetadata(MustMetadata(MetadataSpec{Owner: "devtools"})).
		MustWithVisibility(VisibilityHidden)

	if !result.IsSkipped() {
		t.Fatalf("IsSkipped() = false, want true")
	}

	if got, want := result.Message(), "Skipped."; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}

	if !result.HasArtifact(MustArtifactID("bench.report")) {
		t.Fatalf("HasArtifact(bench.report) = false, want true")
	}

	if got, want := result.WarningCount(), 1; got != want {
		t.Fatalf("WarningCount() = %d, want %d", got, want)
	}

	if got, want := result.Metadata().Owner(), "devtools"; got != want {
		t.Fatalf("Metadata().Owner() = %q, want %q", got, want)
	}

	if !result.Visibility().IsHidden() {
		t.Fatalf("Visibility().IsHidden() = false, want true")
	}

	withoutArtifact := result.WithoutArtifact(MustArtifactID("bench.report"))
	if withoutArtifact.HasArtifact(MustArtifactID("bench.report")) {
		t.Fatalf("WithoutArtifact() still has artifact")
	}

	withoutField := result.WithoutField("mode")
	if withoutField.HasField("mode") {
		t.Fatalf("WithoutField() still has field")
	}
}
