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

// Spec returns a detached construction spec.
func (policy OptionPolicy) Spec() OptionPolicySpec {
	return policy.spec()
}

// IsZero reports whether policy has not been initialized.
func (policy OptionPolicy) IsZero() bool {
	return isZeroOptionPolicy(policy)
}

// IsValid reports whether policy satisfies structural rules.
func (policy OptionPolicy) IsValid() bool {
	return policy.Validate() == nil
}

// Requirement returns the option requirement policy.
func (policy OptionPolicy) Requirement() OptionRequirement {
	return policy.requirement
}

// Scope returns the option declaration scope.
func (policy OptionPolicy) Scope() OptionScope {
	return policy.scope
}

// Occurrence returns the option occurrence policy.
func (policy OptionPolicy) Occurrence() OptionOccurrence {
	return policy.occurrence
}

// EmptyValue returns the empty-value policy.
func (policy OptionPolicy) EmptyValue() OptionEmptyValuePolicy {
	return policy.emptyValue
}

// AllowedSources returns allowed option sources in default precedence order.
func (policy OptionPolicy) AllowedSources() []OptionSource {
	return cloneOptionSources(policy.allowedSources)
}

// AllowedSource returns one allowed source by precedence-order index.
func (policy OptionPolicy) AllowedSource(index int) (OptionSource, bool) {
	if index < 0 || index >= len(policy.allowedSources) {
		return "", false
	}

	return policy.allowedSources[index], true
}

// AllowedSourceCount returns the number of allowed sources.
func (policy OptionPolicy) AllowedSourceCount() int {
	return len(policy.allowedSources)
}

// IsRequired reports whether a resolved value is required.
func (policy OptionPolicy) IsRequired() bool {
	return policy.requirement.IsRequired()
}

// IsOptional reports whether a resolved value may be omitted.
func (policy OptionPolicy) IsOptional() bool {
	return policy.requirement.IsOptional()
}

// IsLocal reports whether the option applies only to the declaring node.
func (policy OptionPolicy) IsLocal() bool {
	return policy.scope.IsLocal()
}

// IsInheritedByChildren reports whether the option is visible to descendants.
func (policy OptionPolicy) IsInheritedByChildren() bool {
	return policy.scope.IsInheritedByChildren()
}

// IsRepeatable reports whether the option may appear multiple times.
func (policy OptionPolicy) IsRepeatable() bool {
	return policy.occurrence.IsMultiple()
}

// AllowsEmptyValue reports whether explicitly supplied empty values are allowed.
func (policy OptionPolicy) AllowsEmptyValue() bool {
	return policy.emptyValue.AllowsEmpty()
}

// AllowsSource reports whether source is allowed to provide a value.
func (policy OptionPolicy) AllowsSource(source OptionSource) bool {
	for _, candidate := range policy.allowedSources {
		if candidate == source {
			return true
		}
	}

	return false
}

// AllowsOnlyExplicitSources reports whether every allowed source is explicit.
func (policy OptionPolicy) AllowsOnlyExplicitSources() bool {
	if len(policy.allowedSources) == 0 {
		return false
	}

	for _, source := range policy.allowedSources {
		if !source.IsExplicit() {
			return false
		}
	}

	return true
}

// AllowsDefaultSource reports whether declaration defaults may provide a value.
func (policy OptionPolicy) AllowsDefaultSource() bool {
	return policy.AllowsSource(OptionSourceDefault)
}

// HighestAllowedSource returns the allowed source with the highest default
// precedence.
func (policy OptionPolicy) HighestAllowedSource() (OptionSource, bool) {
	if len(policy.allowedSources) == 0 {
		return "", false
	}

	highest := policy.allowedSources[0]
	for _, source := range policy.allowedSources[1:] {
		if source.Overrides(highest) {
			highest = source
		}
	}

	return highest, true
}

// LowestAllowedSource returns the allowed source with the lowest default
// precedence.
func (policy OptionPolicy) LowestAllowedSource() (OptionSource, bool) {
	if len(policy.allowedSources) == 0 {
		return "", false
	}

	lowest := policy.allowedSources[0]
	for _, source := range policy.allowedSources[1:] {
		if lowest.Overrides(source) {
			lowest = source
		}
	}

	return lowest, true
}

// spec returns a detached construction spec.
func (policy OptionPolicy) spec() OptionPolicySpec {
	return OptionPolicySpec{
		Requirement:    policy.requirement,
		Scope:          policy.scope,
		Occurrence:     policy.occurrence,
		EmptyValue:     policy.emptyValue,
		AllowedSources: policy.AllowedSources(),
	}
}
