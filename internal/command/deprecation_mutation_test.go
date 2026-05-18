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

func TestDeprecationWithMethods(t *testing.T) {
	t.Parallel()

	deprecation := MustDeprecation(DeprecationSpec{Message: "Deprecated."}).
		MustWithSince("v0.2.0").
		MustWithMessage("Use bench run instead.").
		MustWithReplacement(MustPath("bench", "run"))

	if got, want := deprecation.Since(), "v0.2.0"; got != want {
		t.Fatalf("Since() = %q, want %q", got, want)
	}

	if got, want := deprecation.Message(), "Use bench run instead."; got != want {
		t.Fatalf("Message() = %q, want %q", got, want)
	}

	if !deprecation.HasReplacement() {
		t.Fatalf("HasReplacement() = false, want true")
	}

	if deprecation.WithoutReplacement().HasReplacement() {
		t.Fatalf("WithoutReplacement() still has replacement")
	}
}

func TestDeprecationWithMethodsRejectInvalidValues(t *testing.T) {
	t.Parallel()

	deprecation := MustDeprecation(DeprecationSpec{Message: "Deprecated."})

	if _, err := deprecation.WithSince("bad\x00value"); !errors.Is(err, ErrInvalidDeprecation) {
		t.Fatalf("WithSince() error = %v, want ErrInvalidDeprecation", err)
	}

	if _, err := deprecation.WithMessage("   "); !errors.Is(err, ErrInvalidDeprecation) {
		t.Fatalf("WithMessage() error = %v, want ErrInvalidDeprecation", err)
	}

	if _, err := deprecation.WithReplacement(Path{segments: []string{"Invalid"}}); !errors.Is(err, ErrInvalidDeprecation) {
		t.Fatalf("WithReplacement() error = %v, want ErrInvalidDeprecation", err)
	}
}
