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

// TestArtifactKindValidation verifies open validated kind behavior.
func TestArtifactKindValidation(t *testing.T) {
	t.Parallel()

	kind, err := NewArtifactKind("custom-report")
	if err != nil {
		t.Fatalf("NewArtifactKind() returned unexpected error: %v", err)
	}

	if got, want := kind.String(), "custom-report"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	if _, err := NewArtifactKind(""); !errors.Is(err, ErrEmptyArtifactKind) {
		t.Fatalf("NewArtifactKind(empty) error = %v, want ErrEmptyArtifactKind", err)
	}

	if _, err := NewArtifactKind("Custom"); !errors.Is(err, ErrInvalidArtifactKind) {
		t.Fatalf("NewArtifactKind(invalid) error = %v, want ErrInvalidArtifactKind", err)
	}
}
