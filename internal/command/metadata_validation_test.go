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

func TestValidateMetadataKey(t *testing.T) {
	t.Parallel()

	valid := []string{"owner", "command", "command.registry", "devtools-team", "a1"}
	for _, raw := range valid {
		raw := raw

		t.Run("valid "+raw, func(t *testing.T) {
			t.Parallel()

			if err := validateMetadataKey("field", raw); err != nil {
				t.Fatalf("validateMetadataKey(%q) returned unexpected error: %v", raw, err)
			}
		})
	}

	invalid := []string{"", ".owner", "owner.", "owner..team", "Owner", "owner_team", "1owner", strings.Repeat("x", maxMetadataKeyLength+1)}
	for _, raw := range invalid {
		raw := raw

		t.Run("invalid "+raw, func(t *testing.T) {
			t.Parallel()

			err := validateMetadataKey("field", raw)
			if err == nil {
				t.Fatalf("validateMetadataKey(%q) returned nil error", raw)
			}

			if !errors.Is(err, ErrInvalidMetadata) {
				t.Fatalf("validateMetadataKey(%q) error = %v, want ErrInvalidMetadata", raw, err)
			}
		})
	}
}

func TestValidateMetadataKeyWrapsTextvalidateErrors(t *testing.T) {
	t.Parallel()

	err := validateMetadataKey("field", "Owner")
	if err == nil {
		t.Fatalf("validateMetadataKey() returned nil error")
	}

	if !errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("validateMetadataKey() error = %v, want ErrInvalidMetadata", err)
	}

	if !errors.Is(err, textvalidate.ErrInvalidDottedKebabKey) {
		t.Fatalf("validateMetadataKey() error = %v, want ErrInvalidDottedKebabKey", err)
	}
}

func TestValidateMetadataText(t *testing.T) {
	t.Parallel()

	if err := validateMetadataText("field", "hello\nworld", maxMetadataTextLength); err != nil {
		t.Fatalf("validateMetadataText() returned unexpected error: %v", err)
	}

	tests := []string{
		string([]byte{0xff, 0xfe}),
		strings.Repeat("x", maxMetadataTextLength+1),
		"bad\x00value",
	}

	for _, raw := range tests {
		raw := raw

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			err := validateMetadataText("field", raw, maxMetadataTextLength)
			if err == nil {
				t.Fatalf("validateMetadataText() returned nil error")
			}

			if !errors.Is(err, ErrInvalidMetadata) {
				t.Fatalf("validateMetadataText() error = %v, want ErrInvalidMetadata", err)
			}

			if !errors.Is(err, textvalidate.ErrInvalidCompactText) {
				t.Fatalf("validateMetadataText() error = %v, want ErrInvalidCompactText", err)
			}
		})
	}
}

func TestCloneStringMap(t *testing.T) {
	t.Parallel()

	input := map[string]string{"a": "1"}
	output := cloneStringMap(input)
	output["a"] = "2"

	if got, want := input["a"], "1"; got != want {
		t.Fatalf("input mutated through clone: got %q, want %q", got, want)
	}

	if cloneStringMap(nil) != nil {
		t.Fatalf("cloneStringMap(nil) must return nil")
	}
}
