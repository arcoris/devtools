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

// Validate verifies resolved option value structural rules.
func (value OptionValue) Validate() error {
	if err := value.name.Validate(); err != nil {
		return fmt.Errorf("%w: invalid name: %w", ErrInvalidOptionValue, err)
	}

	if err := value.kind.Validate(); err != nil {
		return fmt.Errorf("%w: invalid kind: %w", ErrInvalidOptionValue, err)
	}

	if err := value.source.Validate(); err != nil {
		return fmt.Errorf("%w: invalid source: %w", ErrInvalidOptionValue, err)
	}

	if len(value.values) == 0 {
		return fmt.Errorf("%w: values must not be empty", ErrInvalidOptionValue)
	}

	if value.kind.IsScalar() && len(value.values) != 1 {
		return fmt.Errorf(
			"%w: scalar kind %q must have exactly one value, got %d",
			ErrInvalidOptionValue,
			value.kind,
			len(value.values),
		)
	}

	elementKind := value.kind.ElementKind()
	for index, raw := range value.values {
		if err := validateOptionRawValue(fmt.Sprintf("value %d", index), raw); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidOptionValue, err)
		}

		if err := validateOptionValueForKind(elementKind, raw); err != nil {
			return fmt.Errorf("%w: value %d: %w", ErrInvalidOptionValue, index, err)
		}

		if elementKind == OptionKindEnum {
			if err := validateOptionAllowedValue(index, raw); err != nil {
				return fmt.Errorf("%w: %w", ErrInvalidOptionValue, err)
			}
		}
	}

	return nil
}

// validateOptionValueSource verifies source validity and option policy.
func validateOptionValueSource(option Option, source OptionSource) error {
	if err := source.Validate(); err != nil {
		return fmt.Errorf("%w: invalid source: %w", ErrInvalidOptionValue, err)
	}

	if !option.Policy().AllowsSource(source) {
		return fmt.Errorf(
			"%w: source %q is not allowed for option %q",
			ErrInvalidOptionValue,
			source,
			option.Name(),
		)
	}

	return nil
}

// resolveOptionValueInputs resolves default-source value omission semantics.
func resolveOptionValueInputs(option Option, source OptionSource, values []string) ([]string, error) {
	resolvedValues := cloneStringSlice(values)
	if !source.IsDefault() {
		if len(resolvedValues) == 0 {
			return nil, fmt.Errorf("%w: values must not be empty", ErrInvalidOptionValue)
		}

		return resolvedValues, nil
	}

	defaultValues := option.DefaultValues()
	if len(defaultValues) == 0 {
		return nil, fmt.Errorf(
			"%w: source %q requires declared default values for option %q",
			ErrInvalidOptionValue,
			source,
			option.Name(),
		)
	}

	if len(resolvedValues) == 0 {
		return defaultValues, nil
	}

	if !stringSlicesEqual(resolvedValues, defaultValues) {
		return nil, fmt.Errorf(
			"%w: default-source values for option %q must match declaration defaults",
			ErrInvalidOptionValue,
			option.Name(),
		)
	}

	return resolvedValues, nil
}

// validateOptionValueAgainstOption validates raw values against option policy
// and allowed-value declarations.
func validateOptionValueAgainstOption(option Option, values []string) error {
	for index, raw := range values {
		if raw == "" && option.Policy().EmptyValue().RejectsEmpty() {
			return fmt.Errorf("%w: value %d must not be empty", ErrInvalidOptionValue, index)
		}

		if !option.AllowsValue(raw) {
			return fmt.Errorf(
				"%w: value %q is not allowed for option %q",
				ErrInvalidOptionValue,
				raw,
				option.Name(),
			)
		}
	}

	return nil
}
