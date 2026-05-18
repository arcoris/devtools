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

// NewDocumentationReference validates spec and returns DocumentationReference.
func NewDocumentationReference(spec DocumentationReferenceSpec) (DocumentationReference, error) {
	if err := validateDocumentationReferenceSpecInput(spec); err != nil {
		return DocumentationReference{}, err
	}

	reference := DocumentationReference{
		key:    normalizeDocumentationReferenceKey(spec.Key),
		kind:   spec.Kind,
		label:  normalizeDocumentationSingleLine(spec.Label),
		target: normalizeDocumentationSingleLine(spec.Target),
	}

	if err := reference.Validate(); err != nil {
		return DocumentationReference{}, err
	}

	return reference, nil
}

// MustDocumentationReference validates spec and returns DocumentationReference.
//
// MustDocumentationReference panics on invalid input. It is intended for static
// command definitions and tests.
func MustDocumentationReference(spec DocumentationReferenceSpec) DocumentationReference {
	reference, err := NewDocumentationReference(spec)
	if err != nil {
		panic(err)
	}

	return reference
}
