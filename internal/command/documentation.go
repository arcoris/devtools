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
	// maxDocumentationSummaryLength is the maximum byte length of a short
	// one-line documentation summary.
	//
	// The summary is intended for command listings, generated indexes, compact
	// help output, and command discovery surfaces.
	maxDocumentationSummaryLength = 240

	// maxDocumentationDescriptionLength is the maximum byte length of a long
	// documentation description.
	//
	// Long descriptions may be rendered in detailed help or generated
	// documentation. They are still intended to be command-local text, not full
	// standalone manuals.
	maxDocumentationDescriptionLength = 16 * 1024

	// maxDocumentationNoteLength is the maximum byte length of one documentation
	// note.
	maxDocumentationNoteLength = 4096

	// maxDocumentationReferenceKeyLength is the maximum byte length of one
	// documentation reference key.
	maxDocumentationReferenceKeyLength = 255

	// maxDocumentationReferenceKeyDepth is the maximum number of dot-separated
	// segments accepted in one reference key.
	maxDocumentationReferenceKeyDepth = 32

	// maxDocumentationReferenceLabelLength is the maximum byte length of one
	// documentation reference label.
	maxDocumentationReferenceLabelLength = 512

	// maxDocumentationReferenceTargetLength is the maximum byte length of one
	// documentation reference target.
	maxDocumentationReferenceTargetLength = 2048
)

var (
	// ErrInvalidDocumentation reports that command documentation is malformed.
	ErrInvalidDocumentation = errors.New("command documentation is invalid")

	// ErrInvalidDocumentationReference reports that a documentation reference is
	// malformed.
	ErrInvalidDocumentationReference = errors.New("command documentation reference is invalid")
)

// DocumentationSpec describes command documentation before validation.
//
// DocumentationSpec is a construction DTO. Documentation stores detached copies
// of all mutable state so callers cannot mutate constructed values through
// shared slices.
//
// Documentation is adapter-neutral. It can be rendered by CLI help, generated
// Markdown, JSON discovery, GitHub summaries, or internal documentation tools,
// but it does not depend on any rendering framework.
type DocumentationSpec struct {
	// Summary is an optional short one-line human-facing description.
	//
	// Summary should answer what the command does in one sentence or sentence
	// fragment. It is suitable for command listings and short help output.
	Summary string

	// Description is an optional longer human-facing description.
	//
	// Description may contain multiple lines. It should explain behavior,
	// important constraints, and operational context. Detailed examples should
	// live in example tests or generated documentation, not in this field.
	Description string

	// Usage is optional structured usage syntax.
	//
	// A zero Usage means usage was not declared at this documentation layer.
	Usage Usage

	// Notes are optional compact documentation notes.
	//
	// Notes are rendered after the main description by documentation or help
	// adapters. They should be short, concrete, and command-local.
	Notes []string

	// References are optional structured references to related commands,
	// documents, or external resources.
	References []DocumentationReference
}

// Documentation is validated framework-neutral documentation for a command-tree
// node.
//
// Documentation deliberately does not own command identity, command grouping,
// execution behavior, or metadata ownership. Those concerns belong to ID, Path,
// Group, Topic, Action, and Metadata respectively.
//
// Documentation is immutable-style:
//
//   - constructors normalize and copy input slices;
//   - accessors return detached copies;
//   - With* methods return validated copies;
//   - duplicate notes and duplicate reference keys are rejected.
type Documentation struct {
	// summary stores the canonical short one-line documentation summary.
	summary string

	// description stores canonical compact multi-line documentation text.
	description string

	// usage stores optional structured usage syntax.
	usage Usage

	// notes stores canonical compact documentation notes.
	notes []string

	// references stores validated documentation references in declaration order.
	references []DocumentationReference
}
