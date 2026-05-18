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

import "testing"

func TestDeprecationAccessorsAndSpec(t *testing.T) {
	t.Parallel()

	deprecation := MustDeprecation(DeprecationSpec{
		Since:       "v0.2.0",
		Message:     "Use bench run instead.",
		Replacement: MustPath("bench", "run"),
	})

	if deprecation.IsZero() {
		t.Fatalf("IsZero() = true, want false")
	}

	if !deprecation.HasSince() {
		t.Fatalf("HasSince() = false, want true")
	}

	if !deprecation.HasReplacement() {
		t.Fatalf("HasReplacement() = false, want true")
	}

	spec := deprecation.Spec()
	spec.Message = "changed"

	if got, want := deprecation.Message(), "Use bench run instead."; got != want {
		t.Fatalf("deprecation mutated through Spec: got %q, want %q", got, want)
	}
}

func TestDeprecationZeroValue(t *testing.T) {
	t.Parallel()

	var deprecation Deprecation
	if !deprecation.IsZero() {
		t.Fatalf("zero deprecation IsZero() = false, want true")
	}
}
