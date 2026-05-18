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
	// maxInvocationArgumentLength is the maximum byte length of one parsed
	// invocation argument.
	maxInvocationArgumentLength = 4096

	// maxInvocationRawArgumentLength is the maximum byte length of one raw
	// invocation argument.
	maxInvocationRawArgumentLength = 4096

	// maxInvocationWorkingDirLength is the maximum byte length of an invocation
	// working directory.
	maxInvocationWorkingDirLength = 8192

	// maxInvocationEnvNameLength is the maximum byte length of one captured
	// environment variable name.
	maxInvocationEnvNameLength = 255

	// maxInvocationEnvValueLength is the maximum byte length of one captured
	// environment variable value.
	maxInvocationEnvValueLength = 8192

	// maxInvocationFieldKeyLength is the maximum byte length of one invocation
	// metadata field key.
	maxInvocationFieldKeyLength = 255

	// maxInvocationFieldKeyDepth is the maximum number of dot-separated
	// segments in one invocation metadata field key.
	maxInvocationFieldKeyDepth = 32

	// maxInvocationFieldValueLength is the maximum byte length of one invocation
	// metadata field value.
	maxInvocationFieldValueLength = 4096
)

var (
	// ErrInvalidInvocation reports that command invocation metadata is
	// malformed.
	ErrInvalidInvocation = errors.New("command invocation is invalid")
)

// InvocationSpec describes adapter-neutral command invocation metadata before
// validation.
//
// InvocationSpec is a construction DTO. Invocation stores detached copies of
// mutable input state, so callers cannot mutate a constructed value through
// shared slices or maps.
//
// Invocation does not parse command-line flags and does not define option
// semantics. It stores already-collected invocation metadata for diagnostics,
// lifecycle events, reports, audit records, reproducibility, and tests.
type InvocationSpec struct {
	// Arguments contains parsed positional arguments.
	//
	// These are adapter-neutral positional values after command resolution and
	// flag parsing have already happened.
	Arguments []string

	// RawArguments contains raw adapter arguments, if available.
	//
	// RawArguments may include command names, flags, or adapter-specific tokens.
	// They are intended for diagnostics and audit metadata only.
	RawArguments []string

	// WorkingDir is the working directory from which the command was invoked.
	//
	// The command kernel does not require the path to exist because Invocation
	// is metadata, not filesystem validation.
	WorkingDir string

	// Env contains selected environment metadata.
	//
	// Env is intentionally not required to contain the full process environment.
	// Callers SHOULD include only values that are useful for diagnostics,
	// reproducibility, or policy decisions.
	Env map[string]string

	// Fields contains adapter-neutral invocation metadata.
	//
	// Field keys use dotted kebab-case. Values are compact UTF-8 text without
	// disallowed control characters.
	Fields map[string]string
}

// Invocation describes adapter-neutral metadata for one command invocation.
//
// Invocation is separate from Context:
//
//   - Invocation describes how the command was invoked;
//   - Context combines base cancellation context, command node, invocation,
//     runtime metadata, and lifecycle metadata.
//
// Invocation is immutable-style:
//
//   - constructors copy input slices and maps;
//   - accessors return detached copies;
//   - With* methods return validated copies;
//   - callers cannot mutate internal state through returned values.
type Invocation struct {
	// arguments stores parsed positional arguments.
	arguments []string

	// rawArguments stores raw adapter arguments for diagnostics and audit logs.
	rawArguments []string

	// workingDir stores the invocation working directory.
	workingDir string

	// env stores selected captured environment metadata.
	env map[string]string

	// fields stores adapter-neutral invocation metadata.
	fields map[string]string
}
