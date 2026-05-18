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

func TestNewActionWarningAcceptsValidWarning(t *testing.T) {
	t.Parallel()

	warning, err := NewActionWarning("partial", "Some optional checks were skipped.")
	if err != nil {
		t.Fatalf("NewActionWarning() returned unexpected error: %v", err)
	}

	if got, want := warning.Kind, "partial"; got != want {
		t.Fatalf("Kind = %q, want %q", got, want)
	}

	if (ActionWarning{}).IsZero() == false {
		t.Fatalf("zero warning IsZero() = false, want true")
	}
}

func TestActionWarningRejectsInvalidWarning(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		warning ActionWarning
	}{
		{name: "empty kind", warning: ActionWarning{Kind: "", Message: "message"}},
		{name: "invalid kind", warning: ActionWarning{Kind: "Partial", Message: "message"}},
		{name: "empty message", warning: ActionWarning{Kind: "partial", Message: ""}},
		{name: "blank message", warning: ActionWarning{Kind: "partial", Message: "   "}},
		{name: "message control", warning: ActionWarning{Kind: "partial", Message: "bad\x00message"}},
		{name: "message too long", warning: ActionWarning{Kind: "partial", Message: strings.Repeat("x", maxActionMessageLength+1)}},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.warning.Validate()
			if err == nil {
				t.Fatalf("Validate() returned nil error")
			}

			if !errors.Is(err, ErrInvalidActionResult) {
				t.Fatalf("Validate() error = %v, want ErrInvalidActionResult", err)
			}
		})
	}
}

func TestActionWarningWrapsReusableValidatorErrors(t *testing.T) {
	t.Parallel()

	err := ActionWarning{Kind: "Partial", Message: "message"}.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, textvalidate.ErrInvalidDottedKebabKey) {
		t.Fatalf("Validate() error = %v, want ErrInvalidDottedKebabKey", err)
	}

	err = ActionWarning{Kind: "partial", Message: "bad\x00message"}.Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("Validate() error = %v, want ErrInvalidCompactText", err)
	}
}

func TestMustActionWarningPanicsForInvalidWarning(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustActionWarning("", "message")
	})
}
