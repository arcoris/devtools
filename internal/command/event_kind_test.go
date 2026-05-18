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
)

// TestEventKindValidation verifies EventKind behavior.
func TestEventKindValidation(t *testing.T) {
	t.Parallel()

	kind, err := NewEventKind("custom.lifecycle")
	if err != nil {
		t.Fatalf("NewEventKind() returned unexpected error: %v", err)
	}

	if got, want := kind.String(), "custom.lifecycle"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	if _, err := NewEventKind(""); !errors.Is(err, ErrEmptyEventKind) {
		t.Fatalf("NewEventKind(empty) error = %v, want ErrEmptyEventKind", err)
	}

	if _, err := NewEventKind("Custom"); !errors.Is(err, ErrInvalidEventKind) {
		t.Fatalf("NewEventKind(invalid) error = %v, want ErrInvalidEventKind", err)
	}
}
