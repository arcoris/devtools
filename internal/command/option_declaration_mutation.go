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

// WithMetavar returns a validated copy with Metavar replaced.
func (option Option) WithMetavar(metavar string) (Option, error) {
	spec := option.spec()
	spec.Metavar = metavar

	return NewOption(spec)
}

// MustWithMetavar returns a validated copy with Metavar replaced and panics on
// invalid input.
func (option Option) MustWithMetavar(metavar string) Option {
	next, err := option.WithMetavar(metavar)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutMetavar returns a validated copy using Kind.ValueMetavar().
func (option Option) WithoutMetavar() Option {
	spec := option.spec()
	spec.Metavar = ""

	return MustOption(spec)
}

// WithDefaultValues returns a validated copy with declaration defaults replaced.
func (option Option) WithDefaultValues(values ...string) (Option, error) {
	spec := option.spec()
	spec.DefaultValues = cloneStringSlice(values)

	return NewOption(spec)
}

// MustWithDefaultValues returns a validated copy with declaration defaults
// replaced and panics on invalid input.
func (option Option) MustWithDefaultValues(values ...string) Option {
	next, err := option.WithDefaultValues(values...)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutDefault returns a validated copy without declaration default values.
func (option Option) WithoutDefault() Option {
	spec := option.spec()
	spec.DefaultValues = nil

	return MustOption(spec)
}

// WithAllowedValues returns a validated copy with allowed values replaced.
func (option Option) WithAllowedValues(values ...string) (Option, error) {
	spec := option.spec()
	spec.AllowedValues = cloneStringSlice(values)

	return NewOption(spec)
}

// MustWithAllowedValues returns a validated copy with allowed values replaced
// and panics on invalid input.
func (option Option) MustWithAllowedValues(values ...string) Option {
	next, err := option.WithAllowedValues(values...)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutAllowedValues returns a validated copy without allowed values.
func (option Option) WithoutAllowedValues() Option {
	spec := option.spec()
	spec.AllowedValues = nil

	return MustOption(spec)
}
