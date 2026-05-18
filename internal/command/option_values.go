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
	"math"
	"strconv"
	"time"

	"arcoris.dev/devtools/internal/textvalidate"
)

// validateOptionShorthand validates one optional short flag name.
//
// Shorthand names are intentionally ASCII-only and limited to one letter or
// digit so every adapter can safely map them to "-x" style flags.
func validateOptionShorthand(raw string) error {
	if raw == "" {
		return nil
	}

	if len(raw) != maxOptionShorthandLength {
		return fmt.Errorf("%w: shorthand %q must be exactly one ASCII character", ErrInvalidOption, raw)
	}

	ch := raw[0]
	if textvalidate.IsASCIILowerLetter(ch) ||
		textvalidate.IsASCIIUpperLetter(ch) ||
		textvalidate.IsASCIIDigit(ch) {
		return nil
	}

	return fmt.Errorf("%w: shorthand %q must be ASCII letter or digit", ErrInvalidOption, raw)
}

// validateOptionMetavar validates help/documentation value placeholder.
func validateOptionMetavar(raw string) error {
	if raw == "" {
		return fmt.Errorf("%w: metavar must not be empty", ErrInvalidOption)
	}

	if err := textvalidate.ValidateTokenText(raw, maxOptionMetavarLength); err != nil {
		return fmt.Errorf("%w: invalid metavar %q: %w", ErrInvalidOption, raw, err)
	}

	return nil
}

// validateOptionAllowedValues validates option allowed-value declarations.
func validateOptionAllowedValues(kind OptionKind, values []string) error {
	if len(values) == 0 {
		if kind.RequiresAllowedValues() {
			return fmt.Errorf("%w: kind %q requires allowed values", ErrInvalidOption, kind)
		}

		return nil
	}

	if !kind.CanHaveAllowedValues() {
		return fmt.Errorf("%w: kind %q must not declare allowed values", ErrInvalidOption, kind)
	}

	seen := make(map[string]struct{}, len(values))
	for index, value := range values {
		if err := validateOptionAllowedValue(index, value); err != nil {
			return err
		}

		if _, exists := seen[value]; exists {
			return fmt.Errorf("%w: duplicate allowed value %q", ErrInvalidOption, value)
		}

		seen[value] = struct{}{}
	}

	return nil
}

// validateOptionAllowedValue validates one allowed value.
//
// Allowed values are machine-facing enum-like values. They intentionally use
// command-name segment grammar rather than arbitrary free text.
func validateOptionAllowedValue(index int, value string) error {
	if err := validateCommandNameSegment(value); err != nil {
		return fmt.Errorf("%w: allowed value %d: %w", ErrInvalidOption, index, err)
	}

	return nil
}

// validateOptionDefaultValues validates declaration defaults against kind,
// policy, and allowed values.
func validateOptionDefaultValues(kind OptionKind, policy OptionPolicy, allowedValues []string, defaults []string) error {
	if len(defaults) == 0 {
		return nil
	}

	if !policy.AllowsDefaultSource() {
		return fmt.Errorf("%w: default values are declared but default source is not allowed", ErrInvalidOption)
	}

	if kind.IsScalar() && len(defaults) > 1 {
		return fmt.Errorf("%w: scalar kind %q must not have multiple default values", ErrInvalidOption, kind)
	}

	for index, value := range defaults {
		if value == "" && policy.EmptyValue().RejectsEmpty() {
			return fmt.Errorf("%w: default value %d must not be empty", ErrInvalidOption, index)
		}

		if err := validateOptionRawValue(fmt.Sprintf("default value %d", index), value); err != nil {
			return err
		}

		if err := validateOptionValueForKind(kind.ElementKind(), value); err != nil {
			return fmt.Errorf("%w: default value %d: %w", ErrInvalidOption, index, err)
		}

		if len(allowedValues) != 0 && !containsString(allowedValues, value) {
			return fmt.Errorf("%w: default value %q is not in allowed values", ErrInvalidOption, value)
		}
	}

	return nil
}

// validateOptionRawValue validates compact UTF-8 option value text.
func validateOptionRawValue(field string, raw string) error {
	if err := textvalidate.ValidateCompactText(raw, maxOptionValueLength); err != nil {
		return fmt.Errorf("%w: %w: %s: %w", ErrInvalidOption, ErrInvalidOptionValue, field, err)
	}

	return nil
}

// validateOptionValueForKind validates one raw option value against a scalar
// OptionKind.
func validateOptionValueForKind(kind OptionKind, raw string) error {
	switch kind {
	case OptionKindBool:
		if raw == "true" || raw == "false" {
			return nil
		}

		return fmt.Errorf("%w: expected bool value true or false", ErrInvalidOptionValue)

	case OptionKindString, OptionKindEnum:
		return nil

	case OptionKindInt:
		_, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("%w: expected int value: %w", ErrInvalidOptionValue, err)
		}

		return nil

	case OptionKindInt64:
		_, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: expected int64 value: %w", ErrInvalidOptionValue, err)
		}

		return nil

	case OptionKindUint:
		_, err := strconv.ParseUint(raw, 10, strconv.IntSize)
		if err != nil {
			return fmt.Errorf("%w: expected uint value: %w", ErrInvalidOptionValue, err)
		}

		return nil

	case OptionKindUint64:
		_, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: expected uint64 value: %w", ErrInvalidOptionValue, err)
		}

		return nil

	case OptionKindFloat64:
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return fmt.Errorf("%w: expected float64 value: %w", ErrInvalidOptionValue, err)
		}

		if math.IsNaN(value) || math.IsInf(value, 0) {
			return fmt.Errorf("%w: expected finite float64 value", ErrInvalidOptionValue)
		}

		return nil

	case OptionKindDuration:
		_, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("%w: expected duration value: %w", ErrInvalidOptionValue, err)
		}

		return nil

	default:
		return fmt.Errorf("%w: unsupported scalar kind %q", ErrInvalidOptionValue, kind)
	}
}
