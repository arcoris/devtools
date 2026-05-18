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

// DeprecationSpec describes deprecation metadata before validation.
//
// DeprecationSpec is a construction DTO. Deprecation stores the validated
// value form, so callers can freely reuse or modify a spec after construction.
type DeprecationSpec struct {
	// Since is an optional compact text value describing when deprecation began.
	//
	// The command kernel does not impose SemVer or date semantics here. Higher
	// policy layers may enforce project-specific formats.
	Since string

	// Message is the required human-facing deprecation message.
	//
	// It should explain what is deprecated and what the user should do.
	Message string

	// Replacement is an optional command path that should be used instead.
	//
	// The root path means "no replacement".
	Replacement Path
}

// Deprecation describes deprecation metadata for a command-tree node.
//
// Deprecation is immutable-style. The replacement Path is already a value
// object, and all fields are validated on construction.
type Deprecation struct {
	// since stores the optional validated deprecation marker.
	since string

	// message stores the required validated deprecation guidance.
	message string

	// replacement stores the optional validated replacement command path.
	//
	// RootPath is the canonical "no replacement" value.
	replacement Path
}
