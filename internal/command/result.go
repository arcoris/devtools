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
	"errors"
	"time"
)

const (
	// maxResultMessageLength is the maximum byte length of a compact result
	// message.
	maxResultMessageLength = 4096

	// maxResultWarningMessageLength is the maximum byte length of one warning
	// message.
	maxResultWarningMessageLength = 4096

	// maxResultHintLength is the maximum byte length of one result hint.
	maxResultHintLength = 4096

	// maxResultFieldKeyLength is the maximum byte length of one result metadata
	// field key.
	maxResultFieldKeyLength = 255

	// maxResultFieldValueLength is the maximum byte length of one result metadata
	// field value.
	maxResultFieldValueLength = 4096

	// maxResultExitCode is the maximum conventional process exit code value.
	//
	// Result is adapter-neutral, but command-line adapters commonly need a
	// bounded process exit code. Keeping the value in the kernel result model
	// makes reporting deterministic without forcing action code to depend on os.
	maxResultExitCode = 255
)

var (
	// ErrInvalidResult reports that a command result is malformed.
	ErrInvalidResult = errors.New("command result is invalid")

	// ErrInvalidResultStatus reports that a result status is not supported.
	ErrInvalidResultStatus = errors.New("command result status is invalid")

	// ErrInvalidResultWarning reports that a result warning is malformed.
	ErrInvalidResultWarning = errors.New("command result warning is invalid")
)

// ResultStatus describes the lifecycle-level outcome of a command execution.
//
// ResultStatus is intentionally separate from ActionStatus:
//
//   - ActionStatus describes a local action handler outcome;
//   - ResultStatus describes the final command lifecycle outcome after binding,
//     policy checks, action execution, artifact collection, diagnostics, and
//     adapter-level result normalization.
//
// Failures should still be represented as errors while the command is running.
// ResultStatusFailed is useful once an execution pipeline has already converted
// a failure into a structured final result.
type ResultStatus string

const (
	// ResultStatusOK means the command completed successfully.
	ResultStatusOK ResultStatus = "ok"

	// ResultStatusSkipped means the command intentionally did not run or had no
	// work to perform.
	ResultStatusSkipped ResultStatus = "skipped"

	// ResultStatusFailed means the command completed with a classified failure.
	ResultStatusFailed ResultStatus = "failed"

	// ResultStatusCanceled means the command was canceled before normal
	// completion.
	ResultStatusCanceled ResultStatus = "canceled"
)

// ResultSpec describes a final command result before validation.
//
// ResultSpec is a construction DTO. Result stores detached copies of mutable
// input state, so callers cannot mutate constructed results through shared
// slices or maps.
//
// Result does not render output, write artifacts, create files, or terminate a
// process. It is an adapter-neutral value object suitable for tests, reports,
// JSON output, CI summaries, and CLI adapters.
type ResultSpec struct {
	// Status is the final lifecycle status.
	//
	// Zero defaults to ResultStatusOK.
	Status ResultStatus

	// Message is an optional compact human-facing result summary.
	Message string

	// StartedAt is an optional command execution start timestamp.
	//
	// Zero means unknown or not recorded.
	StartedAt time.Time

	// FinishedAt is an optional command execution finish timestamp.
	//
	// Zero means unknown or not recorded.
	FinishedAt time.Time

	// ExitCode is an optional adapter-facing process exit code.
	//
	// Nil means not explicitly set. Use Result.RecommendedExitCode when an
	// adapter needs a deterministic process exit code.
	ExitCode *int

	// Artifacts contains artifacts produced or referenced by the command.
	Artifacts []Artifact

	// Warnings contains non-fatal warnings collected during command execution.
	Warnings []ResultWarning

	// Fields contains optional machine-facing result metadata.
	//
	// Field keys use a compact dot-separated key grammar. Values are compact
	// UTF-8 text.
	Fields map[string]string

	// Metadata contains optional machine-facing lifecycle metadata.
	Metadata Metadata

	// Visibility controls default exposure in reports, docs, and discovery.
	//
	// A zero visibility defaults to public.
	Visibility Visibility
}

// Result is a validated framework-neutral command execution result.
//
// Result is immutable-style:
//
//   - constructors normalize default values and copy mutable input state;
//   - accessors return detached copies;
//   - With* methods return validated copies;
//   - callers cannot mutate internal state through returned values.
type Result struct {
	status     ResultStatus
	message    string
	startedAt  time.Time
	finishedAt time.Time
	exitCode   int
	hasExit    bool
	artifacts  []Artifact
	warnings   []ResultWarning
	fields     map[string]string
	metadata   Metadata
	visibility Visibility
}
