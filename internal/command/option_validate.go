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

// Validate verifies option declaration structural rules.
func (option Option) Validate() error {
	if err := option.name.Validate(); err != nil {
		return err
	}

	if err := validateOptionAliases(option.name, option.aliases); err != nil {
		return err
	}

	if option.shorthand != "" {
		if err := validateOptionShorthand(option.shorthand); err != nil {
			return err
		}
	}

	if err := option.kind.Validate(); err != nil {
		return fmt.Errorf("%w: invalid kind: %w", ErrInvalidOption, err)
	}

	if err := validateOptionMetavar(option.metavar); err != nil {
		return err
	}

	if err := option.policy.ValidateForKind(option.kind); err != nil {
		return fmt.Errorf("%w: invalid policy: %w", ErrInvalidOption, err)
	}

	if err := option.documentation.Validate(); err != nil {
		return fmt.Errorf("%w: invalid documentation: %w", ErrInvalidOption, err)
	}

	if err := option.metadata.Validate(); err != nil {
		return fmt.Errorf("%w: invalid metadata: %w", ErrInvalidOption, err)
	}

	if err := option.visibility.Validate(); err != nil {
		return fmt.Errorf("%w: invalid visibility: %w", ErrInvalidOption, err)
	}

	if err := validateOptionAllowedValues(option.kind, option.allowedValues); err != nil {
		return err
	}

	if err := validateOptionDefaultValues(option.kind, option.policy, option.allowedValues, option.defaultValues); err != nil {
		return err
	}

	return nil
}

// validateOptionAliases validates alias names and duplicate/conflict rules.
func validateOptionAliases(name OptionName, aliases []OptionName) error {
	seen := map[OptionName]struct{}{
		name: {},
	}

	for _, alias := range aliases {
		if err := alias.Validate(); err != nil {
			return fmt.Errorf("%w: invalid alias %q: %w", ErrInvalidOption, alias, err)
		}

		if _, exists := seen[alias]; exists {
			return fmt.Errorf("%w: duplicate option name or alias %q", ErrInvalidOption, alias)
		}

		seen[alias] = struct{}{}
	}

	return nil
}
