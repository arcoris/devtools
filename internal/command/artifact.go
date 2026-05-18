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

import "errors"

const (
	// maxArtifactIDLength is the maximum byte length of a stable artifact ID.
	maxArtifactIDLength = 255

	// maxArtifactKindLength is the maximum byte length of an artifact kind.
	maxArtifactKindLength = 63

	// maxArtifactFormatLength is the maximum byte length of an artifact format.
	maxArtifactFormatLength = 63

	// maxArtifactLocationLength is the maximum byte length of an artifact
	// location.
	maxArtifactLocationLength = 8192

	// maxArtifactMediaTypeLength is the maximum byte length of an artifact media
	// type.
	maxArtifactMediaTypeLength = 255

	// maxArtifactDescriptionLength is the maximum byte length of a compact
	// artifact description.
	maxArtifactDescriptionLength = 4096

	// maxArtifactLabelLength is the maximum byte length of one artifact label.
	maxArtifactLabelLength = 255
)

var (
	// ErrEmptyArtifactID reports that an artifact ID was not provided.
	ErrEmptyArtifactID = errors.New("command artifact id is empty")

	// ErrInvalidArtifactID reports that an artifact ID violates the artifact ID
	// grammar.
	ErrInvalidArtifactID = errors.New("command artifact id is invalid")

	// ErrEmptyArtifactKind reports that an artifact kind was not provided.
	ErrEmptyArtifactKind = errors.New("command artifact kind is empty")

	// ErrInvalidArtifactKind reports that an artifact kind violates the artifact
	// kind grammar.
	ErrInvalidArtifactKind = errors.New("command artifact kind is invalid")

	// ErrInvalidArtifactFormat reports that an artifact format violates the
	// artifact format grammar.
	ErrInvalidArtifactFormat = errors.New("command artifact format is invalid")

	// ErrEmptyArtifactLocation reports that an artifact location was not
	// provided.
	ErrEmptyArtifactLocation = errors.New("command artifact location is empty")

	// ErrInvalidArtifactLocation reports that an artifact location violates the
	// artifact location grammar.
	ErrInvalidArtifactLocation = errors.New("command artifact location is invalid")

	// ErrInvalidArtifact reports that an artifact declaration is malformed.
	ErrInvalidArtifact = errors.New("command artifact is invalid")

	// ErrInvalidArtifactDigest reports that an artifact digest is malformed.
	ErrInvalidArtifactDigest = errors.New("command artifact digest is invalid")
)

// ArtifactSpec describes a command-produced artifact before validation.
//
// ArtifactSpec is a construction DTO. Artifact stores detached copies of
// mutable input state, so callers cannot mutate constructed artifacts through
// shared slices.
//
// Artifact is an adapter-neutral reference to an output created, consumed, or
// described by a command. It does not create files, write reports, choose output
// directories, manage retention, upload data, or remove artifacts. Those
// responsibilities belong to execution, storage, and reporting layers.
type ArtifactSpec struct {
	// ID is a stable machine-facing artifact identifier.
	//
	// ID is intended for reports, tests, manifests, generated indexes, and
	// cross-references. It is not a display title.
	ID string

	// Kind describes the artifact category.
	//
	// Kind is open but validated. Common values include "report", "profile",
	// "trace", "coverage", "benchmark", "log", and "metadata".
	Kind ArtifactKind

	// Location is a repository-relative path, absolute path, or URI-like
	// location where the artifact can be found.
	//
	// The command kernel validates this as compact text only. It does not check
	// filesystem existence, URI reachability, permissions, or retention policy.
	Location string

	// Format is an optional machine-facing format key.
	//
	// Common values include "text", "json", "markdown", "html", "csv",
	// "pprof", "trace", and "binary".
	Format ArtifactFormat

	// MediaType is an optional media type such as "text/plain",
	// "application/json", or "text/markdown".
	MediaType string

	// Description is an optional compact human-facing description.
	Description string

	// SizeBytes is an optional non-negative artifact size.
	//
	// Nil means unknown. Zero means a known empty artifact.
	SizeBytes *int64

	// Digest contains optional content digest metadata.
	Digest *ArtifactDigestSpec

	// Labels contains optional machine-facing labels.
	//
	// Labels use the same compact key grammar as artifact IDs.
	Labels []string

	// Metadata contains optional machine-facing metadata.
	Metadata Metadata

	// Visibility controls default exposure in help, reports, docs, and
	// discovery.
	//
	// A zero visibility defaults to public.
	Visibility Visibility
}

// Artifact is a validated framework-neutral command artifact reference.
//
// Artifact is immutable-style:
//
//   - constructors normalize and copy mutable input values;
//   - accessors return detached copies;
//   - With* methods return validated copies;
//   - callers cannot mutate internal state through returned values.
type Artifact struct {
	id          ArtifactID
	kind        ArtifactKind
	location    string
	format      ArtifactFormat
	mediaType   string
	description string
	sizeBytes   int64
	hasSize     bool
	digest      *ArtifactDigest
	labels      []string
	metadata    Metadata
	visibility  Visibility
}
