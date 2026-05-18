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

// OrDefault returns ArgumentRequirementRequired when requirement is zero.
func (requirement ArgumentRequirement) OrDefault() ArgumentRequirement {
	if requirement == "" {
		return ArgumentRequirementRequired
	}

	return requirement
}

// String returns the canonical requirement string.
func (requirement ArgumentRequirement) String() string {
	return string(requirement)
}

// IsKnown reports whether requirement is a supported non-zero state.
func (requirement ArgumentRequirement) IsKnown() bool {
	switch requirement {
	case ArgumentRequirementRequired, ArgumentRequirementOptional:
		return true
	default:
		return false
	}
}

// Validate verifies requirement.
func (requirement ArgumentRequirement) Validate() error {
	if requirement == "" {
		return fmt.Errorf("%w: requirement is empty", ErrInvalidArgument)
	}

	if requirement.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported requirement %q", ErrInvalidArgument, requirement)
}

// IsRequired reports whether requirement is required.
func (requirement ArgumentRequirement) IsRequired() bool {
	return requirement == ArgumentRequirementRequired
}

// IsOptional reports whether requirement is optional.
func (requirement ArgumentRequirement) IsOptional() bool {
	return requirement == ArgumentRequirementOptional
}
