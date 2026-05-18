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
	"errors"
	"strings"
	"testing"

	"arcoris.dev/devtools/internal/textvalidate"
)

func TestNewActionResultNormalizesValidResult(t *testing.T) {
	t.Parallel()

	result, err := NewActionResult(ActionResult{
		Message: "done",
		Fields:  map[string]string{"mode": "ci"},
	})
	if err != nil {
		t.Fatalf("NewActionResult() returned unexpected error: %v", err)
	}

	if got, want := result.Status, ActionStatusOK; got != want {
		t.Fatalf("Status = %q, want %q", got, want)
	}

	if got, want := actionTestResultField(t, result, "mode"), "ci"; got != want {
		t.Fatalf("Field(mode) = %q, want %q", got, want)
	}
}

func TestNewActionResultRejectsInvalidResult(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body ActionResult
	}{
		{name: "unknown status", body: ActionResult{Status: ActionStatus("unknown")}},
		{name: "message control", body: ActionResult{Message: "bad\x00value"}},
		{name: "message too long", body: ActionResult{Message: strings.Repeat("x", maxActionMessageLength+1)}},
		{name: "invalid artifact", body: ActionResult{Artifacts: []ActionArtifact{{Kind: "", Path: "x"}}}},
		{name: "invalid warning", body: ActionResult{Warnings: []ActionWarning{{Kind: "", Message: "x"}}}},
		{name: "invalid field key", body: ActionResult{Fields: map[string]string{"Bad": "value"}}},
		{name: "invalid field value", body: ActionResult{Fields: map[string]string{"mode": "bad\x00value"}}},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewActionResult(test.body)
			if err == nil {
				t.Fatalf("NewActionResult() returned nil error")
			}

			if !errors.Is(err, ErrInvalidActionResult) {
				t.Fatalf("NewActionResult() error = %v, want ErrInvalidActionResult", err)
			}
		})
	}
}

func TestNewActionResultWrapsTextvalidateErrors(t *testing.T) {
	t.Parallel()

	_, err := NewActionResult(ActionResult{Message: "bad\x00value"})
	if err == nil {
		t.Fatalf("NewActionResult() returned nil error")
	}

	if !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("NewActionResult() error = %v, want ErrInvalidActionResult", err)
	}

	if !errors.Is(err, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("NewActionResult() error = %v, want ErrInvalidCompactText", err)
	}
}

func TestActionResultAccessors(t *testing.T) {
	t.Parallel()

	result := MustActionResult(ActionResult{
		Status: ActionStatusSkipped,
		Data:   map[string]int{"count": 1},
		Artifacts: []ActionArtifact{
			MustActionArtifact("report", "report.md", "Report"),
		},
		Warnings: []ActionWarning{
			MustActionWarning("partial", "Partial result."),
		},
		Fields: map[string]string{
			"z.key": "z",
			"a.key": "a",
		},
	})

	if !result.IsSkipped() {
		t.Fatalf("IsSkipped() = false, want true")
	}

	if !result.HasData() {
		t.Fatalf("HasData() = false, want true")
	}

	if !result.HasArtifacts() || result.ArtifactCount() != 1 {
		t.Fatalf("artifact accessors returned unexpected values")
	}

	if !result.HasWarnings() || result.WarningCount() != 1 {
		t.Fatalf("warning accessors returned unexpected values")
	}

	if !result.HasFields() || result.FieldCount() != 2 {
		t.Fatalf("field accessors returned unexpected values")
	}

	if !result.HasField("a.key") {
		t.Fatalf("HasField(a.key) = false, want true")
	}

	assertStringSlicesEqual(t, result.FieldKeys(), []string{"a.key", "z.key"})
}

func TestActionResultCopySemantics(t *testing.T) {
	t.Parallel()

	artifacts := []ActionArtifact{
		MustActionArtifact("report", "report.md", "Report"),
	}
	warnings := []ActionWarning{
		MustActionWarning("partial", "Partial result."),
	}
	fields := map[string]string{"mode": "ci"}

	result := ActionResult{
		Artifacts: artifacts,
		Warnings:  warnings,
		Fields:    fields,
	}.Normalize()

	artifacts[0] = MustActionArtifact("profile", "cpu.pprof", "CPU profile")
	warnings[0] = MustActionWarning("other", "Other warning.")
	fields["mode"] = "changed"

	if got, want := result.ArtifactRefs()[0].Kind, "report"; got != want {
		t.Fatalf("artifact changed through input slice: got %q, want %q", got, want)
	}

	if got, want := result.WarningRefs()[0].Kind, "partial"; got != want {
		t.Fatalf("warning changed through input slice: got %q, want %q", got, want)
	}

	if got, want := actionTestResultField(t, result, "mode"), "ci"; got != want {
		t.Fatalf("field changed through input map: got %q, want %q", got, want)
	}

	outArtifacts := result.ArtifactRefs()
	outArtifacts[0] = MustActionArtifact("trace", "trace.out", "Trace")

	if got, want := result.ArtifactRefs()[0].Kind, "report"; got != want {
		t.Fatalf("artifact changed through ArtifactRefs: got %q, want %q", got, want)
	}

	outWarnings := result.WarningRefs()
	outWarnings[0] = MustActionWarning("other", "Other warning.")

	if got, want := result.WarningRefs()[0].Kind, "partial"; got != want {
		t.Fatalf("warning changed through WarningRefs: got %q, want %q", got, want)
	}

	outFields := result.FieldMap()
	outFields["mode"] = "changed"

	if got, want := actionTestResultField(t, result, "mode"), "ci"; got != want {
		t.Fatalf("field changed through FieldMap: got %q, want %q", got, want)
	}
}

func TestActionResultWithHelpers(t *testing.T) {
	t.Parallel()

	result := ActionResult{}.
		MustWithStatus(ActionStatusSkipped).
		MustWithMessage("skipped").
		MustWithData(map[string]int{"count": 1}).
		MustWithArtifact(MustActionArtifact("report", "bench/reports/check.md", "Check report")).
		MustWithWarning(MustActionWarning("partial", "Some optional tools were skipped.")).
		MustWithField("mode", "ci")

	if !result.IsSkipped() {
		t.Fatalf("IsSkipped() = false, want true")
	}

	if got, want := result.Message, "skipped"; got != want {
		t.Fatalf("Message = %q, want %q", got, want)
	}

	if !result.HasData() {
		t.Fatalf("HasData() = false, want true")
	}

	if got, want := result.ArtifactCount(), 1; got != want {
		t.Fatalf("ArtifactCount() = %d, want %d", got, want)
	}

	if got, want := result.WarningCount(), 1; got != want {
		t.Fatalf("WarningCount() = %d, want %d", got, want)
	}

	if got, want := actionTestResultField(t, result, "mode"), "ci"; got != want {
		t.Fatalf("Field(mode) = %q, want %q", got, want)
	}

	if result.WithoutArtifacts().HasArtifacts() {
		t.Fatalf("WithoutArtifacts() still has artifacts")
	}

	if result.WithoutWarnings().HasWarnings() {
		t.Fatalf("WithoutWarnings() still has warnings")
	}

	if result.WithoutField("mode").HasField("mode") {
		t.Fatalf("WithoutField() still has mode")
	}

	if result.WithoutFields().HasFields() {
		t.Fatalf("WithoutFields() still has fields")
	}
}

func TestActionResultWithCollectionsCopiesInput(t *testing.T) {
	t.Parallel()

	artifacts := []ActionArtifact{MustActionArtifact("report", "report.md", "Report")}
	warnings := []ActionWarning{MustActionWarning("partial", "Partial result.")}
	fields := map[string]string{"mode": "ci"}

	result := ActionResult{}.
		MustWithArtifacts(artifacts).
		MustWithWarnings(warnings).
		MustWithFields(fields)

	artifacts[0] = MustActionArtifact("profile", "cpu.pprof", "CPU profile")
	warnings[0] = MustActionWarning("other", "Other warning.")
	fields["mode"] = "changed"

	if got, want := result.ArtifactRefs()[0].Kind, "report"; got != want {
		t.Fatalf("artifact changed through WithArtifacts input: got %q, want %q", got, want)
	}

	if got, want := result.WarningRefs()[0].Kind, "partial"; got != want {
		t.Fatalf("warning changed through WithWarnings input: got %q, want %q", got, want)
	}

	if got, want := actionTestResultField(t, result, "mode"), "ci"; got != want {
		t.Fatalf("field changed through WithFields input: got %q, want %q", got, want)
	}
}

func TestActionResultWithHelpersRejectInvalidValues(t *testing.T) {
	t.Parallel()

	result := ActionResult{}

	if _, err := result.WithStatus(ActionStatus("unknown")); !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("WithStatus() error = %v, want ErrInvalidActionResult", err)
	}

	if _, err := result.WithMessage("bad\x00value"); !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("WithMessage() error = %v, want ErrInvalidActionResult", err)
	}

	if _, err := result.WithArtifact(ActionArtifact{}); !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("WithArtifact() error = %v, want ErrInvalidActionResult", err)
	}

	if _, err := result.WithArtifacts([]ActionArtifact{{Kind: "", Path: "x"}}); !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("WithArtifacts() error = %v, want ErrInvalidActionResult", err)
	}

	if _, err := result.WithWarning(ActionWarning{}); !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("WithWarning() error = %v, want ErrInvalidActionResult", err)
	}

	if _, err := result.WithWarnings([]ActionWarning{{Kind: "", Message: "x"}}); !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("WithWarnings() error = %v, want ErrInvalidActionResult", err)
	}

	if _, err := result.WithField("Bad", "value"); !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("WithField(key) error = %v, want ErrInvalidActionResult", err)
	}

	if _, err := result.WithFields(map[string]string{"mode": "bad\x00value"}); !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("WithFields(value) error = %v, want ErrInvalidActionResult", err)
	}
}

func TestMustActionResultPanicsForInvalidResult(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustActionResult(ActionResult{Status: ActionStatus("unknown")})
	})
}
