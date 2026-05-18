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

// TestNormalizeRuntimePanicMessage verifies panic message normalization.
func TestNormalizeRuntimePanicMessage(t *testing.T) {
	t.Parallel()

	if got, want := normalizeRuntimePanicMessage("  boom  "), "boom"; got != want {
		t.Fatalf("normalizeRuntimePanicMessage() = %q, want %q", got, want)
	}

	if got, want := normalizeRuntimePanicMessage(""), "panic"; got != want {
		t.Fatalf("normalizeRuntimePanicMessage(empty) = %q, want %q", got, want)
	}

	long := strings.Repeat("x", maxRuntimePanicMessageLength+1)
	if got := normalizeRuntimePanicMessage(long); len(got) != maxRuntimePanicMessageLength {
		t.Fatalf("normalized long panic length = %d, want %d", len(got), maxRuntimePanicMessageLength)
	}
}

// TestRuntimeValidationHelpers verifies runtime helper validation.
func TestRuntimeValidationHelpers(t *testing.T) {
	t.Parallel()

	if err := validateRuntimeName("runtime"); err != nil {
		t.Fatalf("validateRuntimeName(valid) returned unexpected error: %v", err)
	}

	invalid := []string{
		"",
		" bad",
		"bad\nname",
		"bad\x00name",
		strings.Repeat("x", maxRuntimeNameLength+1),
	}

	for _, raw := range invalid {
		raw := raw

		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			err := validateRuntimeName(raw)
			if err == nil {
				t.Fatalf("validateRuntimeName(%q) returned nil error", raw)
			}

			if !errors.Is(err, ErrInvalidRuntime) {
				t.Fatalf("validateRuntimeName(%q) error = %v, want ErrInvalidRuntime", raw, err)
			}
		})
	}

	if err := validateRuntimeName("bad\x00name"); !errors.Is(err, textvalidate.ErrInvalidSingleLineText) {
		t.Fatalf("validateRuntimeName(control) error = %v, want ErrInvalidSingleLineText", err)
	}
}
