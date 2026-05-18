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
)

func TestNewMetadataAcceptsValidMetadata(t *testing.T) {
	t.Parallel()

	metadata := mustTestMetadata(t)

	if got, want := metadata.Owner(), "devtools"; got != want {
		t.Fatalf("Owner() = %q, want %q", got, want)
	}

	if got, want := metadata.Area(), "command.registry"; got != want {
		t.Fatalf("Area() = %q, want %q", got, want)
	}

	if got, want := metadata.Since(), "v0.1.0"; got != want {
		t.Fatalf("Since() = %q, want %q", got, want)
	}

	if !metadata.HasDeprecation() {
		t.Fatalf("HasDeprecation() = false, want true")
	}

	if got, want := mustTestAnnotation(t, metadata, "docs.page"), "commands/bench.md"; got != want {
		t.Fatalf("Annotation(docs.page) = %q, want %q", got, want)
	}
}

func TestEmptyMetadata(t *testing.T) {
	t.Parallel()

	metadata := EmptyMetadata()

	if !metadata.IsZero() {
		t.Fatalf("EmptyMetadata().IsZero() = false, want true")
	}

	if err := metadata.Validate(); err != nil {
		t.Fatalf("EmptyMetadata().Validate() returned unexpected error: %v", err)
	}
}

func TestNewMetadataRejectsInvalidMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec MetadataSpec
	}{
		{name: "invalid owner", spec: MetadataSpec{Owner: "DevTools"}},
		{name: "invalid area", spec: MetadataSpec{Area: "command..registry"}},
		{name: "too long since", spec: MetadataSpec{Since: strings.Repeat("x", maxMetadataTextLength+1)}},
		{name: "invalid annotation key", spec: MetadataSpec{Annotations: map[string]string{"Docs.Page": "value"}}},
		{name: "too long annotation value", spec: MetadataSpec{Annotations: map[string]string{"docs.page": strings.Repeat("x", maxAnnotationValueLength+1)}}},
		{name: "invalid deprecation", spec: MetadataSpec{Deprecation: &DeprecationSpec{}}},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewMetadata(test.spec)
			if err == nil {
				t.Fatalf("NewMetadata() returned nil error")
			}

			if !errors.Is(err, ErrInvalidMetadata) {
				t.Fatalf("NewMetadata() error = %v, want ErrInvalidMetadata", err)
			}
		})
	}
}

func TestNewMetadataWrapsInvalidDeprecation(t *testing.T) {
	t.Parallel()

	_, err := NewMetadata(MetadataSpec{
		Deprecation: &DeprecationSpec{},
	})
	if err == nil {
		t.Fatalf("NewMetadata() returned nil error")
	}

	if !errors.Is(err, ErrInvalidMetadata) {
		t.Fatalf("NewMetadata() error = %v, want ErrInvalidMetadata", err)
	}

	if !errors.Is(err, ErrInvalidDeprecation) {
		t.Fatalf("NewMetadata() error = %v, want ErrInvalidDeprecation", err)
	}
}

func TestMustMetadataPanicsForInvalidMetadata(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustMetadata(MetadataSpec{Owner: "Invalid"})
	})
}
