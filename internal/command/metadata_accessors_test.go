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

func TestMetadataAccessors(t *testing.T) {
	t.Parallel()

	metadata := mustTestMetadata(t)

	if !metadata.HasOwner() {
		t.Fatalf("HasOwner() = false, want true")
	}

	if !metadata.HasArea() {
		t.Fatalf("HasArea() = false, want true")
	}

	if !metadata.HasSince() {
		t.Fatalf("HasSince() = false, want true")
	}

	if !metadata.IsDeprecated() {
		t.Fatalf("IsDeprecated() = false, want true")
	}

	deprecation, ok := metadata.Deprecation()
	if !ok {
		t.Fatalf("Deprecation() ok = false, want true")
	}

	if got, want := deprecation.Message(), "Use bench run instead."; got != want {
		t.Fatalf("Deprecation().Message() = %q, want %q", got, want)
	}
}

func TestMetadataSpecReturnsDetachedCopy(t *testing.T) {
	t.Parallel()

	metadata := mustTestMetadata(t)

	spec := metadata.Spec()
	spec.Owner = "changed"
	spec.Annotations["docs.page"] = "changed"
	spec.Deprecation.Message = "changed"

	if got, want := metadata.Owner(), "devtools"; got != want {
		t.Fatalf("metadata mutated through Spec Owner: got %q, want %q", got, want)
	}

	if got, want := mustTestAnnotation(t, metadata, "docs.page"), "commands/bench.md"; got != want {
		t.Fatalf("metadata mutated through Spec Annotations: got %q, want %q", got, want)
	}

	deprecation, ok := metadata.Deprecation()
	if !ok {
		t.Fatalf("Deprecation() ok = false, want true")
	}

	if got, want := deprecation.Message(), "Use bench run instead."; got != want {
		t.Fatalf("metadata mutated through Spec Deprecation: got %q, want %q", got, want)
	}
}

func TestMetadataDeprecationReturnsValueCopy(t *testing.T) {
	t.Parallel()

	metadata := mustTestMetadata(t)
	deprecation, ok := metadata.Deprecation()
	if !ok {
		t.Fatalf("Deprecation() ok = false, want true")
	}

	deprecation = deprecation.MustWithMessage("changed")

	again, ok := metadata.Deprecation()
	if !ok {
		t.Fatalf("Deprecation() ok = false, want true")
	}

	if got, want := again.Message(), "Use bench run instead."; got != want {
		t.Fatalf("metadata changed through deprecation value copy: got %q, want %q", got, want)
	}
}
