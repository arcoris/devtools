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

func TestNewActionArtifactAcceptsValidArtifact(t *testing.T) {
	t.Parallel()

	artifact, err := NewActionArtifact("report", "bench/reports/check.md", "Check report")
	if err != nil {
		t.Fatalf("NewActionArtifact() returned unexpected error: %v", err)
	}

	if got, want := artifact.Kind, "report"; got != want {
		t.Fatalf("Kind = %q, want %q", got, want)
	}

	if (ActionArtifact{}).IsZero() == false {
		t.Fatalf("zero artifact IsZero() = false, want true")
	}
}

func TestActionArtifactRejectsInvalidArtifact(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		artifact ActionArtifact
	}{
		{name: "empty kind", artifact: ActionArtifact{Kind: "", Path: "x"}},
		{name: "invalid kind", artifact: ActionArtifact{Kind: "Report", Path: "x"}},
		{name: "empty path", artifact: ActionArtifact{Kind: "report", Path: ""}},
		{name: "blank path", artifact: ActionArtifact{Kind: "report", Path: "   "}},
		{name: "path control", artifact: ActionArtifact{Kind: "report", Path: "bad\x00path"}},
		{name: "path too long", artifact: ActionArtifact{Kind: "report", Path: strings.Repeat("x", maxActionArtifactPathLength+1)}},
		{name: "description control", artifact: ActionArtifact{Kind: "report", Path: "x", Description: "bad\x00description"}},
		{name: "description too long", artifact: ActionArtifact{Kind: "report", Path: "x", Description: strings.Repeat("x", maxActionDescriptionLength+1)}},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.artifact.Validate()
			if err == nil {
				t.Fatalf("Validate() returned nil error")
			}

			if !errors.Is(err, ErrInvalidActionResult) {
				t.Fatalf("Validate() error = %v, want ErrInvalidActionResult", err)
			}
		})
	}
}

func TestActionArtifactWrapsReusableValidatorErrors(t *testing.T) {
	t.Parallel()

	err := ActionArtifact{Kind: "Report", Path: "x"}.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, textvalidate.ErrInvalidDottedKebabKey) {
		t.Fatalf("Validate() error = %v, want ErrInvalidDottedKebabKey", err)
	}

	err = ActionArtifact{Kind: "report", Path: "bad\x00path"}.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("Validate() error = %v, want ErrInvalidCompactText", err)
	}
}

func TestMustActionArtifactPanicsForInvalidArtifact(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustActionArtifact("", "x", "")
	})
}
