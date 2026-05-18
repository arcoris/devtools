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

func TestNewGroupRejectsInvalidGroups(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		err  error
	}{
		{name: "empty", raw: "", err: ErrEmptyGroup},
		{name: "leading separator", raw: ".benchmark", err: ErrInvalidGroup},
		{name: "trailing separator", raw: "benchmark.", err: ErrInvalidGroup},
		{name: "consecutive separators", raw: "diagnostics..perf", err: ErrInvalidGroup},
		{name: "invalid segment", raw: "diagnostics.Perf", err: ErrInvalidGroup},
		{name: "too long complete group", raw: "a." + strings.Repeat("b", maxGroupLength), err: ErrInvalidGroup},
		{name: "too deep", raw: strings.Repeat("a.", maxGroupDepth) + "a", err: ErrInvalidGroup},
		{name: "too long segment", raw: "a" + strings.Repeat("b", maxGroupSegmentLength), err: ErrInvalidGroup},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			group, err := NewGroup(test.raw)
			if err == nil {
				t.Fatalf("NewGroup(%q) returned nil error and group %q", test.raw, group)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("NewGroup(%q) error = %v, want errors.Is(..., %v)", test.raw, err, test.err)
			}

			if Group(test.raw).IsValid() {
				t.Fatalf("Group(%q).IsValid() = true, want false", test.raw)
			}
		})
	}
}

func TestValidateGroupSegmentWrapsGenericValidation(t *testing.T) {
	t.Parallel()

	err := validateGroupSegment(0, "Invalid")
	if err == nil {
		t.Fatalf("validateGroupSegment() returned nil error")
	}

	if !errors.Is(err, ErrInvalidGroup) {
		t.Fatalf("validateGroupSegment() error = %v, want ErrInvalidGroup", err)
	}

	if !errors.Is(err, ErrInvalidCommandNameSegment) {
		t.Fatalf("validateGroupSegment() error = %v, want ErrInvalidCommandNameSegment", err)
	}
}

func TestValidateGroupSegmentWrapsEmptySegment(t *testing.T) {
	t.Parallel()

	err := validateGroupSegment(1, "")
	if err == nil {
		t.Fatalf("validateGroupSegment() returned nil error")
	}

	if !errors.Is(err, ErrInvalidGroup) {
		t.Fatalf("validateGroupSegment() error = %v, want ErrInvalidGroup", err)
	}

	if !errors.Is(err, ErrEmptyCommandNameSegment) {
		t.Fatalf("validateGroupSegment() error = %v, want ErrEmptyCommandNameSegment", err)
	}
}

func TestNewGroupPartsRejectsNoSegments(t *testing.T) {
	t.Parallel()

	_, err := NewGroupParts()
	if err == nil {
		t.Fatalf("NewGroupParts() returned nil error")
	}

	if !errors.Is(err, ErrEmptyGroup) {
		t.Fatalf("NewGroupParts() error = %v, want ErrEmptyGroup", err)
	}
}

func TestCanonicalGroupLength(t *testing.T) {
	t.Parallel()

	if got, want := canonicalGroupLength([]string{"diagnostics", "perf"}), len("diagnostics.perf"); got != want {
		t.Fatalf("canonicalGroupLength() = %d, want %d", got, want)
	}
}
