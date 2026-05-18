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

// NewOption validates spec and returns an Option.
func NewOption(spec OptionSpec) (Option, error) {
	name, err := NewOptionName(spec.Name)
	if err != nil {
		return Option{}, err
	}

	aliases, err := newOptionNames(spec.Aliases)
	if err != nil {
		return Option{}, err
	}

	if err := spec.Kind.Validate(); err != nil {
		return Option{}, fmt.Errorf("%w: invalid kind: %w", ErrInvalidOption, err)
	}

	policy := spec.Policy
	if isZeroOptionPolicy(policy) {
		policy, err = NewOptionPolicyForKind(spec.Kind, OptionPolicySpec{})
		if err != nil {
			return Option{}, fmt.Errorf("%w: invalid default policy: %w", ErrInvalidOption, err)
		}
	}

	metavar := spec.Metavar
	if metavar == "" {
		metavar = spec.Kind.ValueMetavar()
	}

	option := Option{
		name:          name,
		aliases:       cloneOptionNames(aliases),
		shorthand:     spec.Shorthand,
		kind:          spec.Kind,
		metavar:       metavar,
		defaultValues: cloneStringSlice(spec.DefaultValues),
		allowedValues: cloneStringSlice(spec.AllowedValues),
		policy:        policy,
		documentation: spec.Documentation,
		metadata:      spec.Metadata,
		visibility:    spec.Visibility.OrDefault(),
	}

	if err := option.Validate(); err != nil {
		return Option{}, err
	}

	return option, nil
}

// MustOption validates spec and returns an Option.
//
// MustOption panics on invalid input. It is intended for static option
// declarations and tests where invalid options are programmer errors.
func MustOption(spec OptionSpec) Option {
	option, err := NewOption(spec)
	if err != nil {
		panic(err)
	}

	return option
}
