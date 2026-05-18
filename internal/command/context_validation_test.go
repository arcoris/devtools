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

func TestValidateContextFieldKey(t *testing.T) {
	t.Parallel()

	valid := []string{"source", "ci.mode", "artifact-kind"}
	for _, key := range valid {
		key := key

		t.Run("valid "+key, func(t *testing.T) {
			t.Parallel()

			if err := validateContextFieldKey("field", key); err != nil {
				t.Fatalf("validateContextFieldKey(%q) returned unexpected error: %v", key, err)
			}
		})
	}
}

func TestValidateContextFieldKeyWrapsReusableValidatorError(t *testing.T) {
	t.Parallel()

	err := validateContextFieldKey("field", "Source")
	if err == nil {
		t.Fatalf("validateContextFieldKey() returned nil error")
	}

	if !errors.Is(err, ErrInvalidContext) {
		t.Fatalf("validateContextFieldKey() error = %v, want ErrInvalidContext", err)
	}

	if !errors.Is(err, textvalidate.ErrInvalidDottedKebabKey) {
		t.Fatalf("validateContextFieldKey() error = %v, want ErrInvalidDottedKebabKey", err)
	}
}

func TestValidateContextText(t *testing.T) {
	t.Parallel()

	if err := validateContextText("field", "hello\nworld", maxContextFieldValueLength); err != nil {
		t.Fatalf("validateContextText() returned unexpected error: %v", err)
	}

	tests := []string{
		string([]byte{0xff, 0xfe}),
		strings.Repeat("x", maxContextFieldValueLength+1),
		"bad\x00value",
	}

	for _, raw := range tests {
		raw := raw

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			err := validateContextText("field", raw, maxContextFieldValueLength)
			if err == nil {
				t.Fatalf("validateContextText() returned nil error")
			}

			if !errors.Is(err, ErrInvalidContext) {
				t.Fatalf("validateContextText() error = %v, want ErrInvalidContext", err)
			}

			if !errors.Is(err, textvalidate.ErrInvalidCompactText) {
				t.Fatalf("validateContextText() error = %v, want ErrInvalidCompactText", err)
			}
		})
	}
}

func TestContextValidateWrapsInvocationError(t *testing.T) {
	t.Parallel()

	_, err := NewContext(ContextSpec{
		Node: contextTestCommandNode(),
		Invocation: Invocation{
			arguments: []string{"bad\x00arg"},
		},
	})
	if err == nil {
		t.Fatalf("NewContext() returned nil error")
	}

	if !errors.Is(err, ErrInvalidContext) {
		t.Fatalf("NewContext() error = %v, want ErrInvalidContext", err)
	}

	if !errors.Is(err, ErrInvalidInvocation) {
		t.Fatalf("NewContext() error = %v, want ErrInvalidInvocation", err)
	}
}
