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

// requireScalarKind returns the raw scalar value after checking kind.
func (value OptionValue) requireScalarKind(kind OptionKind) (string, error) {
	if value.kind != kind {
		return "", fmt.Errorf("%w: expected kind %q, got %q", ErrInvalidOptionValue, kind, value.kind)
	}

	if len(value.values) != 1 {
		return "", fmt.Errorf("%w: expected exactly one value, got %d", ErrInvalidOptionValue, len(value.values))
	}

	return value.values[0], nil
}

// requireElementKind verifies the scalar kind of each stored raw value.
func (value OptionValue) requireElementKind(kind OptionKind, label string) error {
	if value.kind.ElementKind() != kind {
		return fmt.Errorf("%w: kind %q is not %s", ErrInvalidOptionValue, value.kind, label)
	}

	return nil
}

// requireStringLikeScalar returns a raw scalar string-like value.
func (value OptionValue) requireStringLikeScalar() (string, error) {
	if !isStringLikeElementKind(value.kind) {
		return "", fmt.Errorf("%w: kind %q is not string-like scalar", ErrInvalidOptionValue, value.kind)
	}

	if len(value.values) != 1 {
		return "", fmt.Errorf("%w: expected exactly one value, got %d", ErrInvalidOptionValue, len(value.values))
	}

	return value.values[0], nil
}

// parseOptionValueList parses scalar or list values after checking element kind.
func parseOptionValueList[T any](
	value OptionValue,
	kind OptionKind,
	label string,
	parse func(raw string) (T, error),
) ([]T, error) {
	if err := value.requireElementKind(kind, label); err != nil {
		return nil, err
	}

	out := make([]T, len(value.values))
	for index, raw := range value.values {
		parsed, err := parse(raw)
		if err != nil {
			return nil, fmt.Errorf("%w: value %d: %w", ErrInvalidOptionValue, index, err)
		}

		out[index] = parsed
	}

	return out, nil
}

// isStringLikeElementKind reports whether kind is a scalar string-like kind.
func isStringLikeElementKind(kind OptionKind) bool {
	switch kind {
	case OptionKindString, OptionKindEnum:
		return true
	default:
		return false
	}
}

// stringSlicesEqual reports whether two string slices contain the same values
// in the same order.
func stringSlicesEqual(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}

	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}

	return true
}
