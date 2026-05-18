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

// WithShorthand returns a validated copy with Shorthand replaced.
func (option Option) WithShorthand(shorthand string) (Option, error) {
	spec := option.spec()
	spec.Shorthand = shorthand

	return NewOption(spec)
}

// MustWithShorthand returns a validated copy with Shorthand replaced and panics
// on invalid input.
func (option Option) MustWithShorthand(shorthand string) Option {
	next, err := option.WithShorthand(shorthand)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutShorthand returns a validated copy without Shorthand.
func (option Option) WithoutShorthand() Option {
	spec := option.spec()
	spec.Shorthand = ""

	return MustOption(spec)
}
