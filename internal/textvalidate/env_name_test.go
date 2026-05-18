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

func TestValidateEnvNameAcceptsValidNames(t *testing.T) {
	t.Parallel()

	tests := []string{
		"GOOS",
		"GOARCH",
		"_PRIVATE",
		"X1",
		"A_B_C",
	}

	for _, name := range tests {
		name := name

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if err := ValidateEnvName(name, 255); err != nil {
				t.Fatalf("ValidateEnvName(%q) returned unexpected error: %v", name, err)
			}
		})
	}
}

func TestValidateEnvNameRejectsInvalidNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		err  error
	}{
		{name: "empty", raw: "", err: ErrEmptyEnvName},
		{name: "lowercase", raw: "goos", err: ErrInvalidEnvName},
		{name: "digit start", raw: "1GOOS", err: ErrInvalidEnvName},
		{name: "hyphen", raw: "GO-OS", err: ErrInvalidEnvName},
		{name: "dot", raw: "GO.OS", err: ErrInvalidEnvName},
		{name: "space", raw: "GO OS", err: ErrInvalidEnvName},
		{name: "non ASCII", raw: "ГООС", err: ErrInvalidEnvName},
		{name: "invalid UTF-8", raw: string([]byte{0xff, 0xfe}), err: ErrInvalidEnvName},
		{name: "too long", raw: strings.Repeat("X", 4), err: ErrInvalidEnvName},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateEnvName(test.raw, 3)
			if err == nil {
				t.Fatalf("ValidateEnvName(%q) returned nil error", test.raw)
			}

			if !errors.Is(err, test.err) {
				t.Fatalf("ValidateEnvName(%q) error = %v, want %v", test.raw, err, test.err)
			}
		})
	}
}

func TestValidateEnvNameRejectsInvalidLimit(t *testing.T) {
	t.Parallel()

	err := ValidateEnvName("GOOS", 0)
	if err == nil {
		t.Fatalf("ValidateEnvName() returned nil error")
	}

	if !errors.Is(err, ErrInvalidEnvName) {
		t.Fatalf("ValidateEnvName() error = %v, want ErrInvalidEnvName", err)
	}
}

func TestEnvNamePredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		ch           byte
		start        bool
		continuation bool
	}{
		{ch: 'A', start: true, continuation: true},
		{ch: 'Z', start: true, continuation: true},
		{ch: '_', start: true, continuation: true},
		{ch: '0', start: false, continuation: true},
		{ch: '9', start: false, continuation: true},
		{ch: 'a', start: false, continuation: false},
		{ch: '-', start: false, continuation: false},
		{ch: '.', start: false, continuation: false},
	}

	for _, test := range tests {
		test := test

		t.Run(fmt.Sprintf("%q", test.ch), func(t *testing.T) {
			t.Parallel()

			if got := IsEnvNameStart(test.ch); got != test.start {
				t.Fatalf("IsEnvNameStart(%q) = %v, want %v", test.ch, got, test.start)
			}

			if got := IsEnvNameContinuation(test.ch); got != test.continuation {
				t.Fatalf("IsEnvNameContinuation(%q) = %v, want %v", test.ch, got, test.continuation)
			}
		})
	}
}
