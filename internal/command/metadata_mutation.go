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

// WithOwner returns a validated copy with Owner replaced.
func (metadata Metadata) WithOwner(owner string) (Metadata, error) {
	spec := metadata.spec()
	spec.Owner = owner

	return NewMetadata(spec)
}

// MustWithOwner returns a validated copy with Owner replaced and panics on
// invalid input.
func (metadata Metadata) MustWithOwner(owner string) Metadata {
	next, err := metadata.WithOwner(owner)
	if err != nil {
		panic(err)
	}

	return next
}

// WithArea returns a validated copy with Area replaced.
func (metadata Metadata) WithArea(area string) (Metadata, error) {
	spec := metadata.spec()
	spec.Area = area

	return NewMetadata(spec)
}

// MustWithArea returns a validated copy with Area replaced and panics on invalid
// input.
func (metadata Metadata) MustWithArea(area string) Metadata {
	next, err := metadata.WithArea(area)
	if err != nil {
		panic(err)
	}

	return next
}

// WithSince returns a validated copy with Since replaced.
func (metadata Metadata) WithSince(since string) (Metadata, error) {
	spec := metadata.spec()
	spec.Since = since

	return NewMetadata(spec)
}

// MustWithSince returns a validated copy with Since replaced and panics on
// invalid input.
func (metadata Metadata) MustWithSince(since string) Metadata {
	next, err := metadata.WithSince(since)
	if err != nil {
		panic(err)
	}

	return next
}

// WithDeprecation returns a validated copy with deprecation metadata replaced.
func (metadata Metadata) WithDeprecation(deprecation DeprecationSpec) (Metadata, error) {
	spec := metadata.spec()
	spec.Deprecation = &deprecation

	return NewMetadata(spec)
}

// MustWithDeprecation returns a validated copy with deprecation metadata
// replaced and panics on invalid input.
func (metadata Metadata) MustWithDeprecation(deprecation DeprecationSpec) Metadata {
	next, err := metadata.WithDeprecation(deprecation)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutDeprecation returns a validated copy without deprecation metadata.
func (metadata Metadata) WithoutDeprecation() Metadata {
	spec := metadata.spec()
	spec.Deprecation = nil

	return MustMetadata(spec)
}
