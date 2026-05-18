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
	"sort"
	"time"
)

// Deadline delegates to the underlying context.Context.
func (commandContext Context) Deadline() (time.Time, bool) {
	return commandContext.base.Deadline()
}

// Done delegates to the underlying context.Context.
func (commandContext Context) Done() <-chan struct{} {
	return commandContext.base.Done()
}

// Err delegates to the underlying context.Context.
func (commandContext Context) Err() error {
	return commandContext.base.Err()
}

// Value delegates to the underlying context.Context.
func (commandContext Context) Value(key any) any {
	return commandContext.base.Value(key)
}

// Base returns the underlying context.Context.
//
// The returned context is not wrapped. Use this when calling APIs that should
// not depend on command-specific metadata.
func (commandContext Context) Base() context.Context {
	return commandContext.base
}

// Node returns the command-tree node associated with this context.
func (commandContext Context) Node() Node {
	return commandContext.node
}

// ID returns the node ID associated with this context.
func (commandContext Context) ID() ID {
	return commandContext.node.ID()
}

// Path returns the node path associated with this context.
func (commandContext Context) Path() Path {
	return commandContext.node.Path()
}

// Invocation returns a detached copy of invocation metadata.
func (commandContext Context) Invocation() Invocation {
	return commandContext.invocation.Clone()
}

// StartedAt returns the context creation timestamp.
func (commandContext Context) StartedAt() time.Time {
	return commandContext.startedAt
}

// Field returns one context metadata field and whether it exists.
func (commandContext Context) Field(key string) (string, bool) {
	value, ok := commandContext.fields[key]

	return value, ok
}

// HasField reports whether a context metadata field exists.
func (commandContext Context) HasField(key string) bool {
	_, ok := commandContext.Field(key)

	return ok
}

// HasFields reports whether any context metadata fields are set.
func (commandContext Context) HasFields() bool {
	return len(commandContext.fields) > 0
}

// FieldCount returns the number of context metadata fields.
func (commandContext Context) FieldCount() int {
	return len(commandContext.fields)
}

// Fields returns a detached copy of all context metadata fields.
func (commandContext Context) Fields() map[string]string {
	return cloneStringMap(commandContext.fields)
}

// FieldKeys returns context metadata keys in deterministic lexical order.
func (commandContext Context) FieldKeys() []string {
	keys := make([]string, 0, len(commandContext.fields))
	for key := range commandContext.fields {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

// IsZero reports whether the command context has not been initialized.
func (commandContext Context) IsZero() bool {
	return commandContext.base == nil && commandContext.node.Kind() == ""
}

// Spec returns a detached construction spec for context.
func (commandContext Context) Spec() ContextSpec {
	return commandContext.spec()
}

// spec returns a detached construction spec.
func (commandContext Context) spec() ContextSpec {
	return ContextSpec{
		Context:    commandContext.base,
		Node:       commandContext.node,
		Invocation: commandContext.invocation.Clone(),
		StartedAt:  commandContext.startedAt,
		Fields:     cloneStringMap(commandContext.fields),
	}
}
