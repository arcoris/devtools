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

// WithName returns a validated copy with Name replaced.
func (value OptionValue) WithName(name string) (OptionValue, error) {
	spec := value.Spec()
	spec.Name = name

	return NewOptionValue(spec)
}

// MustWithName returns a validated copy with Name replaced and panics on
// invalid input.
func (value OptionValue) MustWithName(name string) OptionValue {
	next, err := value.WithName(name)
	if err != nil {
		panic(err)
	}

	return next
}

// WithKind returns a validated copy with Kind replaced.
func (value OptionValue) WithKind(kind OptionKind) (OptionValue, error) {
	spec := value.Spec()
	spec.Kind = kind

	return NewOptionValue(spec)
}

// MustWithKind returns a validated copy with Kind replaced and panics on invalid
// input.
func (value OptionValue) MustWithKind(kind OptionKind) OptionValue {
	next, err := value.WithKind(kind)
	if err != nil {
		panic(err)
	}

	return next
}

// WithSource returns a validated copy with Source replaced.
func (value OptionValue) WithSource(source OptionSource) (OptionValue, error) {
	spec := value.Spec()
	spec.Source = source

	return NewOptionValue(spec)
}

// MustWithSource returns a validated copy with Source replaced and panics on
// invalid input.
func (value OptionValue) MustWithSource(source OptionSource) OptionValue {
	next, err := value.WithSource(source)
	if err != nil {
		panic(err)
	}

	return next
}

// WithValues returns a validated copy with raw values replaced.
func (value OptionValue) WithValues(values ...string) (OptionValue, error) {
	spec := value.Spec()
	spec.Values = cloneStringSlice(values)

	return NewOptionValue(spec)
}

// MustWithValues returns a validated copy with raw values replaced and panics on
// invalid input.
func (value OptionValue) MustWithValues(values ...string) OptionValue {
	next, err := value.WithValues(values...)
	if err != nil {
		panic(err)
	}

	return next
}
