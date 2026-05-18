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

// Bind validates runtime values, applies declaration defaults, canonicalizes
// option aliases to canonical names, and returns a bound input value.
//
// Bind does not mutate the Binding or input values.
func (binding Binding) Bind(spec BindingValueSpec) (BoundInput, error) {
	if err := binding.Validate(); err != nil {
		return BoundInput{}, err
	}

	optionValues, err := binding.bindOptionValues(spec.OptionValues)
	if err != nil {
		return BoundInput{}, err
	}

	arguments, err := binding.bindPositionalValues(spec.PositionalValues)
	if err != nil {
		return BoundInput{}, err
	}

	bound := BoundInput{
		options:   optionValues,
		arguments: arguments,
	}

	if err := bound.Validate(); err != nil {
		return BoundInput{}, err
	}

	return bound, nil
}

// MustBind validates runtime values and returns a BoundInput.
//
// MustBind panics on invalid input. It is intended for tests and controlled
// internal wiring.
func (binding Binding) MustBind(spec BindingValueSpec) BoundInput {
	bound, err := binding.Bind(spec)
	if err != nil {
		panic(err)
	}

	return bound
}

// bindOptionValues validates, canonicalizes, and default-fills option values.
func (binding Binding) bindOptionValues(values []OptionValue) ([]OptionValue, error) {
	if err := validateBindingOptions(binding.options); err != nil {
		return nil, err
	}

	byName := make(map[OptionName]OptionValue, len(values))

	for index, value := range values {
		if err := value.Validate(); err != nil {
			return nil, fmt.Errorf("%w: option value %d: %w", ErrInvalidBindingValue, index, err)
		}

		option, ok := binding.Option(value.Name())
		if !ok {
			return nil, fmt.Errorf(
				"%w: unknown option %q",
				ErrInvalidBindingValue,
				value.Name(),
			)
		}

		if value.Kind() != option.Kind() {
			return nil, fmt.Errorf(
				"%w: option %q has kind %q, got value kind %q",
				ErrInvalidBindingValue,
				option.Name(),
				option.Kind(),
				value.Kind(),
			)
		}

		canonical, err := NewOptionValueFromOption(option, value.Source(), value.Values()...)
		if err != nil {
			return nil, fmt.Errorf("%w: option %q: %w", ErrInvalidBindingValue, option.Name(), err)
		}

		if _, exists := byName[option.Name()]; exists {
			return nil, fmt.Errorf(
				"%w: duplicate value for option %q",
				ErrInvalidBindingValue,
				option.Name(),
			)
		}

		byName[option.Name()] = canonical
	}

	for _, option := range binding.options {
		if _, exists := byName[option.Name()]; exists {
			continue
		}

		if option.HasDefault() {
			value, err := NewDefaultOptionValue(option)
			if err != nil {
				return nil, fmt.Errorf("%w: option %q default: %w", ErrInvalidBindingValue, option.Name(), err)
			}

			byName[option.Name()] = value

			continue
		}

		if option.IsRequired() {
			return nil, fmt.Errorf(
				"%w: required option %q is missing",
				ErrInvalidBindingValue,
				option.Name(),
			)
		}
	}

	out := make([]OptionValue, 0, len(byName))
	for _, option := range binding.options {
		if value, exists := byName[option.Name()]; exists {
			out = append(out, value)
		}
	}

	return out, nil
}

// bindPositionalValues validates, default-fills, and names positional values.
func (binding Binding) bindPositionalValues(values []string) ([]BoundArgument, error) {
	if err := validateBindingArguments(binding.arguments); err != nil {
		return nil, err
	}

	values = cloneBindingStrings(values)

	if len(binding.arguments) == 0 {
		if len(values) != 0 {
			return nil, fmt.Errorf(
				"%w: command accepts no positional arguments, got %d",
				ErrInvalidBindingValue,
				len(values),
			)
		}

		return nil, nil
	}

	out := make([]BoundArgument, 0, len(binding.arguments))
	cursor := 0

	for _, argument := range binding.arguments {
		var consumed []string

		if argument.IsVariadic() {
			consumed = cloneBindingStrings(values[cursor:])
			cursor = len(values)
		} else if cursor < len(values) {
			consumed = []string{values[cursor]}
			cursor++
		}

		if len(consumed) == 0 && argument.HasDefault() {
			consumed = argument.DefaultValues()
		}

		if err := argument.ValidateValues(consumed...); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidBindingValue, err)
		}

		out = append(out, BoundArgument{
			argument: argument,
			values:   cloneBindingStrings(consumed),
		})
	}

	if cursor < len(values) {
		return nil, fmt.Errorf(
			"%w: too many positional arguments: expected at most %d declaration slot(s), got %d value(s)",
			ErrInvalidBindingValue,
			len(binding.arguments),
			len(values),
		)
	}

	return out, nil
}
