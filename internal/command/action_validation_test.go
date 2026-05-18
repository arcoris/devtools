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

func TestValidateActionFieldKeys(t *testing.T) {
	t.Parallel()

	valid := []string{
		"mode",
		"ci.mode",
		"artifact-kind",
		"a1",
	}

	for _, raw := range valid {
		raw := raw

		t.Run("valid "+raw, func(t *testing.T) {
			t.Parallel()

			if err := validateActionRequestFieldKey("field", raw); err != nil {
				t.Fatalf("validateActionRequestFieldKey(%q) returned unexpected error: %v", raw, err)
			}

			if err := validateActionResultFieldKey("field", raw); err != nil {
				t.Fatalf("validateActionResultFieldKey(%q) returned unexpected error: %v", raw, err)
			}
		})
	}
}

func TestValidateActionFieldKeyUsesDomainSentinel(t *testing.T) {
	t.Parallel()

	requestErr := validateActionRequestFieldKey("field", "Bad")
	if requestErr == nil {
		t.Fatalf("validateActionRequestFieldKey() returned nil error")
	}

	if !errors.Is(requestErr, ErrInvalidActionRequest) {
		t.Fatalf("request key error = %v, want ErrInvalidActionRequest", requestErr)
	}

	if !errors.Is(requestErr, textvalidate.ErrInvalidDottedKebabKey) {
		t.Fatalf("request key error = %v, want ErrInvalidDottedKebabKey", requestErr)
	}

	if errors.Is(requestErr, ErrInvalidActionResult) {
		t.Fatalf("request key error = %v, must not wrap ErrInvalidActionResult", requestErr)
	}

	resultErr := validateActionResultFieldKey("field", "Bad")
	if resultErr == nil {
		t.Fatalf("validateActionResultFieldKey() returned nil error")
	}

	if !errors.Is(resultErr, ErrInvalidActionResult) {
		t.Fatalf("result key error = %v, want ErrInvalidActionResult", resultErr)
	}
}

func TestValidateActionTextUsesDomainSentinel(t *testing.T) {
	t.Parallel()

	if err := validateActionRequestText("field", "hello\nworld", maxActionMessageLength); err != nil {
		t.Fatalf("validateActionRequestText() returned unexpected error: %v", err)
	}

	requestErr := validateActionRequestText("field", "bad\x00value", maxActionMessageLength)
	if requestErr == nil {
		t.Fatalf("validateActionRequestText() returned nil error")
	}

	if !errors.Is(requestErr, ErrInvalidActionRequest) {
		t.Fatalf("request text error = %v, want ErrInvalidActionRequest", requestErr)
	}

	if !errors.Is(requestErr, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("request text error = %v, want ErrInvalidCompactText", requestErr)
	}

	if errors.Is(requestErr, ErrInvalidActionResult) {
		t.Fatalf("request text error = %v, must not wrap ErrInvalidActionResult", requestErr)
	}

	resultErr := validateActionResultText("field", strings.Repeat("x", maxActionMessageLength+1), maxActionMessageLength)
	if resultErr == nil {
		t.Fatalf("validateActionResultText() returned nil error")
	}

	if !errors.Is(resultErr, ErrInvalidActionResult) {
		t.Fatalf("result text error = %v, want ErrInvalidActionResult", resultErr)
	}
}

func TestValidateActionFields(t *testing.T) {
	t.Parallel()

	if err := validateActionRequestFields(map[string]string{"ci.mode": "smoke"}); err != nil {
		t.Fatalf("validateActionRequestFields() returned unexpected error: %v", err)
	}

	if err := validateActionResultFields(map[string]string{"ci.mode": "smoke"}); err != nil {
		t.Fatalf("validateActionResultFields() returned unexpected error: %v", err)
	}

	if err := validateActionRequestFields(map[string]string{"Bad": "smoke"}); !errors.Is(err, ErrInvalidActionRequest) {
		t.Fatalf("validateActionRequestFields(key) error = %v, want ErrInvalidActionRequest", err)
	}

	if err := validateActionResultFields(map[string]string{"ci.mode": "bad\x00value"}); !errors.Is(err, ErrInvalidActionResult) {
		t.Fatalf("validateActionResultFields(value) error = %v, want ErrInvalidActionResult", err)
	}
}
