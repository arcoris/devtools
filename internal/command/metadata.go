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
	// maxMetadataKeyLength is the maximum allowed byte length of a metadata key.
	maxMetadataKeyLength = 255

	// maxMetadataKeyDepth is the maximum number of metadata key segments.
	maxMetadataKeyDepth = 32

	// maxMetadataTextLength is the maximum allowed byte length for compact
	// metadata text fields such as Since.
	maxMetadataTextLength = 255

	// maxDeprecationMessageLength is the maximum allowed byte length of a
	// deprecation message.
	maxDeprecationMessageLength = 2048

	// maxAnnotationValueLength is the maximum allowed byte length of one
	// annotation value.
	maxAnnotationValueLength = 4096
)

var (
	// ErrInvalidMetadata reports that command metadata is malformed.
	ErrInvalidMetadata = errors.New("command metadata is invalid")

	// ErrInvalidDeprecation reports that command deprecation metadata is malformed.
	ErrInvalidDeprecation = errors.New("command deprecation metadata is invalid")
)

// MetadataSpec describes command metadata before validation.
//
// MetadataSpec is a construction DTO. Metadata stores detached copies of all
// mutable state, especially Annotations, so callers cannot mutate constructed
// metadata through shared map state.
type MetadataSpec struct {
	// Owner is an optional stable machine-facing owner key.
	//
	// Owner is not a display name. It should identify a responsible owner such
	// as a team, package, or internal area. When set, it uses the metadata key
	// grammar.
	Owner string

	// Area is an optional stable machine-facing area key.
	//
	// Area is intended for ownership, documentation grouping, policy routing,
	// and reporting. When set, it uses the metadata key grammar.
	Area string

	// Since is an optional compact text value describing when the node was
	// introduced.
	//
	// The command kernel does not impose SemVer or date semantics here. Higher
	// policy layers may enforce project-specific formats.
	Since string

	// Deprecation describes optional deprecation metadata.
	//
	// A nil value means the node is not deprecated. A non-nil value must be a
	// valid DeprecationSpec and is copied into an immutable-style Deprecation
	// value during construction.
	Deprecation *DeprecationSpec

	// Annotations contains optional extensible machine-facing metadata.
	//
	// Annotation keys use the metadata key grammar. Annotation values are
	// arbitrary compact UTF-8 text without disallowed control characters.
	Annotations map[string]string
}

// Metadata is framework-neutral metadata attached to a command-tree node.
//
// Metadata is intentionally not a map-only bag. It has a small set of promoted
// fields for common command-model concerns and a controlled annotation map for
// project-specific extensions.
//
// Metadata is immutable-style:
//
//   - constructors copy input maps;
//   - accessors return detached copies;
//   - With* methods return validated copies;
//   - callers cannot mutate internal map state.
//
// Metadata does not define command taxonomy. Group, Topic, Visibility, and
// future Maturity values are separate concepts. This type only stores generic
// node metadata such as owner, area, introduction marker, deprecation
// information, and annotations.
type Metadata struct {
	// owner stores the optional validated owner key.
	owner string

	// area stores the optional validated area key.
	area string

	// since stores the optional validated introduction marker.
	since string

	// deprecation stores optional validated deprecation metadata.
	//
	// A nil pointer is the canonical "not deprecated" state. A pointer is used
	// instead of an embedded zero value because an empty Deprecation is not a
	// valid deprecation declaration: deprecations require a message.
	deprecation *Deprecation

	// annotations stores optional validated extension metadata.
	//
	// The map is never returned directly to callers. Constructors and accessors
	// clone it so Metadata remains safe to pass by value.
	annotations map[string]string
}
