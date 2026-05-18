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

// WithMetadata returns a validated copy with Metadata replaced.
func (result Result) WithMetadata(metadata Metadata) (Result, error) {
	spec := result.spec()
	spec.Metadata = metadata

	return NewResult(spec)
}

// MustWithMetadata returns a validated copy with Metadata replaced and panics on
// invalid input.
func (result Result) MustWithMetadata(metadata Metadata) Result {
	next, err := result.WithMetadata(metadata)
	if err != nil {
		panic(err)
	}

	return next
}

// WithVisibility returns a validated copy with Visibility replaced.
func (result Result) WithVisibility(visibility Visibility) (Result, error) {
	spec := result.spec()
	spec.Visibility = visibility

	return NewResult(spec)
}

// MustWithVisibility returns a validated copy with Visibility replaced and
// panics on invalid input.
func (result Result) MustWithVisibility(visibility Visibility) Result {
	next, err := result.WithVisibility(visibility)
	if err != nil {
		panic(err)
	}

	return next
}
