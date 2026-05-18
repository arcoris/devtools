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

// DefaultOptionPolicy returns the default policy for ordinary command options.
func DefaultOptionPolicy() OptionPolicy {
	return MustOptionPolicy(OptionPolicySpec{})
}

// NewOptionPolicy validates spec and returns OptionPolicy.
func NewOptionPolicy(spec OptionPolicySpec) (OptionPolicy, error) {
	policy := OptionPolicy{
		requirement:    spec.Requirement.OrDefault(),
		scope:          spec.Scope.OrDefault(),
		occurrence:     spec.Occurrence,
		emptyValue:     spec.EmptyValue.OrDefault(),
		allowedSources: normalizeOptionSources(spec.AllowedSources),
	}

	if policy.occurrence.IsZero() {
		policy.occurrence = OptionOccurrenceSingle
	}

	if len(policy.allowedSources) == 0 {
		policy.allowedSources = KnownOptionSources()
	}

	if err := policy.Validate(); err != nil {
		return OptionPolicy{}, err
	}

	return policy, nil
}

// NewOptionPolicyForKind validates spec against kind and returns OptionPolicy.
//
// This constructor is preferred when the concrete option kind is already known,
// because occurrence defaults can be derived from the value shape.
func NewOptionPolicyForKind(kind OptionKind, spec OptionPolicySpec) (OptionPolicy, error) {
	policy := OptionPolicy{
		requirement:    spec.Requirement.OrDefault(),
		scope:          spec.Scope.OrDefault(),
		occurrence:     spec.Occurrence.OrDefaultForKind(kind),
		emptyValue:     spec.EmptyValue.OrDefault(),
		allowedSources: normalizeOptionSources(spec.AllowedSources),
	}

	if len(policy.allowedSources) == 0 {
		policy.allowedSources = KnownOptionSources()
	}

	if err := policy.ValidateForKind(kind); err != nil {
		return OptionPolicy{}, err
	}

	return policy, nil
}

// MustOptionPolicy validates spec and returns OptionPolicy.
//
// MustOptionPolicy panics on invalid input. It is intended for static option
// declarations and tests.
func MustOptionPolicy(spec OptionPolicySpec) OptionPolicy {
	policy, err := NewOptionPolicy(spec)
	if err != nil {
		panic(err)
	}

	return policy
}

// MustOptionPolicyForKind validates spec against kind and returns OptionPolicy.
//
// MustOptionPolicyForKind panics on invalid input. It is intended for static
// option declarations and tests.
func MustOptionPolicyForKind(kind OptionKind, spec OptionPolicySpec) OptionPolicy {
	policy, err := NewOptionPolicyForKind(kind, spec)
	if err != nil {
		panic(err)
	}

	return policy
}
