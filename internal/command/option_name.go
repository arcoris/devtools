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

// OptionName is the canonical long name of a command option without the leading
// "--".
//
// OptionName is a stable machine-facing value. It is intentionally strict so it
// remains safe in CLI flags, generated documentation anchors, config keys,
// JSON fields, reports, and diagnostics.
type OptionName string

// NewOptionName validates raw and returns it as an OptionName.
func NewOptionName(raw string) (OptionName, error) {
	name := OptionName(raw)
	if err := name.Validate(); err != nil {
		return "", err
	}

	return name, nil
}

// ParseOptionName is an alias for NewOptionName.
func ParseOptionName(raw string) (OptionName, error) {
	return NewOptionName(raw)
}

// MustOptionName validates raw and returns it as an OptionName.
//
// MustOptionName panics on invalid input. It is intended for static option
// declarations and tests.
func MustOptionName(raw string) OptionName {
	name, err := NewOptionName(raw)
	if err != nil {
		panic(err)
	}

	return name
}

// String returns the canonical option name.
func (name OptionName) String() string {
	return string(name)
}

// LongFlag returns the CLI long flag spelling.
func (name OptionName) LongFlag() string {
	if name == "" {
		return ""
	}

	return "--" + string(name)
}

// IsZero reports whether the option name has not been set.
func (name OptionName) IsZero() bool {
	return name == ""
}

// IsValid reports whether the option name satisfies the option-name grammar.
func (name OptionName) IsValid() bool {
	return name.Validate() == nil
}

// Validate verifies option-name structural rules.
func (name OptionName) Validate() error {
	raw := string(name)
	if raw == "" {
		return ErrEmptyOptionName
	}

	if len(raw) > maxOptionNameLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidOptionName,
			len(raw),
			maxOptionNameLength,
		)
	}

	if err := validateCommandNameSegment(raw); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidOptionName, err)
	}

	return nil
}

// newOptionNames validates raw option names.
func newOptionNames(values []string) ([]OptionName, error) {
	names := make([]OptionName, 0, len(values))

	for index, raw := range values {
		name, err := NewOptionName(raw)
		if err != nil {
			return nil, fmt.Errorf("%w: alias %d: %w", ErrInvalidOption, index, err)
		}

		names = append(names, name)
	}

	return names, nil
}
