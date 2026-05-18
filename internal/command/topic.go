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
	// TopicSeparator separates hierarchical command topic segments.
	//
	// Topic is a stable machine-facing key used for finer command
	// classification than Group. It is not a display title.
	TopicSeparator = "."

	// maxTopicLength is the maximum allowed byte length of a complete topic key.
	//
	// Topic keys are compact metadata values used for filtering, generated docs,
	// reporting, discovery, and policy lookups.
	maxTopicLength = 255

	// maxTopicDepth is the maximum number of hierarchical topic segments.
	maxTopicDepth = 32

	// maxTopicSegmentLength is the maximum allowed byte length of one topic
	// segment.
	maxTopicSegmentLength = 63
)

var (
	// ErrEmptyTopic reports that a command topic key was not provided.
	ErrEmptyTopic = errors.New("command topic is empty")

	// ErrInvalidTopic reports that a command topic key violates the topic grammar.
	ErrInvalidTopic = errors.New("command topic is invalid")
)

// Topic is a stable machine-facing key for fine-grained command classification.
//
// Topic is intentionally not an enum in this package. The command kernel only
// defines the value type and validation rules. Concrete topic constants belong
// in command-definition packages or higher policy layers.
//
// Topic is more specific than Group:
//
//   - Group answers "where should this command be displayed?";
//   - Topic answers "which functional subject does this command belong to?".
//
// For example, a command may be displayed in the "diagnostics" group and carry
// finer topics such as "profiling", "tracing", or "perf". This package does
// not enforce that relationship; it only validates topic keys.
//
// The topic grammar is strict:
//
//   - a topic is one or more dot-separated segments;
//   - each segment uses the generic command-name segment grammar;
//   - separators must not appear at the beginning, at the end, or consecutively;
//   - a complete topic must not exceed maxTopicLength bytes;
//   - a complete topic must not exceed maxTopicDepth segments;
//   - one segment must not exceed maxTopicSegmentLength bytes.
//
// Valid shape examples:
//
//   - "profiling"
//   - "tracing"
//   - "perf"
//   - "diagnostics.perf"
//   - "benchmark.compare"
type Topic string
