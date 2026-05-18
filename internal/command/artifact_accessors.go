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

import (
	"sort"
	"strings"
)

// ID returns the stable artifact identifier.
func (artifact Artifact) ID() ArtifactID {
	return artifact.id
}

// Kind returns the artifact kind.
func (artifact Artifact) Kind() ArtifactKind {
	return artifact.kind
}

// Location returns the artifact location.
func (artifact Artifact) Location() string {
	return artifact.location
}

// Format returns the optional artifact format.
func (artifact Artifact) Format() ArtifactFormat {
	return artifact.format
}

// HasFormat reports whether Format is set.
func (artifact Artifact) HasFormat() bool {
	return !artifact.format.IsZero()
}

// MediaType returns the optional media type.
func (artifact Artifact) MediaType() string {
	return artifact.mediaType
}

// HasMediaType reports whether MediaType is set.
func (artifact Artifact) HasMediaType() bool {
	return artifact.mediaType != ""
}

// Description returns the optional human-facing artifact description.
func (artifact Artifact) Description() string {
	return artifact.description
}

// HasDescription reports whether Description is set.
func (artifact Artifact) HasDescription() bool {
	return artifact.description != ""
}

// SizeBytes returns the known artifact size and whether it is set.
func (artifact Artifact) SizeBytes() (int64, bool) {
	return artifact.sizeBytes, artifact.hasSize
}

// HasSize reports whether SizeBytes is set.
func (artifact Artifact) HasSize() bool {
	return artifact.hasSize
}

// Digest returns the artifact digest and whether it is set.
func (artifact Artifact) Digest() (ArtifactDigest, bool) {
	if artifact.digest == nil {
		return ArtifactDigest{}, false
	}

	return *artifact.digest, true
}

// HasDigest reports whether digest metadata is set.
func (artifact Artifact) HasDigest() bool {
	return artifact.digest != nil
}

// Labels returns detached artifact labels.
func (artifact Artifact) Labels() []string {
	return cloneArtifactStrings(artifact.labels)
}

// HasLabels reports whether labels are present.
func (artifact Artifact) HasLabels() bool {
	return len(artifact.labels) > 0
}

// HasLabel reports whether label is present.
func (artifact Artifact) HasLabel(label string) bool {
	for _, current := range artifact.labels {
		if current == label {
			return true
		}
	}

	return false
}

// SortedLabels returns labels in deterministic lexical order.
func (artifact Artifact) SortedLabels() []string {
	labels := artifact.Labels()
	sort.Strings(labels)

	return labels
}

// Metadata returns artifact metadata.
func (artifact Artifact) Metadata() Metadata {
	return artifact.metadata
}

// Visibility returns artifact visibility.
func (artifact Artifact) Visibility() Visibility {
	return artifact.visibility
}

// IsVisibleByDefault reports whether default reports/docs/discovery should
// expose the artifact.
func (artifact Artifact) IsVisibleByDefault() bool {
	return artifact.visibility.IsDiscoverableByDefault()
}

// IsURI reports whether Location looks like a URI.
func (artifact Artifact) IsURI() bool {
	return strings.Contains(artifact.location, "://")
}

// IsRepositoryRelative reports whether Location is neither absolute nor URI-like.
func (artifact Artifact) IsRepositoryRelative() bool {
	return !artifact.IsURI() && !strings.HasPrefix(artifact.location, "/")
}
