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

// Validate verifies argument declaration structural rules.
func (argument Argument) Validate() error {
	if err := argument.name.Validate(); err != nil {
		return err
	}

	if err := argument.kind.Validate(); err != nil {
		return fmt.Errorf("%w: invalid kind: %w", ErrInvalidArgument, err)
	}

	if argument.kind.IsList() {
		return fmt.Errorf(
			"%w: argument kind must be scalar, got list kind %q",
			ErrInvalidArgument,
			argument.kind,
		)
	}

	if err := argument.requirement.Validate(); err != nil {
		return err
	}

	if err := argument.cardinality.Validate(); err != nil {
		return err
	}

	if err := argument.emptyValue.Validate(); err != nil {
		return err
	}

	if argument.kind.IsBool() && argument.emptyValue.AllowsEmpty() {
		return fmt.Errorf("%w: bool argument must not allow empty explicit values", ErrInvalidArgument)
	}

	if !argument.kind.IsString() && !argument.kind.IsEnum() && argument.emptyValue.AllowsEmpty() {
		return fmt.Errorf("%w: empty explicit values are only valid for string-like arguments", ErrInvalidArgument)
	}

	if err := validateArgumentMetavar(argument.metavar); err != nil {
		return err
	}

	if err := argument.visibility.Validate(); err != nil {
		return fmt.Errorf("%w: invalid visibility: %w", ErrInvalidArgument, err)
	}

	if err := argument.documentation.Validate(); err != nil {
		return fmt.Errorf("%w: invalid documentation: %w", ErrInvalidArgument, err)
	}

	if err := argument.metadata.Validate(); err != nil {
		return fmt.Errorf("%w: invalid metadata: %w", ErrInvalidArgument, err)
	}

	if err := validateArgumentAllowedValues(argument.kind, argument.allowedValues); err != nil {
		return err
	}

	if err := validateArgumentDefaultValues(argument); err != nil {
		return err
	}

	return nil
}

// ValidateValues verifies runtime positional values against this declaration.
func (argument Argument) ValidateValues(values ...string) error {
	if err := argument.Validate(); err != nil {
		return err
	}

	if !argument.AcceptsCount(len(values)) {
		if maxValue, bounded := argument.MaxValues(); bounded {
			return fmt.Errorf(
				"%w: argument %q expects between %d and %d value(s), got %d",
				ErrInvalidArgumentValue,
				argument.name,
				argument.MinValues(),
				maxValue,
				len(values),
			)
		}

		return fmt.Errorf(
			"%w: argument %q expects at least %d value(s), got %d",
			ErrInvalidArgumentValue,
			argument.name,
			argument.MinValues(),
			len(values),
		)
	}

	for index, value := range values {
		if err := argument.validateOneValue(index, value); err != nil {
			return err
		}
	}

	return nil
}

// validateOneValue verifies one runtime value.
func (argument Argument) validateOneValue(index int, value string) error {
	if value == "" && argument.emptyValue.RejectsEmpty() {
		return fmt.Errorf("%w: argument %q value %d must not be empty", ErrInvalidArgumentValue, argument.name, index)
	}

	if err := validateArgumentRawValue(fmt.Sprintf("argument %q value %d", argument.name, index), value); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgumentValue, err)
	}

	if err := validateOptionValueForKind(argument.kind, value); err != nil {
		return fmt.Errorf("%w: argument %q value %d: %w", ErrInvalidArgumentValue, argument.name, index, err)
	}

	if !argument.AllowsValue(value) {
		return fmt.Errorf(
			"%w: argument %q value %q is not allowed",
			ErrInvalidArgumentValue,
			argument.name,
			value,
		)
	}

	return nil
}
