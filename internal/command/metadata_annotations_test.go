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

func TestMetadataAnnotationCopySemantics(t *testing.T) {
	t.Parallel()

	input := map[string]string{"docs.page": "commands.md"}
	metadata := MustMetadata(MetadataSpec{Annotations: input})

	input["docs.page"] = "changed"

	if got, want := mustTestAnnotation(t, metadata, "docs.page"), "commands.md"; got != want {
		t.Fatalf("metadata changed through input map: got %q, want %q", got, want)
	}

	out := metadata.Annotations()
	out["docs.page"] = "changed"

	if got, want := mustTestAnnotation(t, metadata, "docs.page"), "commands.md"; got != want {
		t.Fatalf("metadata changed through output map: got %q, want %q", got, want)
	}
}

func TestMetadataAnnotationKeys(t *testing.T) {
	t.Parallel()

	metadata := MustMetadata(MetadataSpec{
		Annotations: map[string]string{
			"z.key": "z",
			"a.key": "a",
			"m.key": "m",
		},
	})

	assertStringSlicesEqual(t, metadata.AnnotationKeys(), []string{"a.key", "m.key", "z.key"})
}

func TestMetadataAnnotationHelpers(t *testing.T) {
	t.Parallel()

	metadata := EmptyMetadata().
		MustWithAnnotation("docs.page", "commands.md").
		MustWithAnnotation("ci.mode", "smoke")

	if !metadata.HasAnnotations() {
		t.Fatalf("HasAnnotations() = false, want true")
	}

	if got, want := metadata.AnnotationCount(), 2; got != want {
		t.Fatalf("AnnotationCount() = %d, want %d", got, want)
	}

	if !metadata.HasAnnotation("docs.page") {
		t.Fatalf("HasAnnotation(docs.page) = false, want true")
	}

	withoutOne := metadata.WithoutAnnotation("docs.page")
	if withoutOne.HasAnnotation("docs.page") {
		t.Fatalf("WithoutAnnotation() still has docs.page")
	}

	withoutAll := metadata.WithoutAnnotations()
	if withoutAll.HasAnnotations() {
		t.Fatalf("WithoutAnnotations() still has annotations")
	}
}

func TestMetadataWithAnnotationsCopiesInput(t *testing.T) {
	t.Parallel()

	input := map[string]string{"docs.page": "commands.md"}
	metadata := EmptyMetadata().MustWithAnnotations(input)

	input["docs.page"] = "changed"

	if got, want := mustTestAnnotation(t, metadata, "docs.page"), "commands.md"; got != want {
		t.Fatalf("metadata changed through WithAnnotations input: got %q, want %q", got, want)
	}
}
