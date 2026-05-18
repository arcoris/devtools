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
)

// TestEventIDValidation verifies EventID value object behavior.
func TestEventIDValidation(t *testing.T) {
	t.Parallel()

	id, err := NewEventID("command.started.001")
	if err != nil {
		t.Fatalf("NewEventID() returned unexpected error: %v", err)
	}

	if got, want := id.String(), "command.started.001"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	invalid := []struct {
		raw string
		err error
	}{
		{raw: "", err: ErrEmptyEventID},
		{raw: "Command.Started", err: ErrInvalidEventID},
		{raw: ".command", err: ErrInvalidEventID},
		{raw: "command.", err: ErrInvalidEventID},
		{raw: "command..started", err: ErrInvalidEventID},
		{raw: strings.Repeat("x", maxEventIDLength+1), err: ErrInvalidEventID},
	}

	for _, tt := range invalid {
		tt := tt

		t.Run(tt.raw, func(t *testing.T) {
			t.Parallel()

			_, err := NewEventID(tt.raw)
			if err == nil {
				t.Fatalf("NewEventID(%q) returned nil error", tt.raw)
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewEventID(%q) error = %v, want %v", tt.raw, err, tt.err)
			}
		})
	}
}

// TestMustEventIDPanicsForInvalidID verifies fail-fast ID construction.
func TestMustEventIDPanicsForInvalidID(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustEventID did not panic")
		}
	}()

	_ = MustEventID("Bad")
}
