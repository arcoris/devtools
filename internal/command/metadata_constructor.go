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

import "fmt"

// EmptyMetadata returns a valid empty metadata value.
func EmptyMetadata() Metadata {
	return Metadata{}
}

// NewMetadata validates spec and returns Metadata.
//
// NewMetadata copies all mutable input state. In particular, Annotations is
// cloned and Deprecation is converted into a validated value object before
// being stored.
func NewMetadata(spec MetadataSpec) (Metadata, error) {
	var deprecation *Deprecation

	if spec.Deprecation != nil {
		value, err := NewDeprecation(*spec.Deprecation)
		if err != nil {
			return Metadata{}, fmt.Errorf("%w: invalid deprecation: %w", ErrInvalidMetadata, err)
		}

		deprecation = &value
	}

	metadata := Metadata{
		owner:       spec.Owner,
		area:        spec.Area,
		since:       spec.Since,
		deprecation: deprecation,
		annotations: cloneStringMap(spec.Annotations),
	}

	if err := metadata.Validate(); err != nil {
		return Metadata{}, err
	}

	return metadata, nil
}

// MustMetadata validates spec and returns Metadata.
//
// MustMetadata panics on invalid input. It is intended for static command
// definitions and tests where invalid metadata is a programmer error.
func MustMetadata(spec MetadataSpec) Metadata {
	metadata, err := NewMetadata(spec)
	if err != nil {
		panic(err)
	}

	return metadata
}
