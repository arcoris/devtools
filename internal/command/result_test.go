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

func resultTestArtifact(id string) Artifact {
	return MustArtifact(ArtifactSpec{
		ID:       id,
		Kind:     ArtifactKindReport,
		Location: id + ".md",
	})
}

func resultTestMustField(t *testing.T, result Result, key string) string {
	t.Helper()

	value, ok := result.Field(key)
	if !ok {
		t.Fatalf("field %q not found", key)
	}

	return value
}

func resultTestAssertStrings(t *testing.T, got []string, want []string) {
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

// TestNewResultAcceptsValidResult verifies full result construction.
func TestNewResultAcceptsValidResult(t *testing.T) {
	t.Parallel()

	startedAt := time.Date(2026, 5, 18, 10, 0, 0, 0, time.UTC)
	finishedAt := startedAt.Add(2 * time.Second)
	exitCode := 0

	result, err := NewResult(ResultSpec{
		Status:     ResultStatusOK,
		Message:    "Command completed.",
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		ExitCode:   &exitCode,
		Artifacts: []Artifact{
			resultTestArtifact("bench.report"),
		},
		Warnings: []ResultWarning{
			MustResultWarning(ResultWarningSpec{
				Kind:    "partial",
				Message: "Optional profile was not collected.",
				Hints:   []string{"Run with profiling enabled."},
			}),
		},
		Fields: map[string]string{
			"mode": "ci",
		},
		Metadata: MustMetadata(MetadataSpec{
			Owner: "devtools",
		}),
		Visibility: VisibilityPublic,
	})
	if err != nil {
		t.Fatalf("NewResult() returned unexpected error: %v", err)
	}

	if got, want := result.Status(), ResultStatusOK; got != want {
		t.Fatalf("Status() = %q, want %q", got, want)
	}

	if got, want := result.Message(), "Command completed."; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}

	if got, ok := result.Duration(); !ok || got != 2*time.Second {
		t.Fatalf("Duration() = %v, %v; want 2s, true", got, ok)
	}

	if got, ok := result.ExitCode(); !ok || got != 0 {
		t.Fatalf("ExitCode() = %d, %v; want 0, true", got, ok)
	}

	if !result.HasArtifact(MustArtifactID("bench.report")) {
		t.Fatalf("HasArtifact(bench.report) = false, want true")
	}

	if got, want := result.WarningCount(), 1; got != want {
		t.Fatalf("WarningCount() = %d, want %d", got, want)
	}

	if got, ok := result.Field("mode"); !ok || got != "ci" {
		t.Fatalf("Field(mode) = %q, %v; want ci, true", got, ok)
	}
}

// TestNewResultDefaults verifies default status and visibility.
func TestNewResultDefaults(t *testing.T) {
	t.Parallel()

	result, err := NewResult(ResultSpec{})
	if err != nil {
		t.Fatalf("NewResult() returned unexpected error: %v", err)
	}

	if got, want := result.Status(), ResultStatusOK; got != want {
		t.Fatalf("Status() = %q, want %q", got, want)
	}

	if got, want := result.Visibility(), VisibilityPublic; got != want {
		t.Fatalf("Visibility() = %q, want %q", got, want)
	}

	if got, want := result.RecommendedExitCode(), 0; got != want {
		t.Fatalf("RecommendedExitCode() = %d, want %d", got, want)
	}
}

// TestNewResultNormalizesMessage verifies result message normalization.
func TestNewResultNormalizesMessage(t *testing.T) {
	t.Parallel()

	result := MustResult(ResultSpec{
		Message: "  First line.  \r\n  Second line.  ",
	})

	if got, want := result.Message(), "First line.\nSecond line."; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}
}

// TestResultConstructors verifies convenience constructors.
func TestResultConstructors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		result   Result
		status   ResultStatus
		exitCode int
	}{
		{
			name:     "ok",
			result:   OKResult("done"),
			status:   ResultStatusOK,
			exitCode: 0,
		},
		{
			name:     "skipped",
			result:   SkippedResult("skipped"),
			status:   ResultStatusSkipped,
			exitCode: 0,
		},
		{
			name:     "failed",
			result:   FailedResult("failed"),
			status:   ResultStatusFailed,
			exitCode: 1,
		},
		{
			name:     "canceled",
			result:   CanceledResult("canceled"),
			status:   ResultStatusCanceled,
			exitCode: 130,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.result.Status(); got != tt.status {
				t.Fatalf("Status() = %q, want %q", got, tt.status)
			}

			if got := tt.result.RecommendedExitCode(); got != tt.exitCode {
				t.Fatalf("RecommendedExitCode() = %d, want %d", got, tt.exitCode)
			}
		})
	}
}
