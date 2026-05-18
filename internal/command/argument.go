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

const (
	// maxArgumentNameLength is the maximum allowed byte length of one argument
	// name.
	maxArgumentNameLength = 63

	// maxArgumentMetavarLength is the maximum allowed byte length of an
	// argument value metavar.
	maxArgumentMetavarLength = 63

	// maxArgumentValueLength is the maximum allowed byte length of one raw
	// argument value.
	maxArgumentValueLength = 4096
)

var (
	// ErrEmptyArgumentName reports that an argument name was not provided.
	ErrEmptyArgumentName = errors.New("command argument name is empty")

	// ErrInvalidArgumentName reports that an argument name violates the
	// argument-name grammar.
	ErrInvalidArgumentName = errors.New("command argument name is invalid")

	// ErrInvalidArgument reports that an argument declaration is malformed.
	ErrInvalidArgument = errors.New("command argument is invalid")

	// ErrInvalidArgumentValue reports that runtime argument values do not match
	// an argument declaration.
	ErrInvalidArgumentValue = errors.New("command argument value is invalid")
)

// ArgumentRequirement describes whether a positional argument must be supplied.
type ArgumentRequirement string

const (
	// ArgumentRequirementRequired means at least one value must be supplied.
	ArgumentRequirementRequired ArgumentRequirement = "required"

	// ArgumentRequirementOptional means the argument may be omitted.
	ArgumentRequirementOptional ArgumentRequirement = "optional"
)

// ArgumentCardinality describes how many values one positional argument
// declaration may consume.
type ArgumentCardinality string

const (
	// ArgumentCardinalitySingle means the declaration consumes at most one
	// positional value.
	ArgumentCardinalitySingle ArgumentCardinality = "single"

	// ArgumentCardinalityVariadic means the declaration consumes zero or more
	// trailing positional values, depending on Requirement.
	//
	// A command argument list should normally allow at most one variadic
	// argument, and that argument should be the last declaration. That global
	// rule belongs to an argument-list or command-definition validator, not to
	// a single Argument value object.
	ArgumentCardinalityVariadic ArgumentCardinality = "variadic"
)

// ArgumentSpec describes a positional command argument before validation.
//
// ArgumentSpec is a construction DTO. Argument stores detached copies of
// mutable input state, so callers cannot mutate constructed declarations
// through shared slices.
//
// Argument describes declaration-time contract only. It is not a runtime
// invocation value. Runtime positional values live in Invocation and
// ActionRequest; this type describes how those values should be interpreted.
type ArgumentSpec struct {
	// Name is the canonical machine-facing argument name.
	//
	// Name should be stable because it may be used in documentation anchors,
	// generated schemas, diagnostics, reports, and tests.
	Name string

	// Kind describes the scalar value shape accepted by this argument.
	//
	// Positional arguments use scalar OptionKind values. Multiplicity is modeled
	// by Cardinality, not by list-shaped OptionKind values.
	Kind OptionKind

	// Requirement controls whether the argument must be present.
	//
	// Zero defaults to required.
	Requirement ArgumentRequirement

	// Cardinality controls whether the argument consumes one value or multiple
	// trailing values.
	//
	// Zero defaults to single.
	Cardinality ArgumentCardinality

	// EmptyValue controls whether an explicitly supplied empty value is allowed.
	//
	// Zero defaults to reject-empty.
	EmptyValue OptionEmptyValuePolicy

	// Metavar is an optional value placeholder for help and documentation.
	//
	// If empty, a metavar is derived from Name.
	Metavar string

	// DefaultValues contains declaration defaults.
	//
	// Defaults are only valid for optional arguments. Single arguments may have
	// at most one default value. Variadic arguments may have multiple defaults.
	DefaultValues []string

	// AllowedValues contains allowed values for enum-like or constrained string
	// arguments.
	//
	// Enum arguments require allowed values. String arguments may use them for
	// lightweight value constraints. Numeric and duration arguments should use
	// typed range validation in a higher layer.
	AllowedValues []string

	// Documentation contains human-facing argument documentation.
	Documentation Documentation

	// Metadata contains machine-facing argument metadata.
	Metadata Metadata

	// Visibility controls default exposure in help, docs, and discovery.
	//
	// A zero visibility defaults to public.
	Visibility Visibility
}

// Argument is a validated framework-neutral positional command argument
// declaration.
//
// Argument is intentionally separate from Usage and Invocation:
//
//   - Usage describes rendered syntax;
//   - Invocation stores runtime values;
//   - Argument defines the contract for one named positional value.
//
// Argument is immutable-style:
//
//   - constructors normalize defaults and copy slices;
//   - accessors return detached copies;
//   - With* methods return validated copies;
//   - callers cannot mutate internal state through returned values.
type Argument struct {
	name          ArgumentName
	kind          OptionKind
	requirement   ArgumentRequirement
	cardinality   ArgumentCardinality
	emptyValue    OptionEmptyValuePolicy
	metavar       string
	defaultValues []string
	allowedValues []string
	documentation Documentation
	metadata      Metadata
	visibility    Visibility
}
