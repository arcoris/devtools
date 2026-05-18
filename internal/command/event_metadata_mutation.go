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
func (event Event) WithMetadata(metadata Metadata) (Event, error) {
	spec := event.spec()
	spec.Metadata = metadata

	return NewEvent(spec)
}

// MustWithMetadata returns a validated copy with Metadata replaced and panics on
// invalid input.
func (event Event) MustWithMetadata(metadata Metadata) Event {
	next, err := event.WithMetadata(metadata)
	if err != nil {
		panic(err)
	}

	return next
}

// WithVisibility returns a validated copy with Visibility replaced.
func (event Event) WithVisibility(visibility Visibility) (Event, error) {
	spec := event.spec()
	spec.Visibility = visibility

	return NewEvent(spec)
}

// MustWithVisibility returns a validated copy with Visibility replaced and
// panics on invalid input.
func (event Event) MustWithVisibility(visibility Visibility) Event {
	next, err := event.WithVisibility(visibility)
	if err != nil {
		panic(err)
	}

	return next
}
