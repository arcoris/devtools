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

// OptionValueSpec describes one resolved option value before validation.
//
// OptionValueSpec is a construction DTO. OptionValue stores detached copies of
// mutable input state so callers cannot mutate constructed values through
// shared slices.
//
// OptionValue is runtime/resolution data, not an option declaration:
//
//   - Option declares what may be accepted;
//   - OptionSource declares where a value came from;
//   - OptionPolicy declares source/requirement/occurrence behavior;
//   - OptionValue stores the already resolved raw value or values.
//
// OptionValue stores raw strings deliberately. Typed accessors parse those
// strings on demand and keep parser/runtime code framework-neutral.
type OptionValueSpec struct {
	// Name is the canonical option name without "--".
	Name string

	// Kind describes the value shape.
	Kind OptionKind

	// Source describes where the resolved value came from.
	Source OptionSource

	// Values contains one or more raw resolved values.
	//
	// Scalar kinds MUST contain exactly one value. List kinds MUST contain one
	// or more values. Absence should be represented by no OptionValue, not by an
	// OptionValue with an empty value slice.
	Values []string
}

// OptionValue is one validated resolved value for one command option.
//
// OptionValue does not know about a concrete parser. It can represent values
// resolved from command-line flags, config, environment variables, runtime
// injection, interactive prompts, inherited scopes, or declaration defaults.
//
// OptionValue is immutable-style:
//
//   - constructors copy input slices;
//   - accessors return detached copies;
//   - With* methods return validated copies;
//   - callers cannot mutate internal state through returned values.
type OptionValue struct {
	// name stores the canonical option name for this resolved value.
	name OptionName

	// kind stores the resolved value shape.
	kind OptionKind

	// source stores where the value came from.
	source OptionSource

	// values stores one or more raw resolved values.
	values []string
}
