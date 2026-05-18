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

// OptionEmptyValuePolicy describes whether an explicitly supplied empty value
// is accepted.
type OptionEmptyValuePolicy string

const (
	// OptionEmptyValueReject means an explicitly supplied empty value is invalid.
	OptionEmptyValueReject OptionEmptyValuePolicy = "reject-empty"

	// OptionEmptyValueAllow means an explicitly supplied empty value is allowed.
	OptionEmptyValueAllow OptionEmptyValuePolicy = "allow-empty"
)

// NewOptionEmptyValuePolicy validates raw and returns it as an
// OptionEmptyValuePolicy.
func NewOptionEmptyValuePolicy(raw string) (OptionEmptyValuePolicy, error) {
	policy := OptionEmptyValuePolicy(raw)
	if err := policy.Validate(); err != nil {
		return "", err
	}

	return policy, nil
}

// ParseOptionEmptyValuePolicy is an alias for NewOptionEmptyValuePolicy.
func ParseOptionEmptyValuePolicy(raw string) (OptionEmptyValuePolicy, error) {
	return NewOptionEmptyValuePolicy(raw)
}

// MustOptionEmptyValuePolicy validates raw and returns it as an
// OptionEmptyValuePolicy.
func MustOptionEmptyValuePolicy(raw string) OptionEmptyValuePolicy {
	policy, err := NewOptionEmptyValuePolicy(raw)
	if err != nil {
		panic(err)
	}

	return policy
}

// String returns the canonical string representation of the empty-value policy.
func (policy OptionEmptyValuePolicy) String() string {
	return string(policy)
}

// IsZero reports whether the empty-value policy has not been set.
func (policy OptionEmptyValuePolicy) IsZero() bool {
	return policy == ""
}

// OrDefault returns OptionEmptyValueReject when policy is zero.
func (policy OptionEmptyValuePolicy) OrDefault() OptionEmptyValuePolicy {
	if policy.IsZero() {
		return OptionEmptyValueReject
	}

	return policy
}

// IsKnown reports whether policy is one of the supported non-zero states.
func (policy OptionEmptyValuePolicy) IsKnown() bool {
	switch policy {
	case OptionEmptyValueReject, OptionEmptyValueAllow:
		return true
	default:
		return false
	}
}

// IsValid reports whether empty-value policy satisfies policy rules.
func (policy OptionEmptyValuePolicy) IsValid() bool {
	return policy.Validate() == nil
}

// Validate verifies that policy is a supported non-zero state.
func (policy OptionEmptyValuePolicy) Validate() error {
	if policy == "" {
		return fmt.Errorf("%w: empty-value policy is empty", ErrInvalidOptionPolicy)
	}

	if policy.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported empty-value policy %q", ErrInvalidOptionPolicy, policy)
}

// AllowsEmpty reports whether explicitly supplied empty values are allowed.
func (policy OptionEmptyValuePolicy) AllowsEmpty() bool {
	return policy == OptionEmptyValueAllow
}

// RejectsEmpty reports whether explicitly supplied empty values are rejected.
func (policy OptionEmptyValuePolicy) RejectsEmpty() bool {
	return policy == OptionEmptyValueReject
}
