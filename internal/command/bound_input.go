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

import "fmt"

// BoundInput is a validated command input after declaration binding.
//
// BoundInput is the output of Binding.Bind. It contains canonical option values
// and named positional argument values. It is still adapter-neutral and does not
// imply process execution.
type BoundInput struct {
	options   []OptionValue
	arguments []BoundArgument
}

// Options returns detached bound option values.
func (input BoundInput) Options() []OptionValue {
	return cloneBindingOptionValues(input.options)
}

// Arguments returns detached bound positional argument values.
func (input BoundInput) Arguments() []BoundArgument {
	return cloneBoundArguments(input.arguments)
}

// Option returns a bound option value by canonical option name.
func (input BoundInput) Option(name OptionName) (OptionValue, bool) {
	for _, value := range input.options {
		if value.Name() == name {
			return value, true
		}
	}

	return OptionValue{}, false
}

// OptionByName parses raw and returns a bound option value by canonical option
// name.
func (input BoundInput) OptionByName(raw string) (OptionValue, bool) {
	name, err := NewOptionName(raw)
	if err != nil {
		return OptionValue{}, false
	}

	return input.Option(name)
}

// HasOption reports whether a bound option value exists.
func (input BoundInput) HasOption(name OptionName) bool {
	_, ok := input.Option(name)

	return ok
}

// Argument returns a bound positional argument by name.
func (input BoundInput) Argument(name ArgumentName) (BoundArgument, bool) {
	for _, argument := range input.arguments {
		if argument.Name() == name {
			return argument, true
		}
	}

	return BoundArgument{}, false
}

// ArgumentByName parses raw and returns a bound positional argument by name.
func (input BoundInput) ArgumentByName(raw string) (BoundArgument, bool) {
	name, err := NewArgumentName(raw)
	if err != nil {
		return BoundArgument{}, false
	}

	return input.Argument(name)
}

// HasArgument reports whether a bound positional argument exists.
func (input BoundInput) HasArgument(name ArgumentName) bool {
	_, ok := input.Argument(name)

	return ok
}

// OptionCount returns the number of bound option values.
func (input BoundInput) OptionCount() int {
	return len(input.options)
}

// ArgumentCount returns the number of bound argument values.
func (input BoundInput) ArgumentCount() int {
	return len(input.arguments)
}

// Validate verifies bound input structural rules.
func (input BoundInput) Validate() error {
	seenOptions := make(map[OptionName]struct{}, len(input.options))
	for index, value := range input.options {
		if err := value.Validate(); err != nil {
			return fmt.Errorf("%w: option %d: %w", ErrInvalidBindingValue, index, err)
		}

		if _, exists := seenOptions[value.Name()]; exists {
			return fmt.Errorf("%w: duplicate bound option %q", ErrInvalidBindingValue, value.Name())
		}

		seenOptions[value.Name()] = struct{}{}
	}

	seenArguments := make(map[ArgumentName]struct{}, len(input.arguments))
	for index, argument := range input.arguments {
		if err := argument.Validate(); err != nil {
			return fmt.Errorf("%w: argument %d: %w", ErrInvalidBindingValue, index, err)
		}

		if _, exists := seenArguments[argument.Name()]; exists {
			return fmt.Errorf("%w: duplicate bound argument %q", ErrInvalidBindingValue, argument.Name())
		}

		seenArguments[argument.Name()] = struct{}{}
	}

	return nil
}
