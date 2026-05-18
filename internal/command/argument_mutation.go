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

// WithDocumentation returns a validated copy with Documentation replaced.
func (argument Argument) WithDocumentation(documentation Documentation) (Argument, error) {
	spec := argument.spec()
	spec.Documentation = documentation

	return NewArgument(spec)
}

// MustWithDocumentation returns a validated copy with Documentation replaced
// and panics on invalid input.
func (argument Argument) MustWithDocumentation(documentation Documentation) Argument {
	next, err := argument.WithDocumentation(documentation)
	if err != nil {
		panic(err)
	}

	return next
}

// WithMetadata returns a validated copy with Metadata replaced.
func (argument Argument) WithMetadata(metadata Metadata) (Argument, error) {
	spec := argument.spec()
	spec.Metadata = metadata

	return NewArgument(spec)
}

// MustWithMetadata returns a validated copy with Metadata replaced and panics
// on invalid input.
func (argument Argument) MustWithMetadata(metadata Metadata) Argument {
	next, err := argument.WithMetadata(metadata)
	if err != nil {
		panic(err)
	}

	return next
}

// WithVisibility returns a validated copy with Visibility replaced.
func (argument Argument) WithVisibility(visibility Visibility) (Argument, error) {
	spec := argument.spec()
	spec.Visibility = visibility

	return NewArgument(spec)
}

// MustWithVisibility returns a validated copy with Visibility replaced and
// panics on invalid input.
func (argument Argument) MustWithVisibility(visibility Visibility) Argument {
	next, err := argument.WithVisibility(visibility)
	if err != nil {
		panic(err)
	}

	return next
}

// WithDefaultValues returns a validated copy with declaration defaults replaced.
func (argument Argument) WithDefaultValues(values ...string) (Argument, error) {
	spec := argument.spec()
	spec.DefaultValues = cloneStringSlice(values)

	return NewArgument(spec)
}

// MustWithDefaultValues returns a validated copy with declaration defaults
// replaced and panics on invalid input.
func (argument Argument) MustWithDefaultValues(values ...string) Argument {
	next, err := argument.WithDefaultValues(values...)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutDefault returns a validated copy without declaration defaults.
func (argument Argument) WithoutDefault() Argument {
	spec := argument.spec()
	spec.DefaultValues = nil

	return MustArgument(spec)
}

// WithAllowedValues returns a validated copy with allowed values replaced.
func (argument Argument) WithAllowedValues(values ...string) (Argument, error) {
	spec := argument.spec()
	spec.AllowedValues = cloneStringSlice(values)

	return NewArgument(spec)
}

// MustWithAllowedValues returns a validated copy with allowed values replaced
// and panics on invalid input.
func (argument Argument) MustWithAllowedValues(values ...string) Argument {
	next, err := argument.WithAllowedValues(values...)
	if err != nil {
		panic(err)
	}

	return next
}

// spec returns a detached construction spec.
func (argument Argument) spec() ArgumentSpec {
	return ArgumentSpec{
		Name:          argument.name.String(),
		Kind:          argument.kind,
		Requirement:   argument.requirement,
		Cardinality:   argument.cardinality,
		EmptyValue:    argument.emptyValue,
		Metavar:       argument.metavar,
		DefaultValues: cloneStringSlice(argument.defaultValues),
		AllowedValues: cloneStringSlice(argument.allowedValues),
		Documentation: argument.documentation,
		Metadata:      argument.metadata,
		Visibility:    argument.visibility,
	}
}
