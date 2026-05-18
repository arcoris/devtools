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

package textvalidate

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateDottedKebabKeyAcceptsValidKeys(t *testing.T) {
	t.Parallel()

	tests := []string{
		"owner",
		"command.registry",
		"devtools-team",
		"a1",
		"a.b.c",
	}

	for _, key := range tests {
		key := key

		t.Run(key, func(t *testing.T) {
			t.Parallel()

			if err := ValidateDottedKebabKey(key, 255, 32); err != nil {
				t.Fatalf("ValidateDottedKebabKey(%q) returned unexpected error: %v", key, err)
			}
		})
	}
}

func TestValidateDottedKebabKeyRejectsInvalidKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		key  string
		err  error
	}{
		{name: "empty", key: "", err: ErrEmptyDottedKebabKey},
		{name: "leading separator", key: ".owner", err: ErrInvalidDottedKebabKey},
		{name: "trailing separator", key: "owner.", err: ErrInvalidDottedKebabKey},
		{name: "empty segment", key: "owner..team", err: ErrInvalidDottedKebabKey},
		{name: "invalid segment", key: "Owner", err: ErrInvalidDottedKebabKey},
		{name: "too long", key: strings.Repeat("a", 4), err: ErrInvalidDottedKebabKey},
		{name: "too deep", key: "a.b.c", err: ErrInvalidDottedKebabKey},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateDottedKebabKey(test.key, 3, 2)
			if err == nil {
				t.Fatalf("ValidateDottedKebabKey(%q) returned nil error", test.key)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("ValidateDottedKebabKey(%q) error = %v, want %v", test.key, err, test.err)
			}
		})
	}
}

func TestValidateDottedKebabKeyAllowsUnlimitedDepth(t *testing.T) {
	t.Parallel()

	if err := ValidateDottedKebabKey("a.b.c.d", 255, 0); err != nil {
		t.Fatalf("ValidateDottedKebabKey() returned unexpected error: %v", err)
	}
}

func TestValidateDottedKebabKeyWrapsSegmentError(t *testing.T) {
	t.Parallel()

	err := ValidateDottedKebabKey("owner.Bad", 255, 32)
	if err == nil {
		t.Fatalf("ValidateDottedKebabKey() returned nil error")
	}

	if !errors.Is(err, ErrInvalidDottedKebabKey) {
		t.Fatalf("ValidateDottedKebabKey() error = %v, want ErrInvalidDottedKebabKey", err)
	}

	if !errors.Is(err, ErrInvalidKebabSegment) {
		t.Fatalf("ValidateDottedKebabKey() error = %v, want ErrInvalidKebabSegment", err)
	}
}

func TestValidateDottedKebabKeyRejectsInvalidLimits(t *testing.T) {
	t.Parallel()

	err := ValidateDottedKebabKey("owner", 0, 0)
	if err == nil {
		t.Fatalf("ValidateDottedKebabKey() returned nil error")
	}

	if !errors.Is(err, ErrInvalidDottedKebabKey) {
		t.Fatalf("ValidateDottedKebabKey() error = %v, want ErrInvalidDottedKebabKey", err)
	}
}
