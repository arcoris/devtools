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

import (
	"fmt"
	"strings"

	"arcoris.dev/devtools/internal/textvalidate"
)

// defaultArgumentMetavar derives a stable help placeholder from the argument
// name.
func defaultArgumentMetavar(name ArgumentName) string {
	if name == "" {
		return "VALUE"
	}

	return strings.ToUpper(strings.ReplaceAll(name.String(), "-", "_"))
}

// validateArgumentMetavar validates help/documentation value placeholder.
func validateArgumentMetavar(raw string) error {
	if err := textvalidate.ValidateTokenText(raw, maxArgumentMetavarLength); err != nil {
		return fmt.Errorf("%w: invalid metavar %q: %w", ErrInvalidArgument, raw, err)
	}

	return nil
}

// validateArgumentAllowedValues validates declaration allowed values.
func validateArgumentAllowedValues(kind OptionKind, values []string) error {
	if len(values) == 0 {
		if kind.RequiresAllowedValues() {
			return fmt.Errorf("%w: kind %q requires allowed values", ErrInvalidArgument, kind)
		}

		return nil
	}

	if !kind.CanHaveAllowedValues() {
		return fmt.Errorf("%w: kind %q must not declare allowed values", ErrInvalidArgument, kind)
	}

	seen := make(map[string]struct{}, len(values))
	for index, value := range values {
		if err := validateOptionAllowedValue(index, value); err != nil {
			return fmt.Errorf("%w: %w", ErrInvalidArgument, err)
		}

		if _, exists := seen[value]; exists {
			return fmt.Errorf("%w: duplicate allowed value %q", ErrInvalidArgument, value)
		}

		seen[value] = struct{}{}
	}

	return nil
}

// validateArgumentDefaultValues validates declaration defaults.
func validateArgumentDefaultValues(argument Argument) error {
	if len(argument.defaultValues) == 0 {
		return nil
	}

	if argument.requirement.IsRequired() {
		return fmt.Errorf("%w: required argument %q must not declare defaults", ErrInvalidArgument, argument.name)
	}

	if argument.cardinality.IsSingle() && len(argument.defaultValues) > 1 {
		return fmt.Errorf("%w: single argument %q must not have multiple defaults", ErrInvalidArgument, argument.name)
	}

	for index, value := range argument.defaultValues {
		if value == "" && argument.emptyValue.RejectsEmpty() {
			return fmt.Errorf("%w: default value %d must not be empty", ErrInvalidArgument, index)
		}

		if err := validateArgumentRawValue(fmt.Sprintf("default value %d", index), value); err != nil {
			return err
		}

		if err := validateOptionValueForKind(argument.kind, value); err != nil {
			return fmt.Errorf("%w: default value %d: %w", ErrInvalidArgument, index, err)
		}

		if !argument.AllowsValue(value) {
			return fmt.Errorf("%w: default value %q is not in allowed values", ErrInvalidArgument, value)
		}
	}

	return nil
}

// validateArgumentRawValue validates compact UTF-8 argument value text.
func validateArgumentRawValue(field string, raw string) error {
	if err := textvalidate.ValidateCompactText(raw, maxArgumentValueLength); err != nil {
		return fmt.Errorf("%w: invalid %s: %w", ErrInvalidArgument, field, err)
	}

	return nil
}

// normalizeArgumentValues returns a detached copy of values preserving
// declaration order.
func normalizeArgumentValues(values []string) []string {
	if values == nil {
		return nil
	}

	out := make([]string, len(values))
	copy(out, values)

	return out
}
