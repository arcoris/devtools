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
	// ErrInvalidBinding reports that a command input binding declaration is
	// malformed.
	ErrInvalidBinding = errors.New("command binding is invalid")

	// ErrInvalidBindingValue reports that runtime values do not satisfy a
	// command input binding declaration.
	ErrInvalidBindingValue = errors.New("command binding value is invalid")
)

// BindingSpec describes command input binding before validation.
//
// BindingSpec is a construction DTO. Binding stores detached copies of mutable
// input state, so callers cannot mutate constructed bindings through shared
// slices.
//
// Binding is declaration-level input contract. It does not parse os.Args,
// inspect Cobra flags, read environment variables, load configuration, or
// perform interactive prompts. Adapter and resolver layers must do that work
// and then pass already resolved option values and positional values here.
type BindingSpec struct {
	// Options contains command option declarations.
	//
	// Option canonical names, aliases, and shorthands must be unique inside one
	// binding.
	Options []Option

	// Arguments contains positional argument declarations in consumption order.
	//
	// At most one variadic argument is allowed, and it must be the last
	// declaration. A required argument must not follow an optional argument
	// because that would make positional consumption ambiguous.
	Arguments []Argument
}

// BindingValueSpec describes runtime values to bind against a Binding.
type BindingValueSpec struct {
	// OptionValues contains already resolved option values.
	//
	// Values may use canonical option names or alias names. Binding canonicalizes
	// them to declared option names.
	OptionValues []OptionValue

	// PositionalValues contains runtime positional values in invocation order.
	PositionalValues []string
}

// Binding is a validated adapter-neutral command input contract.
//
// Binding connects declaration-level inputs:
//
//   - Options: named non-positional values;
//   - Arguments: ordered positional values.
//
// Binding deliberately does not own parsing, shell syntax, config lookup,
// environment lookup, or prompting. It validates and canonicalizes values after
// those adapter-specific stages have already happened.
type Binding struct {
	options   []Option
	arguments []Argument
}
