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
	"strconv"
	"time"
)

// Bool parses the scalar value as bool.
func (value OptionValue) Bool() (bool, error) {
	raw, err := value.requireScalarKind(OptionKindBool)
	if err != nil {
		return false, err
	}

	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrInvalidOptionValue, err)
	}

	return parsed, nil
}

// Text returns the scalar value for string-like scalar kinds.
func (value OptionValue) Text() (string, error) {
	raw, err := value.requireStringLikeScalar()
	if err != nil {
		return "", err
	}

	return raw, nil
}

// Int parses the scalar value as int.
func (value OptionValue) Int() (int, error) {
	raw, err := value.requireScalarKind(OptionKindInt)
	if err != nil {
		return 0, err
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalidOptionValue, err)
	}

	return parsed, nil
}

// Int64 parses the scalar value as int64.
func (value OptionValue) Int64() (int64, error) {
	raw, err := value.requireScalarKind(OptionKindInt64)
	if err != nil {
		return 0, err
	}

	parsed, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalidOptionValue, err)
	}

	return parsed, nil
}

// Uint parses the scalar value as uint.
func (value OptionValue) Uint() (uint, error) {
	raw, err := value.requireScalarKind(OptionKindUint)
	if err != nil {
		return 0, err
	}

	parsed, err := strconv.ParseUint(raw, 10, strconv.IntSize)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalidOptionValue, err)
	}

	return uint(parsed), nil
}

// Uint64 parses the scalar value as uint64.
func (value OptionValue) Uint64() (uint64, error) {
	raw, err := value.requireScalarKind(OptionKindUint64)
	if err != nil {
		return 0, err
	}

	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalidOptionValue, err)
	}

	return parsed, nil
}

// Float64 parses the scalar value as float64.
func (value OptionValue) Float64() (float64, error) {
	raw, err := value.requireScalarKind(OptionKindFloat64)
	if err != nil {
		return 0, err
	}

	if err := validateOptionValueForKind(OptionKindFloat64, raw); err != nil {
		return 0, err
	}

	parsed, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalidOptionValue, err)
	}

	return parsed, nil
}

// Duration parses the scalar value as time.Duration.
func (value OptionValue) Duration() (time.Duration, error) {
	raw, err := value.requireScalarKind(OptionKindDuration)
	if err != nil {
		return 0, err
	}

	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalidOptionValue, err)
	}

	return parsed, nil
}

// TextValues returns string-like values for string or enum scalar/list kinds.
func (value OptionValue) TextValues() ([]string, error) {
	if !isStringLikeElementKind(value.kind.ElementKind()) {
		return nil, fmt.Errorf(
			"%w: kind %q is not string-like",
			ErrInvalidOptionValue,
			value.kind,
		)
	}

	return value.Values(), nil
}

// IntValues parses scalar or list values as int values.
func (value OptionValue) IntValues() ([]int, error) {
	return parseOptionValueList(value, OptionKindInt, "int-like", func(raw string) (int, error) {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return 0, err
		}

		return parsed, nil
	})
}

// Int64Values parses scalar or list values as int64 values.
func (value OptionValue) Int64Values() ([]int64, error) {
	return parseOptionValueList(value, OptionKindInt64, "int64-like", func(raw string) (int64, error) {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return 0, err
		}

		return parsed, nil
	})
}

// UintValues parses scalar or list values as uint values.
func (value OptionValue) UintValues() ([]uint, error) {
	return parseOptionValueList(value, OptionKindUint, "uint-like", func(raw string) (uint, error) {
		parsed, err := strconv.ParseUint(raw, 10, strconv.IntSize)
		if err != nil {
			return 0, err
		}

		return uint(parsed), nil
	})
}

// Uint64Values parses scalar or list values as uint64 values.
func (value OptionValue) Uint64Values() ([]uint64, error) {
	return parseOptionValueList(value, OptionKindUint64, "uint64-like", func(raw string) (uint64, error) {
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return 0, err
		}

		return parsed, nil
	})
}

// Float64Values parses scalar or list values as float64 values.
func (value OptionValue) Float64Values() ([]float64, error) {
	return parseOptionValueList(value, OptionKindFloat64, "float64-like", func(raw string) (float64, error) {
		if err := validateOptionValueForKind(OptionKindFloat64, raw); err != nil {
			return 0, err
		}

		parsed, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return 0, err
		}

		return parsed, nil
	})
}

// DurationValues parses scalar or list values as duration values.
func (value OptionValue) DurationValues() ([]time.Duration, error) {
	return parseOptionValueList(value, OptionKindDuration, "duration-like", func(raw string) (time.Duration, error) {
		parsed, err := time.ParseDuration(raw)
		if err != nil {
			return 0, err
		}

		return parsed, nil
	})
}
