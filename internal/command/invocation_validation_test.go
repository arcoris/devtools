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

func TestValidateInvocationEnvName(t *testing.T) {
	t.Parallel()

	valid := []string{"GOOS", "GOARCH", "_PRIVATE", "X1", "A_B_C"}
	for _, name := range valid {
		name := name

		t.Run("valid "+name, func(t *testing.T) {
			t.Parallel()

			if err := validateInvocationEnvName(name); err != nil {
				t.Fatalf("validateInvocationEnvName(%q) returned unexpected error: %v", name, err)
			}
		})
	}
}

func TestValidateInvocationEnvNameWrapsReusableValidatorError(t *testing.T) {
	t.Parallel()

	err := validateInvocationEnvName("goos")
	if err == nil {
		t.Fatalf("validateInvocationEnvName() returned nil error")
	}

	if !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("validateInvocationEnvName() error = %v, want ErrInvalidInvocation", err)
	}

	if !errors.Is(err, textvalidate.ErrInvalidEnvName) {
		t.Fatalf("validateInvocationEnvName() error = %v, want ErrInvalidEnvName", err)
	}
}

func TestValidateInvocationFieldKey(t *testing.T) {
	t.Parallel()

	valid := []string{"source", "ci.mode", "artifact-kind"}
	for _, key := range valid {
		key := key

		t.Run("valid "+key, func(t *testing.T) {
			t.Parallel()

			if err := validateInvocationFieldKey(key); err != nil {
				t.Fatalf("validateInvocationFieldKey(%q) returned unexpected error: %v", key, err)
			}
		})
	}
}

func TestValidateInvocationFieldKeyWrapsReusableValidatorError(t *testing.T) {
	t.Parallel()

	err := validateInvocationFieldKey("Source")
	if err == nil {
		t.Fatalf("validateInvocationFieldKey() returned nil error")
	}

	if !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("validateInvocationFieldKey() error = %v, want ErrInvalidInvocation", err)
	}

	if !errors.Is(err, textvalidate.ErrInvalidDottedKebabKey) {
		t.Fatalf("validateInvocationFieldKey() error = %v, want ErrInvalidDottedKebabKey", err)
	}
}

func TestValidateInvocationText(t *testing.T) {
	t.Parallel()

	if err := validateInvocationText("field", "hello\nworld", maxInvocationFieldValueLength); err != nil {
		t.Fatalf("validateInvocationText() returned unexpected error: %v", err)
	}

	tests := []string{
		string([]byte{0xff, 0xfe}),
		strings.Repeat("x", maxInvocationFieldValueLength+1),
		"bad\x00value",
	}

	for _, raw := range tests {
		raw := raw

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			err := validateInvocationText("field", raw, maxInvocationFieldValueLength)
			if err == nil {
				t.Fatalf("validateInvocationText() returned nil error")
			}

			if !errors.Is(err, ErrInvalidInvocation) {
				t.Fatalf("validateInvocationText() error = %v, want ErrInvalidInvocation", err)
			}

			if !errors.Is(err, textvalidate.ErrInvalidCompactText) {
				t.Fatalf("validateInvocationText() error = %v, want ErrInvalidCompactText", err)
			}
		})
	}
}

func TestValidateInvocationFieldsAndEnv(t *testing.T) {
	t.Parallel()

	if err := validateInvocationEnv(map[string]string{"GOOS": "linux"}); err != nil {
		t.Fatalf("validateInvocationEnv() returned unexpected error: %v", err)
	}

	if err := validateInvocationFields(map[string]string{"ci.mode": "smoke"}); err != nil {
		t.Fatalf("validateInvocationFields() returned unexpected error: %v", err)
	}

	if err := validateInvocationEnv(map[string]string{"goos": "linux"}); !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("validateInvocationEnv(name) error = %v, want ErrInvalidInvocation", err)
	}

	if err := validateInvocationFields(map[string]string{"Bad": "value"}); !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("validateInvocationFields(key) error = %v, want ErrInvalidInvocation", err)
	}
}
