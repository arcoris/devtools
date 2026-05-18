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

	"arcoris.dev/devtools/internal/textvalidate"
)

// TestEventValidationHelpers verifies lower-level validators.
func TestEventValidationHelpers(t *testing.T) {
	t.Parallel()

	if err := validateEventKey("field", "command.started", ErrInvalidEvent); err != nil {
		t.Fatalf("validateEventKey(valid) returned unexpected error: %v", err)
	}

	if err := validateEventKey("field", "Command.Started", ErrInvalidEvent); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("validateEventKey(invalid) error = %v, want ErrInvalidEvent", err)
	}

	if err := validateEventBlock("message", "Line one.\nLine two.", maxEventMessageLength); err != nil {
		t.Fatalf("validateEventBlock(valid) returned unexpected error: %v", err)
	}

	if err := validateEventBlock("message", "bad\x00message", maxEventMessageLength); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("validateEventBlock(control) error = %v, want ErrInvalidEvent", err)
	}

	if err := validateEventBlock("message", "bad\x00message", maxEventMessageLength); !errors.Is(err, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("validateEventBlock(control) error = %v, want ErrInvalidCompactText", err)
	}

	invalidUTF8 := string([]byte{0xff, 0xfe})
	if err := validateEventBlock("message", invalidUTF8, maxEventMessageLength); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("validateEventBlock(invalid UTF-8) error = %v, want ErrInvalidEvent", err)
	}
}
