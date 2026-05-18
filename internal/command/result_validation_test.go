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

// TestResultValidationHelpers verifies lower-level validators.
func TestResultValidationHelpers(t *testing.T) {
	t.Parallel()

	if err := validateResultKey("field", "ci.mode"); err != nil {
		t.Fatalf("validateResultKey(valid) returned unexpected error: %v", err)
	}

	if err := validateResultKey("field", "Ci.Mode"); !errors.Is(err, ErrInvalidResult) {
		t.Fatalf("validateResultKey(invalid) error = %v, want ErrInvalidResult", err)
	}

	if err := validateResultBlock("message", "Line one.\nLine two.", maxResultMessageLength); err != nil {
		t.Fatalf("validateResultBlock(valid) returned unexpected error: %v", err)
	}

	if err := validateResultBlock("message", "bad\x00message", maxResultMessageLength); !errors.Is(err, ErrInvalidResult) {
		t.Fatalf("validateResultBlock(control) error = %v, want ErrInvalidResult", err)
	}

	if err := validateResultBlock("message", "bad\x00message", maxResultMessageLength); !errors.Is(err, textvalidate.ErrInvalidCompactText) {
		t.Fatalf("validateResultBlock(control) error = %v, want ErrInvalidCompactText", err)
	}

	invalidUTF8 := string([]byte{0xff, 0xfe})
	if err := validateResultBlock("message", invalidUTF8, maxResultMessageLength); !errors.Is(err, ErrInvalidResult) {
		t.Fatalf("validateResultBlock(invalid UTF-8) error = %v, want ErrInvalidResult", err)
	}
}
