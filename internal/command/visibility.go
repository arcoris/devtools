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

var (
	// ErrEmptyVisibility reports that a command visibility value was not provided.
	ErrEmptyVisibility = errors.New("command visibility is empty")

	// ErrInvalidVisibility reports that a command visibility value is not one
	// of the supported visibility states.
	ErrInvalidVisibility = errors.New("command visibility is invalid")
)

// Visibility describes how a command-tree node should be exposed by adapters,
// generated help, generated documentation, and discovery surfaces.
//
// Visibility is intentionally a closed kernel-level enum. Unlike Group or
// Topic, visibility is not project taxonomy. It is command infrastructure
// policy, so arbitrary user-defined values would make adapter behavior
// ambiguous.
//
// Visibility is adapter-neutral. A Cobra adapter, documentation generator, or
// command-discovery endpoint may interpret these values, but this type does not
// depend on any concrete adapter.
//
// The supported states are:
//
//   - public: visible in default help, default docs, and command discovery;
//   - hidden: omitted from default help/discovery, but still a valid command
//     for users who know it explicitly;
//   - internal: reserved for internal, generated, test-only, or debug-only
//     surfaces. Normal user-facing adapters should not expose or invoke these
//     commands by default.
//
// The zero value is treated as unset. Use DefaultVisibility when constructing
// command nodes that should default to public exposure.
type Visibility string

const (
	// VisibilityPublic marks a command-tree node as part of the normal public
	// user-facing surface.
	VisibilityPublic Visibility = "public"

	// VisibilityHidden marks a command-tree node as intentionally omitted from
	// default help and discovery output while still remaining addressable by
	// explicit command path.
	VisibilityHidden Visibility = "hidden"

	// VisibilityInternal marks a command-tree node as internal infrastructure.
	//
	// Internal nodes are not part of the normal user-facing CLI contract.
	// Adapters should hide them from default help, default generated
	// documentation, default command discovery, and normal invocation paths
	// unless an explicit internal/debug mode is enabled.
	VisibilityInternal Visibility = "internal"
)

var knownVisibilities = []Visibility{
	VisibilityPublic,
	VisibilityHidden,
	VisibilityInternal,
}
