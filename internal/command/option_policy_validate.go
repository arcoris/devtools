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

// Validate verifies policy structural rules without checking a concrete kind.
func (policy OptionPolicy) Validate() error {
	if err := policy.requirement.Validate(); err != nil {
		return err
	}

	if err := policy.scope.Validate(); err != nil {
		return err
	}

	if err := policy.occurrence.Validate(); err != nil {
		return err
	}

	if err := policy.emptyValue.Validate(); err != nil {
		return err
	}

	return validateAllowedOptionSources(policy.allowedSources)
}

// ValidateForKind verifies policy structural rules and consistency with kind.
func (policy OptionPolicy) ValidateForKind(kind OptionKind) error {
	if err := kind.Validate(); err != nil {
		return fmt.Errorf("%w: invalid option kind: %w", ErrInvalidOptionPolicy, err)
	}

	if err := policy.Validate(); err != nil {
		return err
	}

	if kind.IsList() && policy.occurrence.IsSingle() {
		return fmt.Errorf(
			"%w: list option kind %q should use multiple occurrence policy",
			ErrInvalidOptionPolicy,
			kind,
		)
	}

	if kind.IsBool() && policy.emptyValue.AllowsEmpty() {
		return fmt.Errorf(
			"%w: bool option kind must not allow empty explicit values",
			ErrInvalidOptionPolicy,
		)
	}

	if !kind.IsString() && !kind.IsList() && policy.emptyValue.AllowsEmpty() {
		return fmt.Errorf(
			"%w: empty explicit values are only valid for string-like options",
			ErrInvalidOptionPolicy,
		)
	}

	return nil
}

// validateAllowedOptionSources validates source list semantics.
func validateAllowedOptionSources(sources []OptionSource) error {
	if len(sources) == 0 {
		return fmt.Errorf("%w: allowed sources must not be empty", ErrInvalidOptionPolicy)
	}

	seen := make(map[OptionSource]struct{}, len(sources))
	for index, source := range sources {
		if err := source.Validate(); err != nil {
			return fmt.Errorf("%w: allowed source %d: %w", ErrInvalidOptionPolicy, index, err)
		}

		if _, exists := seen[source]; exists {
			return fmt.Errorf("%w: duplicate allowed source %q", ErrInvalidOptionPolicy, source)
		}

		seen[source] = struct{}{}
	}

	return nil
}
