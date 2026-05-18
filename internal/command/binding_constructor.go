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

// EmptyBinding returns a valid binding with no options and no arguments.
func EmptyBinding() Binding {
	return Binding{}
}

// NewBinding validates spec and returns a Binding.
func NewBinding(spec BindingSpec) (Binding, error) {
	binding := Binding{
		options:   cloneBindingOptions(spec.Options),
		arguments: cloneBindingArguments(spec.Arguments),
	}

	if err := binding.Validate(); err != nil {
		return Binding{}, err
	}

	return binding, nil
}

// MustBinding validates spec and returns a Binding.
//
// MustBinding panics on invalid input. It is intended for static command
// definitions and tests where invalid bindings are programmer errors.
func MustBinding(spec BindingSpec) Binding {
	binding, err := NewBinding(spec)
	if err != nil {
		panic(err)
	}

	return binding
}
