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

// OptionRequirement describes whether an option value must be present after
// option resolution.
type OptionRequirement string

const (
	// OptionRequirementOptional means the option may be omitted.
	OptionRequirementOptional OptionRequirement = "optional"

	// OptionRequirementRequired means the option must have a resolved value.
	OptionRequirementRequired OptionRequirement = "required"
)

// NewOptionRequirement validates raw and returns it as an OptionRequirement.
func NewOptionRequirement(raw string) (OptionRequirement, error) {
	requirement := OptionRequirement(raw)
	if err := requirement.Validate(); err != nil {
		return "", err
	}

	return requirement, nil
}

// ParseOptionRequirement is an alias for NewOptionRequirement.
func ParseOptionRequirement(raw string) (OptionRequirement, error) {
	return NewOptionRequirement(raw)
}

// MustOptionRequirement validates raw and returns it as an OptionRequirement.
func MustOptionRequirement(raw string) OptionRequirement {
	requirement, err := NewOptionRequirement(raw)
	if err != nil {
		panic(err)
	}

	return requirement
}

// String returns the canonical string representation of the requirement.
func (requirement OptionRequirement) String() string {
	return string(requirement)
}

// IsZero reports whether the requirement has not been set.
func (requirement OptionRequirement) IsZero() bool {
	return requirement == ""
}

// OrDefault returns OptionRequirementOptional when requirement is zero.
func (requirement OptionRequirement) OrDefault() OptionRequirement {
	if requirement.IsZero() {
		return OptionRequirementOptional
	}

	return requirement
}

// IsKnown reports whether requirement is one of the supported non-zero states.
func (requirement OptionRequirement) IsKnown() bool {
	switch requirement {
	case OptionRequirementOptional, OptionRequirementRequired:
		return true
	default:
		return false
	}
}

// IsValid reports whether requirement satisfies policy rules.
func (requirement OptionRequirement) IsValid() bool {
	return requirement.Validate() == nil
}

// Validate verifies that requirement is a supported non-zero state.
func (requirement OptionRequirement) Validate() error {
	if requirement == "" {
		return fmt.Errorf("%w: requirement is empty", ErrInvalidOptionPolicy)
	}

	if requirement.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported requirement %q", ErrInvalidOptionPolicy, requirement)
}

// IsOptional reports whether the option is optional.
func (requirement OptionRequirement) IsOptional() bool {
	return requirement == OptionRequirementOptional
}

// IsRequired reports whether the option is required.
func (requirement OptionRequirement) IsRequired() bool {
	return requirement == OptionRequirementRequired
}
