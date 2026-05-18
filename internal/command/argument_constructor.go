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

import "strings"

// NewArgument validates spec and returns an Argument.
func NewArgument(spec ArgumentSpec) (Argument, error) {
	name, err := NewArgumentName(spec.Name)
	if err != nil {
		return Argument{}, err
	}

	requirement := spec.Requirement.OrDefault()
	cardinality := spec.Cardinality.OrDefault()
	emptyValue := spec.EmptyValue.OrDefault()
	visibility := spec.Visibility.OrDefault()

	metavar := strings.TrimSpace(spec.Metavar)
	if metavar == "" {
		metavar = defaultArgumentMetavar(name)
	}

	argument := Argument{
		name:          name,
		kind:          spec.Kind,
		requirement:   requirement,
		cardinality:   cardinality,
		emptyValue:    emptyValue,
		metavar:       metavar,
		defaultValues: cloneStringSlice(spec.DefaultValues),
		allowedValues: normalizeArgumentValues(spec.AllowedValues),
		documentation: spec.Documentation,
		metadata:      spec.Metadata,
		visibility:    visibility,
	}

	if err := argument.Validate(); err != nil {
		return Argument{}, err
	}

	return argument, nil
}

// MustArgument validates spec and returns an Argument.
//
// MustArgument panics on invalid input. It is intended for static command
// definitions and tests where invalid argument declarations are programmer
// errors.
func MustArgument(spec ArgumentSpec) Argument {
	argument, err := NewArgument(spec)
	if err != nil {
		panic(err)
	}

	return argument
}
