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
	"fmt"
	"strings"
	"testing"
)

func TestValidateSingleLineTextAcceptsValidText(t *testing.T) {
	t.Parallel()

	if err := ValidateSingleLineText("bench run [flags]", 64); err != nil {
		t.Fatalf("ValidateSingleLineText() returned unexpected error: %v", err)
	}

	if err := ValidateSingleLineText("", 0); err != nil {
		t.Fatalf("ValidateSingleLineText(empty) returned unexpected error: %v", err)
	}
}

func TestValidateSingleLineTextRejectsInvalidText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		text  string
		limit int
	}{
		{name: "invalid UTF-8", text: string([]byte{0xff, 0xfe}), limit: 64},
		{name: "too long", text: strings.Repeat("x", 65), limit: 64},
		{name: "tab", text: "bench\trun", limit: 64},
		{name: "newline", text: "bench\nrun", limit: 64},
		{name: "carriage return", text: "bench\rrun", limit: 64},
		{name: "control rune", text: "bench\x00run", limit: 64},
		{name: "negative limit", text: "bench", limit: -1},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateSingleLineText(test.text, test.limit)
			if err == nil {
				t.Fatalf("ValidateSingleLineText() returned nil error")
			}

			if !errors.Is(err, ErrInvalidSingleLineText) {
				t.Fatalf("ValidateSingleLineText() error = %v, want ErrInvalidSingleLineText", err)
			}
		})
	}
}

func TestValidateTokenTextAcceptsValidText(t *testing.T) {
	t.Parallel()

	tests := []string{
		"bench",
		"[flags]",
		"<path>",
		"--",
		"--suite",
	}

	for _, text := range tests {
		text := text

		t.Run(text, func(t *testing.T) {
			t.Parallel()

			if err := ValidateTokenText(text, 64); err != nil {
				t.Fatalf("ValidateTokenText(%q) returned unexpected error: %v", text, err)
			}
		})
	}
}

func TestValidateTokenTextRejectsInvalidText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		text  string
		limit int
		err   error
	}{
		{name: "empty", text: "", limit: 64, err: ErrEmptyTokenText},
		{name: "invalid UTF-8", text: string([]byte{0xff, 0xfe}), limit: 64, err: ErrInvalidTokenText},
		{name: "too long", text: strings.Repeat("x", 65), limit: 64, err: ErrInvalidTokenText},
		{name: "space", text: "bench run", limit: 64, err: ErrInvalidTokenText},
		{name: "tab", text: "bench\trun", limit: 64, err: ErrInvalidTokenText},
		{name: "newline", text: "bench\nrun", limit: 64, err: ErrInvalidTokenText},
		{name: "control rune", text: "bench\x00run", limit: 64, err: ErrInvalidTokenText},
		{name: "invalid limit", text: "bench", limit: 0, err: ErrInvalidTokenText},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateTokenText(test.text, test.limit)
			if err == nil {
				t.Fatalf("ValidateTokenText() returned nil error")
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("ValidateTokenText() error = %v, want %v", err, test.err)
			}
		})
	}
}

func TestIsDisallowedSingleLineTextRune(t *testing.T) {
	t.Parallel()

	tests := []struct {
		r    rune
		want bool
	}{
		{r: '\x00', want: true},
		{r: '\x1f', want: true},
		{r: '\t', want: true},
		{r: '\n', want: true},
		{r: '\r', want: true},
		{r: ' ', want: false},
		{r: 'a', want: false},
	}

	for _, test := range tests {
		test := test

		t.Run(fmt.Sprintf("U+%04X", test.r), func(t *testing.T) {
			t.Parallel()

			if got := IsDisallowedSingleLineTextRune(test.r); got != test.want {
				t.Fatalf("IsDisallowedSingleLineTextRune(%q) = %v, want %v", test.r, got, test.want)
			}
		})
	}
}
