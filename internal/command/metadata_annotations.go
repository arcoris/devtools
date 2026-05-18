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

import "sort"

// Annotation returns one annotation value and whether it exists.
func (metadata Metadata) Annotation(key string) (string, bool) {
	if metadata.annotations == nil {
		return "", false
	}

	value, ok := metadata.annotations[key]

	return value, ok
}

// HasAnnotation reports whether an annotation key exists.
func (metadata Metadata) HasAnnotation(key string) bool {
	_, ok := metadata.Annotation(key)

	return ok
}

// HasAnnotations reports whether any annotations are set.
func (metadata Metadata) HasAnnotations() bool {
	return len(metadata.annotations) > 0
}

// AnnotationCount returns the number of annotations.
func (metadata Metadata) AnnotationCount() int {
	return len(metadata.annotations)
}

// Annotations returns a detached copy of all annotations.
func (metadata Metadata) Annotations() map[string]string {
	return cloneStringMap(metadata.annotations)
}

// AnnotationKeys returns annotation keys in deterministic lexical order.
func (metadata Metadata) AnnotationKeys() []string {
	keys := make([]string, 0, len(metadata.annotations))
	for key := range metadata.annotations {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// WithAnnotations returns a validated copy with annotations replaced.
func (metadata Metadata) WithAnnotations(annotations map[string]string) (Metadata, error) {
	spec := metadata.spec()
	spec.Annotations = cloneStringMap(annotations)

	return NewMetadata(spec)
}

// MustWithAnnotations returns a validated copy with annotations replaced and
// panics on invalid input.
func (metadata Metadata) MustWithAnnotations(annotations map[string]string) Metadata {
	next, err := metadata.WithAnnotations(annotations)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutAnnotations returns a validated copy without annotations.
func (metadata Metadata) WithoutAnnotations() Metadata {
	spec := metadata.spec()
	spec.Annotations = nil

	return MustMetadata(spec)
}

// WithAnnotation returns a validated copy with one annotation added or replaced.
func (metadata Metadata) WithAnnotation(key string, value string) (Metadata, error) {
	spec := metadata.spec()
	if spec.Annotations == nil {
		spec.Annotations = make(map[string]string)
	}

	spec.Annotations[key] = value

	return NewMetadata(spec)
}

// MustWithAnnotation returns a validated copy with one annotation added or
// replaced and panics on invalid input.
func (metadata Metadata) MustWithAnnotation(key string, value string) Metadata {
	next, err := metadata.WithAnnotation(key, value)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutAnnotation returns a validated copy without one annotation key.
func (metadata Metadata) WithoutAnnotation(key string) Metadata {
	spec := metadata.spec()
	delete(spec.Annotations, key)

	return MustMetadata(spec)
}
