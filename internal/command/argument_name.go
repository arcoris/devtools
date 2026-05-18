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

// ArgumentName is the canonical machine-facing name of a positional argument.
type ArgumentName string

// NewArgumentName validates raw and returns it as an ArgumentName.
func NewArgumentName(raw string) (ArgumentName, error) {
	name := ArgumentName(raw)
	if err := name.Validate(); err != nil {
		return "", err
	}

	return name, nil
}

// MustArgumentName validates raw and returns it as an ArgumentName.
//
// MustArgumentName panics on invalid input. It is intended for static command
// definitions and tests.
func MustArgumentName(raw string) ArgumentName {
	name, err := NewArgumentName(raw)
	if err != nil {
		panic(err)
	}

	return name
}

// String returns the canonical argument name.
func (name ArgumentName) String() string {
	return string(name)
}

// IsZero reports whether the argument name has not been set.
func (name ArgumentName) IsZero() bool {
	return name == ""
}

// IsValid reports whether the argument name satisfies the argument-name grammar.
func (name ArgumentName) IsValid() bool {
	return name.Validate() == nil
}

// Validate verifies argument-name structural rules.
func (name ArgumentName) Validate() error {
	raw := string(name)
	if raw == "" {
		return ErrEmptyArgumentName
	}

	if len(raw) > maxArgumentNameLength {
		return fmt.Errorf(
			"%w: length %d exceeds maximum length %d",
			ErrInvalidArgumentName,
			len(raw),
			maxArgumentNameLength,
		)
	}

	if err := validateCommandNameSegment(raw); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidArgumentName, err)
	}

	return nil
}
