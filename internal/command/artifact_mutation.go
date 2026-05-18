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

// WithLocation returns a validated copy with Location replaced.
func (artifact Artifact) WithLocation(location string) (Artifact, error) {
	spec := artifact.spec()
	spec.Location = location

	return NewArtifact(spec)
}

// MustWithLocation returns a validated copy with Location replaced and panics on
// invalid input.
func (artifact Artifact) MustWithLocation(location string) Artifact {
	next, err := artifact.WithLocation(location)
	if err != nil {
		panic(err)
	}

	return next
}

// WithFormat returns a validated copy with Format replaced.
func (artifact Artifact) WithFormat(format ArtifactFormat) (Artifact, error) {
	spec := artifact.spec()
	spec.Format = format

	return NewArtifact(spec)
}

// MustWithFormat returns a validated copy with Format replaced and panics on
// invalid input.
func (artifact Artifact) MustWithFormat(format ArtifactFormat) Artifact {
	next, err := artifact.WithFormat(format)
	if err != nil {
		panic(err)
	}

	return next
}

// WithDigest returns a validated copy with Digest replaced.
func (artifact Artifact) WithDigest(digest ArtifactDigestSpec) (Artifact, error) {
	spec := artifact.spec()
	spec.Digest = &digest

	return NewArtifact(spec)
}

// MustWithDigest returns a validated copy with Digest replaced and panics on
// invalid input.
func (artifact Artifact) MustWithDigest(digest ArtifactDigestSpec) Artifact {
	next, err := artifact.WithDigest(digest)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutDigest returns a validated copy without digest metadata.
func (artifact Artifact) WithoutDigest() Artifact {
	spec := artifact.spec()
	spec.Digest = nil

	return MustArtifact(spec)
}

// WithLabel returns a validated copy with one label appended.
func (artifact Artifact) WithLabel(label string) (Artifact, error) {
	spec := artifact.spec()
	spec.Labels = append(spec.Labels, label)

	return NewArtifact(spec)
}

// MustWithLabel returns a validated copy with one label appended and panics on
// invalid input.
func (artifact Artifact) MustWithLabel(label string) Artifact {
	next, err := artifact.WithLabel(label)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutLabel returns a validated copy without label.
func (artifact Artifact) WithoutLabel(label string) Artifact {
	spec := artifact.spec()
	nextLabels := make([]string, 0, len(spec.Labels))

	for _, current := range spec.Labels {
		if current == label {
			continue
		}

		nextLabels = append(nextLabels, current)
	}

	spec.Labels = nextLabels

	return MustArtifact(spec)
}

// WithMetadata returns a validated copy with Metadata replaced.
func (artifact Artifact) WithMetadata(metadata Metadata) (Artifact, error) {
	spec := artifact.spec()
	spec.Metadata = metadata

	return NewArtifact(spec)
}

// MustWithMetadata returns a validated copy with Metadata replaced and panics on
// invalid input.
func (artifact Artifact) MustWithMetadata(metadata Metadata) Artifact {
	next, err := artifact.WithMetadata(metadata)
	if err != nil {
		panic(err)
	}

	return next
}

// WithVisibility returns a validated copy with Visibility replaced.
func (artifact Artifact) WithVisibility(visibility Visibility) (Artifact, error) {
	spec := artifact.spec()
	spec.Visibility = visibility

	return NewArtifact(spec)
}

// MustWithVisibility returns a validated copy with Visibility replaced and
// panics on invalid input.
func (artifact Artifact) MustWithVisibility(visibility Visibility) Artifact {
	next, err := artifact.WithVisibility(visibility)
	if err != nil {
		panic(err)
	}

	return next
}

// spec returns a detached construction spec.
func (artifact Artifact) spec() ArtifactSpec {
	var size *int64
	if artifact.hasSize {
		value := artifact.sizeBytes
		size = &value
	}

	var digest *ArtifactDigestSpec
	if artifact.digest != nil {
		spec := artifact.digest.spec()
		digest = &spec
	}

	return ArtifactSpec{
		ID:          artifact.id.String(),
		Kind:        artifact.kind,
		Location:    artifact.location,
		Format:      artifact.format,
		MediaType:   artifact.mediaType,
		Description: artifact.description,
		SizeBytes:   size,
		Digest:      digest,
		Labels:      cloneArtifactStrings(artifact.labels),
		Metadata:    artifact.metadata,
		Visibility:  artifact.visibility,
	}
}
