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

// Validate verifies binding declaration structural rules.
func (binding Binding) Validate() error {
	if err := validateBindingOptions(binding.options); err != nil {
		return err
	}

	if err := validateBindingArguments(binding.arguments); err != nil {
		return err
	}

	return nil
}

// ValidateOptionValues verifies resolved option values against declared options.
func (binding Binding) ValidateOptionValues(values ...OptionValue) error {
	_, err := binding.bindOptionValues(values)

	return err
}

// ValidatePositionalValues verifies positional values against declared
// arguments.
func (binding Binding) ValidatePositionalValues(values ...string) error {
	_, err := binding.bindPositionalValues(values)

	return err
}

// validateBindingOptions validates option declarations and uniqueness rules.
func validateBindingOptions(options []Option) error {
	seenNames := make(map[OptionName]struct{})
	seenShorthands := make(map[string]OptionName)

	for index, option := range options {
		if err := option.Validate(); err != nil {
			return fmt.Errorf("%w: option %d: %w", ErrInvalidBinding, index, err)
		}

		allNames := append([]OptionName{option.Name()}, option.Aliases()...)
		for _, name := range allNames {
			if owner, exists := seenNames[name]; exists {
				_ = owner
				return fmt.Errorf("%w: duplicate option name or alias %q", ErrInvalidBinding, name)
			}

			seenNames[name] = struct{}{}
		}

		if option.HasShorthand() {
			if owner, exists := seenShorthands[option.Shorthand()]; exists {
				return fmt.Errorf(
					"%w: duplicate shorthand %q used by options %q and %q",
					ErrInvalidBinding,
					option.Shorthand(),
					owner,
					option.Name(),
				)
			}

			seenShorthands[option.Shorthand()] = option.Name()
		}
	}

	return nil
}

// validateBindingArguments validates positional argument declarations and
// consumption-order rules.
func validateBindingArguments(arguments []Argument) error {
	seen := make(map[ArgumentName]struct{}, len(arguments))
	optionalSeen := false

	for index, argument := range arguments {
		if err := argument.Validate(); err != nil {
			return fmt.Errorf("%w: argument %d: %w", ErrInvalidBinding, index, err)
		}

		if _, exists := seen[argument.Name()]; exists {
			return fmt.Errorf("%w: duplicate argument name %q", ErrInvalidBinding, argument.Name())
		}

		seen[argument.Name()] = struct{}{}

		if optionalSeen && argument.IsRequired() {
			return fmt.Errorf(
				"%w: required argument %q must not follow optional argument",
				ErrInvalidBinding,
				argument.Name(),
			)
		}

		if argument.IsOptional() {
			optionalSeen = true
		}

		if argument.IsVariadic() && index != len(arguments)-1 {
			return fmt.Errorf(
				"%w: variadic argument %q must be the last argument",
				ErrInvalidBinding,
				argument.Name(),
			)
		}
	}

	return nil
}
