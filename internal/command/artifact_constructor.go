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

// NewArtifact validates spec and returns an Artifact.
func NewArtifact(spec ArtifactSpec) (Artifact, error) {
	id, err := NewArtifactID(spec.ID)
	if err != nil {
		return Artifact{}, err
	}

	if err := spec.Kind.Validate(); err != nil {
		return Artifact{}, err
	}

	visibility := spec.Visibility.OrDefault()

	var sizeBytes int64
	var hasSize bool
	if spec.SizeBytes != nil {
		sizeBytes = *spec.SizeBytes
		hasSize = true
	}

	var digest *ArtifactDigest
	if spec.Digest != nil {
		value, err := NewArtifactDigest(*spec.Digest)
		if err != nil {
			return Artifact{}, fmt.Errorf("%w: %w", ErrInvalidArtifact, err)
		}

		digest = &value
	}

	artifact := Artifact{
		id:          id,
		kind:        spec.Kind,
		location:    normalizeArtifactText(spec.Location),
		format:      spec.Format,
		mediaType:   normalizeArtifactText(spec.MediaType),
		description: normalizeArtifactBlock(spec.Description),
		sizeBytes:   sizeBytes,
		hasSize:     hasSize,
		digest:      digest,
		labels:      cloneArtifactStrings(spec.Labels),
		metadata:    spec.Metadata,
		visibility:  visibility,
	}

	if err := artifact.Validate(); err != nil {
		return Artifact{}, err
	}

	return artifact, nil
}

// MustArtifact validates spec and returns an Artifact.
//
// MustArtifact panics on invalid input. It is intended for static command
// definitions and tests where invalid artifact declarations are programmer
// errors.
func MustArtifact(spec ArtifactSpec) Artifact {
	artifact, err := NewArtifact(spec)
	if err != nil {
		panic(err)
	}

	return artifact
}
