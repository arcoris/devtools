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
	"context"
	"errors"
	"time"
)

var (
	// ErrInvalidContext reports that command execution context metadata is
	// malformed.
	ErrInvalidContext = errors.New("command context is invalid")
)

// ContextSpec describes command execution context before validation.
//
// ContextSpec is a construction DTO. Context stores detached copies of mutable
// data so callers cannot mutate a constructed value through shared slices or
// maps.
//
// This type is not a replacement for context.Context. It wraps a base
// context.Context with command-tree metadata needed by command lifecycle,
// action execution, logging, diagnostics, reports, and tests.
type ContextSpec struct {
	// Context is the base cancellation/deadline/value context.
	//
	// If nil, context.Background() is used.
	Context context.Context

	// Node is the command node associated with this execution context.
	//
	// The node MAY be root, family, or command depending on lifecycle stage.
	// Action execution will normally require NodeCommand, but generic command
	// context can also be useful for traversal, pre-run validation, help
	// generation, and diagnostics.
	Node Node

	// Invocation describes adapter-neutral invocation metadata.
	Invocation Invocation

	// StartedAt is the timestamp when command execution context was created.
	//
	// If zero, NewContext uses time.Now().UTC().
	StartedAt time.Time

	// Fields contains optional machine-facing context metadata.
	//
	// Field keys use dotted kebab-case. Values are compact UTF-8 text without
	// disallowed control characters.
	Fields map[string]string
}

// Context is an adapter-neutral command execution context.
//
// Context embeds command-tree metadata around a base context.Context without
// depending on Cobra, pflag, os.Args, terminal output, or process exit codes.
//
// Context implements the context.Context interface by delegating Deadline,
// Done, Err, and Value to its base context. This allows it to be passed to code
// that expects context.Context while still exposing command-specific metadata
// through dedicated methods.
//
// Context is immutable-style:
//
//   - constructors copy input slices and maps;
//   - accessors return detached copies;
//   - With* methods return validated copies;
//   - callers cannot mutate internal state through returned values.
type Context struct {
	// base stores the cancellation/deadline/value context.
	//
	// It is never nil for values created by NewContext.
	base context.Context

	// node stores the validated command-tree node associated with this context.
	node Node

	// invocation stores validated adapter-neutral invocation metadata.
	invocation Invocation

	// startedAt stores the context creation timestamp.
	startedAt time.Time

	// fields stores validated lifecycle/runtime metadata.
	//
	// The map is never returned directly to callers.
	fields map[string]string
}
