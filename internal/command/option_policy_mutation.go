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

// WithRequirement returns a validated copy with requirement replaced.
func (policy OptionPolicy) WithRequirement(requirement OptionRequirement) (OptionPolicy, error) {
	spec := policy.spec()
	spec.Requirement = requirement

	return NewOptionPolicy(spec)
}

// MustWithRequirement returns a validated copy with requirement replaced and
// panics on invalid input.
func (policy OptionPolicy) MustWithRequirement(requirement OptionRequirement) OptionPolicy {
	next, err := policy.WithRequirement(requirement)
	if err != nil {
		panic(err)
	}

	return next
}

// WithScope returns a validated copy with scope replaced.
func (policy OptionPolicy) WithScope(scope OptionScope) (OptionPolicy, error) {
	spec := policy.spec()
	spec.Scope = scope

	return NewOptionPolicy(spec)
}

// MustWithScope returns a validated copy with scope replaced and panics on
// invalid input.
func (policy OptionPolicy) MustWithScope(scope OptionScope) OptionPolicy {
	next, err := policy.WithScope(scope)
	if err != nil {
		panic(err)
	}

	return next
}

// WithOccurrence returns a validated copy with occurrence replaced.
func (policy OptionPolicy) WithOccurrence(occurrence OptionOccurrence) (OptionPolicy, error) {
	spec := policy.spec()
	spec.Occurrence = occurrence

	return NewOptionPolicy(spec)
}

// MustWithOccurrence returns a validated copy with occurrence replaced and
// panics on invalid input.
func (policy OptionPolicy) MustWithOccurrence(occurrence OptionOccurrence) OptionPolicy {
	next, err := policy.WithOccurrence(occurrence)
	if err != nil {
		panic(err)
	}

	return next
}

// WithEmptyValue returns a validated copy with empty-value policy replaced.
func (policy OptionPolicy) WithEmptyValue(emptyValue OptionEmptyValuePolicy) (OptionPolicy, error) {
	spec := policy.spec()
	spec.EmptyValue = emptyValue

	return NewOptionPolicy(spec)
}

// MustWithEmptyValue returns a validated copy with empty-value policy replaced
// and panics on invalid input.
func (policy OptionPolicy) MustWithEmptyValue(emptyValue OptionEmptyValuePolicy) OptionPolicy {
	next, err := policy.WithEmptyValue(emptyValue)
	if err != nil {
		panic(err)
	}

	return next
}

// WithAllowedSources returns a validated copy with allowed sources replaced.
func (policy OptionPolicy) WithAllowedSources(sources ...OptionSource) (OptionPolicy, error) {
	spec := policy.spec()
	spec.AllowedSources = cloneOptionSources(sources)

	return NewOptionPolicy(spec)
}

// MustWithAllowedSources returns a validated copy with allowed sources replaced
// and panics on invalid input.
func (policy OptionPolicy) MustWithAllowedSources(sources ...OptionSource) OptionPolicy {
	next, err := policy.WithAllowedSources(sources...)
	if err != nil {
		panic(err)
	}

	return next
}
