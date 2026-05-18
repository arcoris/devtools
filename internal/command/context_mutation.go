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
	"time"
)

// WithBase returns a validated copy using base as the underlying context.
//
// If base is nil, context.Background() is used.
func (commandContext Context) WithBase(base context.Context) (Context, error) {
	spec := commandContext.spec()
	spec.Context = base

	return NewContext(spec)
}

// MustWithBase returns a validated copy using base as the underlying context
// and panics on invalid input.
func (commandContext Context) MustWithBase(base context.Context) Context {
	next, err := commandContext.WithBase(base)
	if err != nil {
		panic(err)
	}

	return next
}

// WithNode returns a validated copy associated with node.
func (commandContext Context) WithNode(node Node) (Context, error) {
	spec := commandContext.spec()
	spec.Node = node

	return NewContext(spec)
}

// MustWithNode returns a validated copy associated with node and panics on
// invalid input.
func (commandContext Context) MustWithNode(node Node) Context {
	next, err := commandContext.WithNode(node)
	if err != nil {
		panic(err)
	}

	return next
}

// WithInvocation returns a validated copy with invocation metadata replaced.
func (commandContext Context) WithInvocation(invocation Invocation) (Context, error) {
	spec := commandContext.spec()
	spec.Invocation = invocation

	return NewContext(spec)
}

// MustWithInvocation returns a validated copy with invocation metadata replaced
// and panics on invalid input.
func (commandContext Context) MustWithInvocation(invocation Invocation) Context {
	next, err := commandContext.WithInvocation(invocation)
	if err != nil {
		panic(err)
	}

	return next
}

// WithStartedAt returns a validated copy with StartedAt replaced.
func (commandContext Context) WithStartedAt(startedAt time.Time) (Context, error) {
	spec := commandContext.spec()
	spec.StartedAt = startedAt

	return NewContext(spec)
}

// MustWithStartedAt returns a validated copy with StartedAt replaced and panics
// on invalid input.
func (commandContext Context) MustWithStartedAt(startedAt time.Time) Context {
	next, err := commandContext.WithStartedAt(startedAt)
	if err != nil {
		panic(err)
	}

	return next
}

// WithFields returns a validated copy with context metadata fields replaced.
func (commandContext Context) WithFields(fields map[string]string) (Context, error) {
	spec := commandContext.spec()
	spec.Fields = cloneStringMap(fields)

	return NewContext(spec)
}

// MustWithFields returns a validated copy with context metadata fields replaced
// and panics on invalid input.
func (commandContext Context) MustWithFields(fields map[string]string) Context {
	next, err := commandContext.WithFields(fields)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutFields returns a validated copy without context metadata fields.
func (commandContext Context) WithoutFields() Context {
	spec := commandContext.spec()
	spec.Fields = nil

	return MustContext(spec)
}

// WithField returns a validated copy with one context metadata field added or
// replaced.
func (commandContext Context) WithField(key string, value string) (Context, error) {
	spec := commandContext.spec()
	if spec.Fields == nil {
		spec.Fields = make(map[string]string)
	}

	spec.Fields[key] = value

	return NewContext(spec)
}

// MustWithField returns a validated copy with one metadata field added or
// replaced and panics on invalid input.
func (commandContext Context) MustWithField(key string, value string) Context {
	next, err := commandContext.WithField(key, value)
	if err != nil {
		panic(err)
	}

	return next
}

// WithoutField returns a validated copy without one metadata field.
func (commandContext Context) WithoutField(key string) Context {
	spec := commandContext.spec()
	delete(spec.Fields, key)

	return MustContext(spec)
}
