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

// NewOptionValue validates spec and returns an OptionValue.
func NewOptionValue(spec OptionValueSpec) (OptionValue, error) {
	name, err := NewOptionName(spec.Name)
	if err != nil {
		return OptionValue{}, fmt.Errorf("%w: invalid name: %w", ErrInvalidOptionValue, err)
	}

	value := OptionValue{
		name:   name,
		kind:   spec.Kind,
		source: spec.Source,
		values: cloneStringSlice(spec.Values),
	}

	if err := value.Validate(); err != nil {
		return OptionValue{}, err
	}

	return value, nil
}

// MustOptionValue validates spec and returns an OptionValue.
//
// MustOptionValue panics on invalid input. It is intended for tests and
// controlled static wiring.
func MustOptionValue(spec OptionValueSpec) OptionValue {
	value, err := NewOptionValue(spec)
	if err != nil {
		panic(err)
	}

	return value
}

// NewScalarOptionValue constructs a scalar resolved option value.
func NewScalarOptionValue(name string, kind OptionKind, source OptionSource, raw string) (OptionValue, error) {
	if !kind.IsScalar() {
		return OptionValue{}, fmt.Errorf("%w: scalar value requires scalar kind, got %q", ErrInvalidOptionValue, kind)
	}

	return NewOptionValue(OptionValueSpec{
		Name:   name,
		Kind:   kind,
		Source: source,
		Values: []string{raw},
	})
}

// MustScalarOptionValue constructs a scalar resolved option value and panics on
// invalid input.
func MustScalarOptionValue(name string, kind OptionKind, source OptionSource, raw string) OptionValue {
	value, err := NewScalarOptionValue(name, kind, source, raw)
	if err != nil {
		panic(err)
	}

	return value
}

// NewListOptionValue constructs a list-shaped resolved option value.
func NewListOptionValue(name string, kind OptionKind, source OptionSource, values ...string) (OptionValue, error) {
	if !kind.IsList() {
		return OptionValue{}, fmt.Errorf("%w: list value requires list kind, got %q", ErrInvalidOptionValue, kind)
	}

	return NewOptionValue(OptionValueSpec{
		Name:   name,
		Kind:   kind,
		Source: source,
		Values: values,
	})
}

// MustListOptionValue constructs a list-shaped resolved option value and panics
// on invalid input.
func MustListOptionValue(name string, kind OptionKind, source OptionSource, values ...string) OptionValue {
	value, err := NewListOptionValue(name, kind, source, values...)
	if err != nil {
		panic(err)
	}

	return value
}

// NewOptionValueFromOption validates values against an Option declaration and
// returns a resolved OptionValue.
//
// When source is OptionSourceDefault and values are omitted, the option's
// declaration default values are used. When source is OptionSourceDefault and
// values are provided explicitly, they MUST match the option's declared default
// values. This keeps provenance honest: a value marked as "default" must be the
// declaration default.
func NewOptionValueFromOption(option Option, source OptionSource, values ...string) (OptionValue, error) {
	if err := option.Validate(); err != nil {
		return OptionValue{}, fmt.Errorf("%w: invalid option declaration: %w", ErrInvalidOptionValue, err)
	}

	if err := validateOptionValueSource(option, source); err != nil {
		return OptionValue{}, err
	}

	resolvedValues, err := resolveOptionValueInputs(option, source, values)
	if err != nil {
		return OptionValue{}, err
	}

	if err := validateOptionValueAgainstOption(option, resolvedValues); err != nil {
		return OptionValue{}, err
	}

	return NewOptionValue(OptionValueSpec{
		Name:   option.Name().String(),
		Kind:   option.Kind(),
		Source: source,
		Values: resolvedValues,
	})
}

// NewDefaultOptionValue returns an OptionValue from option declaration defaults.
func NewDefaultOptionValue(option Option) (OptionValue, error) {
	return NewOptionValueFromOption(option, OptionSourceDefault)
}

// MustDefaultOptionValue returns an OptionValue from option declaration defaults
// and panics on invalid input.
func MustDefaultOptionValue(option Option) OptionValue {
	value, err := NewDefaultOptionValue(option)
	if err != nil {
		panic(err)
	}

	return value
}
