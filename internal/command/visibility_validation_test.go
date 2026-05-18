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
	"testing"
)

func TestNewVisibilityRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		err  error
	}{
		{name: "empty", raw: "", err: ErrEmptyVisibility},
		{name: "uppercase", raw: "Public", err: ErrInvalidVisibility},
		{name: "unknown", raw: "private", err: ErrInvalidVisibility},
		{name: "space", raw: "public hidden", err: ErrInvalidVisibility},
		{name: "underscore", raw: "internal_only", err: ErrInvalidVisibility},
		{name: "hyphen", raw: "internal-only", err: ErrInvalidVisibility},
		{name: "whitespace", raw: " public ", err: ErrInvalidVisibility},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewVisibility(test.raw)
			if err == nil {
				t.Fatalf("NewVisibility(%q) returned nil error and value %q", test.raw, got)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("NewVisibility(%q) error = %v, want errors.Is(..., %v)", test.raw, err, test.err)
			}

			if Visibility(test.raw).IsValid() {
				t.Fatalf("Visibility(%q).IsValid() = true, want false", test.raw)
			}
		})
	}
}

func TestVisibilityValidateRejectsZero(t *testing.T) {
	t.Parallel()

	err := Visibility("").Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrEmptyVisibility) {
		t.Fatalf("Validate() error = %v, want ErrEmptyVisibility", err)
	}
}

func TestVisibilityValidateRejectsUnknown(t *testing.T) {
	t.Parallel()

	err := Visibility("private").Validate()
	if err == nil {
		t.Fatalf("Validate() returned nil error")
	}

	if !errors.Is(err, ErrInvalidVisibility) {
		t.Fatalf("Validate() error = %v, want ErrInvalidVisibility", err)
	}
}
