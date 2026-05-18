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

import "errors"

var (
	// ErrInvalidOptionPolicy reports that an option policy is malformed.
	ErrInvalidOptionPolicy = errors.New("command option policy is invalid")
)

// OptionPolicySpec describes option declaration policy before validation.
//
// OptionPolicySpec is a construction DTO. OptionPolicy stores detached copies
// of mutable state, especially AllowedSources.
type OptionPolicySpec struct {
	// Requirement controls whether a resolved value must exist.
	//
	// Zero defaults to optional.
	Requirement OptionRequirement

	// Scope controls where the option declaration applies in the command tree.
	//
	// Zero defaults to local.
	Scope OptionScope

	// Occurrence controls whether the option may appear once or multiple times.
	//
	// Zero defaults from OptionKind when ValidateForKind is used. Without a
	// concrete OptionKind, zero defaults to single.
	Occurrence OptionOccurrence

	// EmptyValue controls whether explicitly supplied empty values are accepted.
	//
	// Zero defaults to reject-empty.
	EmptyValue OptionEmptyValuePolicy

	// AllowedSources declares which sources may provide this option's value.
	//
	// Empty means all known option sources are allowed. The stored order is
	// normalized into default precedence order.
	AllowedSources []OptionSource
}

// OptionPolicy describes declaration-level behavior for one option.
//
// OptionPolicy is intentionally separate from OptionKind and OptionSource:
//
//   - OptionKind describes value shape;
//   - OptionSource describes where a resolved value came from;
//   - OptionPolicy describes how a concrete option declaration should behave.
//
// OptionPolicy is immutable-style:
//
//   - constructors normalize zero fields into explicit defaults;
//   - slices are copied;
//   - accessors return detached copies;
//   - With* methods return validated copies.
type OptionPolicy struct {
	// requirement stores whether a resolved value must be present.
	requirement OptionRequirement

	// scope stores where the option declaration applies in the command tree.
	scope OptionScope

	// occurrence stores whether the option may appear once or multiple times.
	occurrence OptionOccurrence

	// emptyValue stores whether explicitly supplied empty values are accepted.
	emptyValue OptionEmptyValuePolicy

	// allowedSources stores source classes accepted by this declaration.
	allowedSources []OptionSource
}
