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
func (option Option) WithName(name string) (Option, error) {
	spec := option.spec()
	spec.Name = name

	return NewOption(spec)
}

// MustWithName returns a validated copy with Name replaced and panics on
// invalid input.
func (option Option) MustWithName(name string) Option {
	next, err := option.WithName(name)
	if err != nil {
		panic(err)
	}

	return next
}

// WithKind returns a validated copy with Kind replaced.
func (option Option) WithKind(kind OptionKind) (Option, error) {
	spec := option.spec()
	spec.Kind = kind

	return NewOption(spec)
}

// MustWithKind returns a validated copy with Kind replaced and panics on invalid
// input.
func (option Option) MustWithKind(kind OptionKind) Option {
	next, err := option.WithKind(kind)
	if err != nil {
		panic(err)
	}

	return next
}
