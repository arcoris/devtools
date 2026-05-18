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
	// maxOptionNameLength is the maximum allowed byte length of one option name.
	maxOptionNameLength = 63

	// maxOptionShorthandLength is the exact byte length expected for a
	// shorthand option name.
	maxOptionShorthandLength = 1

	// maxOptionMetavarLength is the maximum allowed byte length of an option
	// value metavar.
	maxOptionMetavarLength = 63

	// maxOptionValueLength is the maximum allowed byte length of one raw option
	// value declaration.
	maxOptionValueLength = 4096
)

var (
	// ErrEmptyOptionName reports that an option name was not provided.
	ErrEmptyOptionName = errors.New("command option name is empty")

	// ErrInvalidOptionName reports that an option name violates the option-name
	// grammar.
	ErrInvalidOptionName = errors.New("command option name is invalid")

	// ErrInvalidOption reports that an option declaration is malformed.
	ErrInvalidOption = errors.New("command option is invalid")

	// ErrInvalidOptionValue reports that an option value declaration is
	// malformed.
	ErrInvalidOptionValue = errors.New("command option value is invalid")
)

// OptionSpec describes a command option declaration before validation.
//
// OptionSpec is a construction DTO. Option stores detached copies of mutable
// input state, so callers cannot mutate a constructed declaration through
// shared slices.
//
// OptionSpec intentionally describes declaration-time metadata only. It does
// not represent a parsed value, resolved value, or parser binding. Runtime
// values should be modeled separately by the option-resolution layer.
type OptionSpec struct {
	// Name is the canonical long option name without the leading "--".
	//
	// Name is stable and machine-facing. It uses the same compact segment
	// grammar as command path segments:
	//
	//   - "verbose"
	//   - "output"
	//   - "bench-time"
	//   - "package"
	Name string

	// Aliases are alternative long option names without the leading "--".
	//
	// Aliases are compatibility spellings, not stable identity. They are useful
	// for migrations and adapter-level convenience.
	Aliases []string

	// Shorthand is an optional one-character short option name without "-".
	//
	// Example: "v" for "-v".
	Shorthand string

	// Kind describes the value shape accepted by this option.
	Kind OptionKind

	// Metavar is an optional value placeholder used by help and documentation.
	//
	// If empty, Kind.ValueMetavar() is used.
	Metavar string

	// DefaultValues contains declaration default values.
	//
	// Empty means no declaration default. Scalar kinds may have at most one
	// default value. List kinds may have multiple default values.
	DefaultValues []string

	// AllowedValues contains allowed values for enum-like or constrained string
	// options.
	//
	// Enum kinds require this field. String kinds may use it for lightweight
	// value constraints. Other kinds should use typed range validation in a
	// higher layer.
	AllowedValues []string

	// Policy controls declaration behavior such as requirement, scope,
	// occurrence, empty-value handling, and allowed value sources.
	//
	// A zero policy defaults to DefaultOptionPolicy adjusted for Kind.
	Policy OptionPolicy

	// Documentation contains human-facing option documentation.
	Documentation Documentation

	// Metadata contains machine-facing option metadata.
	Metadata Metadata

	// Visibility controls default exposure in help, docs, and discovery.
	//
	// A zero visibility defaults to public.
	Visibility Visibility
}

// Option is a validated framework-neutral command option declaration.
//
// Option is not a parser-specific flag. It does not depend on Cobra, pflag,
// urfave/cli, environment loaders, config loaders, or shell completion
// frameworks. Adapter layers can translate Option into concrete parser flags.
//
// Option connects declaration-level concerns:
//
//   - Name and aliases: option identity and compatibility spellings;
//   - Kind: value shape;
//   - Policy: requirement, scope, occurrence, and source policy;
//   - DefaultValues and AllowedValues: declaration-level value constraints;
//   - Documentation, Metadata, Visibility: presentation and lifecycle metadata.
//
// Resolved values are intentionally not modeled here. A future resolution layer
// should combine Option, OptionSource, parser input, defaults, config,
// environment, and validation into a separate resolved-value type.
type Option struct {
	// name stores the canonical stable long option name.
	name OptionName

	// aliases stores alternative long option names in declaration order.
	aliases []OptionName

	// shorthand stores the optional one-character short option name.
	shorthand string

	// kind stores the declared value shape.
	kind OptionKind

	// metavar stores the canonical value placeholder for help and docs.
	metavar string

	// defaultValues stores declaration defaults in declaration order.
	defaultValues []string

	// allowedValues stores declaration allowed values in declaration order.
	allowedValues []string

	// policy stores declaration behavior for resolution and adapters.
	policy OptionPolicy

	// documentation stores human-facing option documentation.
	documentation Documentation

	// metadata stores machine-facing option metadata.
	metadata Metadata

	// visibility stores default help/docs/discovery exposure.
	visibility Visibility
}
