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

func TestMetadataWithMethods(t *testing.T) {
	t.Parallel()

	metadata := EmptyMetadata().
		MustWithOwner("devtools").
		MustWithArea("command").
		MustWithSince("v0.1.0").
		MustWithAnnotation("docs.page", "commands.md").
		MustWithDeprecation(DeprecationSpec{
			Message:     "Use check instead.",
			Replacement: MustPath("check"),
		})

	if got, want := metadata.Owner(), "devtools"; got != want {
		t.Fatalf("Owner() = %q, want %q", got, want)
	}

	if got, want := metadata.Area(), "command"; got != want {
		t.Fatalf("Area() = %q, want %q", got, want)
	}

	if got, want := metadata.Since(), "v0.1.0"; got != want {
		t.Fatalf("Since() = %q, want %q", got, want)
	}

	if !metadata.HasDeprecation() {
		t.Fatalf("HasDeprecation() = false, want true")
	}

	if metadata.WithoutDeprecation().HasDeprecation() {
		t.Fatalf("WithoutDeprecation() still has deprecation")
	}
}

func TestMetadataWithMethodsRejectInvalidValues(t *testing.T) {
	t.Parallel()

	metadata := EmptyMetadata()

	if _, err := metadata.WithOwner("Invalid"); !errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("WithOwner() error = %v, want ErrInvalidMetadata", err)
	}

	if _, err := metadata.WithArea("area..bad"); !errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("WithArea() error = %v, want ErrInvalidMetadata", err)
	}

	if _, err := metadata.WithSince("bad\x00value"); !errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("WithSince() error = %v, want ErrInvalidMetadata", err)
	}

	if _, err := metadata.WithAnnotation("Invalid", "value"); !errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("WithAnnotation() error = %v, want ErrInvalidMetadata", err)
	}

	if _, err := metadata.WithAnnotations(map[string]string{"Invalid": "value"}); !errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("WithAnnotations() error = %v, want ErrInvalidMetadata", err)
	}

	if _, err := metadata.WithDeprecation(DeprecationSpec{}); !errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("WithDeprecation() error = %v, want ErrInvalidMetadata", err)
	}
}
