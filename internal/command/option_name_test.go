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

// TestMustOptionNamePanicsForInvalidName verifies fail-fast option-name construction.
func TestMustOptionNamePanicsForInvalidName(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustOptionName did not panic")
		}
	}()

	_ = MustOptionName("Bad")
}

// TestOptionNameValidation verifies option-name value object behavior.
func TestOptionNameValidation(t *testing.T) {
	t.Parallel()

	name, err := NewOptionName("bench-time")
	if err != nil {
		t.Fatalf("NewOptionName() returned unexpected error: %v", err)
	}

	if got, want := name.String(), "bench-time"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	if got, want := name.LongFlag(), "--bench-time"; got != want {
		t.Fatalf("LongFlag() = %q, want %q", got, want)
	}

	invalid := []struct {
		raw string
		err error
	}{
		{raw: "", err: ErrEmptyOptionName},
		{raw: "Bench", err: ErrInvalidOptionName},
		{raw: "bench_time", err: ErrInvalidOptionName},
		{raw: "1bench", err: ErrInvalidOptionName},
		{raw: strings.Repeat("x", maxOptionNameLength+1), err: ErrInvalidOptionName},
	}

	for _, tt := range invalid {
		tt := tt

		t.Run(tt.raw, func(t *testing.T) {
			t.Parallel()

			_, err := NewOptionName(tt.raw)
			if err == nil {
				t.Fatalf("NewOptionName(%q) returned nil error", tt.raw)
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewOptionName(%q) error = %v, want %v", tt.raw, err, tt.err)
			}
		})
	}
}
