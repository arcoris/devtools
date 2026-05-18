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

func TestValidateCompactTextAcceptsValidText(t *testing.T) {
	t.Parallel()

	if err := ValidateCompactText("hello\nworld\tok\r", 32); err != nil {
		t.Fatalf("ValidateCompactText() returned unexpected error: %v", err)
	}

	if err := ValidateCompactText("", 0); err != nil {
		t.Fatalf("ValidateCompactText(empty) returned unexpected error: %v", err)
	}
}

func TestValidateCompactTextRejectsInvalidText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		text  string
		limit int
	}{
		{name: "invalid UTF-8", text: string([]byte{0xff, 0xfe}), limit: 32},
		{name: "too long", text: strings.Repeat("x", 33), limit: 32},
		{name: "control rune", text: "bad\x00value", limit: 32},
		{name: "negative limit", text: "ok", limit: -1},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateCompactText(test.text, test.limit)
			if err == nil {
				t.Fatalf("ValidateCompactText() returned nil error")
			}

			if !errors.Is(err, ErrInvalidCompactText) {
				t.Fatalf("ValidateCompactText() error = %v, want ErrInvalidCompactText", err)
			}
		})
	}
}

func TestIsDisallowedCompactTextControlRune(t *testing.T) {
	t.Parallel()

	tests := []struct {
		r    rune
		want bool
	}{
		{r: '\x00', want: true},
		{r: '\x1f', want: true},
		{r: '\t', want: false},
		{r: '\n', want: false},
		{r: '\r', want: false},
		{r: 'a', want: false},
	}

	for _, test := range tests {
		test := test

		t.Run(string(test.r), func(t *testing.T) {
			t.Parallel()

			if got := IsDisallowedCompactTextControlRune(test.r); got != test.want {
				t.Fatalf("IsDisallowedCompactTextControlRune(%q) = %v, want %v", test.r, got, test.want)
			}
		})
	}
}
