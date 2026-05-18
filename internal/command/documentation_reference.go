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

// DocumentationReferenceKind classifies a documentation reference target.
//
// This is intentionally a small closed set because renderers need predictable
// behavior for references. Project-specific classification can be added above
// this layer through Metadata or generated documentation logic.
type DocumentationReferenceKind string

const (
	// DocumentationReferenceCommandID references another command by stable ID.
	DocumentationReferenceCommandID DocumentationReferenceKind = "command-id"

	// DocumentationReferenceCommandPath references another command by path text.
	DocumentationReferenceCommandPath DocumentationReferenceKind = "command-path"

	// DocumentationReferenceDocument references a project-local document.
	DocumentationReferenceDocument DocumentationReferenceKind = "document"

	// DocumentationReferenceURL references an external URL-like target.
	DocumentationReferenceURL DocumentationReferenceKind = "url"
)

// DocumentationReferenceSpec describes a documentation reference before
// validation.
type DocumentationReferenceSpec struct {
	// Key is a stable machine-facing reference key.
	//
	// Key is used for replacement, de-duplication, deterministic indexing, and
	// generated anchors.
	Key string

	// Kind classifies the target.
	Kind DocumentationReferenceKind

	// Label is a human-facing label.
	Label string

	// Target is the machine-facing reference target.
	//
	// The command kernel validates it as compact single-line UTF-8 text. It does
	// not check whether a command, document, or URL exists. Cross-reference
	// resolution is a higher-level registry or documentation-generation concern.
	Target string
}

// DocumentationReference is a validated reference from command documentation to
// related material.
type DocumentationReference struct {
	// key stores the canonical stable machine-facing reference key.
	key string

	// kind stores the supported reference target class.
	kind DocumentationReferenceKind

	// label stores the canonical human-facing reference label.
	label string

	// target stores the canonical machine-facing reference target.
	target string
}
