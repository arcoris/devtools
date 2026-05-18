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

// Name returns the canonical argument name.
func (argument Argument) Name() ArgumentName {
	return argument.name
}

// Kind returns the scalar value kind accepted by the argument.
func (argument Argument) Kind() OptionKind {
	return argument.kind
}

// Requirement returns the argument requirement policy.
func (argument Argument) Requirement() ArgumentRequirement {
	return argument.requirement
}

// Cardinality returns the argument cardinality policy.
func (argument Argument) Cardinality() ArgumentCardinality {
	return argument.cardinality
}

// EmptyValue returns the explicit empty-value policy.
func (argument Argument) EmptyValue() OptionEmptyValuePolicy {
	return argument.emptyValue
}

// Metavar returns the value placeholder used in help and documentation.
func (argument Argument) Metavar() string {
	return argument.metavar
}

// DefaultValues returns detached declaration default values.
func (argument Argument) DefaultValues() []string {
	return cloneStringSlice(argument.defaultValues)
}

// HasDefault reports whether declaration default values are present.
func (argument Argument) HasDefault() bool {
	return len(argument.defaultValues) > 0
}

// DefaultValue returns the first declaration default value and whether it is
// present.
//
// For variadic arguments, use DefaultValues for lossless access.
func (argument Argument) DefaultValue() (string, bool) {
	if len(argument.defaultValues) == 0 {
		return "", false
	}

	return argument.defaultValues[0], true
}

// AllowedValues returns detached allowed values.
func (argument Argument) AllowedValues() []string {
	return cloneStringSlice(argument.allowedValues)
}

// HasAllowedValues reports whether allowed values are declared.
func (argument Argument) HasAllowedValues() bool {
	return len(argument.allowedValues) > 0
}

// AllowsValue reports whether value is allowed by the declaration-level
// allowed-value set.
//
// If no allowed values are declared, AllowsValue returns true. Type parsing and
// range validation are separate concerns.
func (argument Argument) AllowsValue(value string) bool {
	if len(argument.allowedValues) == 0 {
		return true
	}

	for _, allowed := range argument.allowedValues {
		if allowed == value {
			return true
		}
	}

	return false
}

// Documentation returns argument documentation.
func (argument Argument) Documentation() Documentation {
	return argument.documentation
}

// Metadata returns argument metadata.
func (argument Argument) Metadata() Metadata {
	return argument.metadata
}

// Visibility returns argument visibility.
func (argument Argument) Visibility() Visibility {
	return argument.visibility
}

// IsRequired reports whether at least one value must be supplied.
func (argument Argument) IsRequired() bool {
	return argument.requirement.IsRequired()
}

// IsOptional reports whether the argument may be omitted.
func (argument Argument) IsOptional() bool {
	return argument.requirement.IsOptional()
}

// IsVariadic reports whether the argument consumes multiple trailing values.
func (argument Argument) IsVariadic() bool {
	return argument.cardinality.IsVariadic()
}

// IsSingle reports whether the argument consumes at most one value.
func (argument Argument) IsSingle() bool {
	return argument.cardinality.IsSingle()
}

// IsVisibleByDefault reports whether default help/docs/discovery should expose
// the argument.
func (argument Argument) IsVisibleByDefault() bool {
	return argument.visibility.IsDiscoverableByDefault()
}

// MinValues returns the minimum number of runtime values required by this
// declaration.
func (argument Argument) MinValues() int {
	if argument.requirement.IsRequired() {
		return 1
	}

	return 0
}

// MaxValues returns the maximum number of runtime values accepted by this
// declaration and whether the maximum is bounded.
func (argument Argument) MaxValues() (int, bool) {
	if argument.cardinality.IsVariadic() {
		return 0, false
	}

	return 1, true
}

// AcceptsCount reports whether count satisfies requirement and cardinality.
func (argument Argument) AcceptsCount(count int) bool {
	if count < argument.MinValues() {
		return false
	}

	if maxValue, bounded := argument.MaxValues(); bounded && count > maxValue {
		return false
	}

	return true
}
