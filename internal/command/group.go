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
	// GroupSeparator separates hierarchical command group segments.
	//
	// Group uses the same separator style as ID because both are stable
	// machine-facing keys. A group is not a display title. Display titles belong
	// to documentation, help rendering, or adapter-specific presentation layers.
	GroupSeparator = "."

	// maxGroupLength is the maximum allowed byte length of a complete group key.
	//
	// Group keys are compact metadata values used for help grouping,
	// documentation generation, reporting, filtering, and policy lookup. They
	// are not free-form descriptions.
	maxGroupLength = 255

	// maxGroupDepth is the maximum number of hierarchical group segments.
	maxGroupDepth = 32

	// maxGroupSegmentLength is the maximum allowed byte length of one group
	// segment.
	maxGroupSegmentLength = 63
)

var (
	// ErrEmptyGroup reports that a command group key was not provided.
	ErrEmptyGroup = errors.New("command group is empty")

	// ErrInvalidGroup reports that a command group key violates the group grammar.
	ErrInvalidGroup = errors.New("command group is invalid")
)

// Group is a stable machine-facing classification key for command-tree nodes.
//
// Group is intentionally not an enum in this package. The command kernel only
// defines the value type and validation rules. Concrete group constants belong
// in command-definition packages or higher policy layers.
//
// Group is not a human-facing title. For example, a group key such as
// "benchmark" may be rendered as "Benchmark Commands" by a help adapter, but
// that rendering is presentation logic and does not belong to this value type.
//
// The group grammar is strict:
//
//   - a group is one or more dot-separated segments;
//   - each segment uses the generic command-name segment grammar;
//   - separators must not appear at the beginning, at the end, or consecutively;
//   - a complete group must not exceed maxGroupLength bytes;
//   - a complete group must not exceed maxGroupDepth segments;
//   - one segment must not exceed maxGroupSegmentLength bytes.
//
// Valid shape examples:
//
//   - "quality"
//   - "benchmark"
//   - "diagnostics"
//   - "diagnostics.perf"
//   - "config.schema"
type Group string
