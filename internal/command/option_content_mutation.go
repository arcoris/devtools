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
func (option Option) WithDocumentation(documentation Documentation) (Option, error) {
	spec := option.spec()
	spec.Documentation = documentation

	return NewOption(spec)
}

// MustWithDocumentation returns a validated copy with Documentation replaced
// and panics on invalid input.
func (option Option) MustWithDocumentation(documentation Documentation) Option {
	next, err := option.WithDocumentation(documentation)
	if err != nil {
		panic(err)
	}

	return next
}

// WithMetadata returns a validated copy with Metadata replaced.
func (option Option) WithMetadata(metadata Metadata) (Option, error) {
	spec := option.spec()
	spec.Metadata = metadata

	return NewOption(spec)
}

// MustWithMetadata returns a validated copy with Metadata replaced and panics
// on invalid input.
func (option Option) MustWithMetadata(metadata Metadata) Option {
	next, err := option.WithMetadata(metadata)
	if err != nil {
		panic(err)
	}

	return next
}

// WithVisibility returns a validated copy with Visibility replaced.
func (option Option) WithVisibility(visibility Visibility) (Option, error) {
	spec := option.spec()
	spec.Visibility = visibility

	return NewOption(spec)
}

// MustWithVisibility returns a validated copy with Visibility replaced and
// panics on invalid input.
func (option Option) MustWithVisibility(visibility Visibility) Option {
	next, err := option.WithVisibility(visibility)
	if err != nil {
		panic(err)
	}

	return next
}
