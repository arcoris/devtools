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

// Spec returns a detached construction spec.
func (reference DocumentationReference) Spec() DocumentationReferenceSpec {
	return DocumentationReferenceSpec{
		Key:    reference.key,
		Kind:   reference.kind,
		Label:  reference.label,
		Target: reference.target,
	}
}

// Key returns the stable machine-facing reference key.
func (reference DocumentationReference) Key() string {
	return reference.key
}

// Kind returns the reference kind.
func (reference DocumentationReference) Kind() DocumentationReferenceKind {
	return reference.kind
}

// Label returns the human-facing reference label.
func (reference DocumentationReference) Label() string {
	return reference.label
}

// Target returns the machine-facing reference target.
func (reference DocumentationReference) Target() string {
	return reference.target
}

// String returns the reference target.
func (reference DocumentationReference) String() string {
	return reference.target
}

// IsZero reports whether reference has not been set.
func (reference DocumentationReference) IsZero() bool {
	return reference.key == "" &&
		reference.kind == "" &&
		reference.label == "" &&
		reference.target == ""
}

// IsValid reports whether reference satisfies structural rules.
func (reference DocumentationReference) IsValid() bool {
	return reference.Validate() == nil
}
