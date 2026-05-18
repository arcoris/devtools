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

// TestArgumentNameValidation verifies argument-name value object behavior.
func TestArgumentNameValidation(t *testing.T) {
	t.Parallel()

	name, err := NewArgumentName("bench-time")
	if err != nil {
		t.Fatalf("NewArgumentName() returned unexpected error: %v", err)
	}

	if got, want := name.String(), "bench-time"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	invalid := []struct {
		raw string
		err error
	}{
		{raw: "", err: ErrEmptyArgumentName},
		{raw: "Bench", err: ErrInvalidArgumentName},
		{raw: "bench_time", err: ErrInvalidArgumentName},
		{raw: "1bench", err: ErrInvalidArgumentName},
		{raw: strings.Repeat("x", maxArgumentNameLength+1), err: ErrInvalidArgumentName},
	}

	for _, tt := range invalid {
		tt := tt

		t.Run(tt.raw, func(t *testing.T) {
			t.Parallel()

			_, err := NewArgumentName(tt.raw)
			if err == nil {
				t.Fatalf("NewArgumentName(%q) returned nil error", tt.raw)
			}

			if !errors.Is(err, tt.err) {
				t.Fatalf("NewArgumentName(%q) error = %v, want %v", tt.raw, err, tt.err)
			}
		})
	}
}

// TestMustArgumentNamePanicsForInvalidName verifies fail-fast name construction.
func TestMustArgumentNamePanicsForInvalidName(t *testing.T) {
	t.Parallel()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("MustArgumentName did not panic")
		}
	}()

	_ = MustArgumentName("Bad")
}
