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

// Spec returns a detached construction spec for deprecation.
func (deprecation Deprecation) Spec() DeprecationSpec {
	return deprecation.spec()
}

// IsZero reports whether deprecation has no fields set.
func (deprecation Deprecation) IsZero() bool {
	return deprecation.since == "" &&
		deprecation.message == "" &&
		deprecation.replacement.IsRoot()
}

// Since returns the compact deprecation marker.
func (deprecation Deprecation) Since() string {
	return deprecation.since
}

// HasSince reports whether Since is set.
func (deprecation Deprecation) HasSince() bool {
	return deprecation.since != ""
}

// Message returns the human-facing deprecation message.
func (deprecation Deprecation) Message() string {
	return deprecation.message
}

// Replacement returns the replacement command path and whether it is set.
func (deprecation Deprecation) Replacement() (Path, bool) {
	if deprecation.replacement.IsRoot() {
		return RootPath(), false
	}

	return deprecation.replacement, true
}

// HasReplacement reports whether a replacement command path is set.
func (deprecation Deprecation) HasReplacement() bool {
	_, ok := deprecation.Replacement()

	return ok
}

// spec returns a detached construction spec.
func (deprecation Deprecation) spec() DeprecationSpec {
	return DeprecationSpec{
		Since:       deprecation.since,
		Message:     deprecation.message,
		Replacement: deprecation.replacement,
	}
}
