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

// NewContext validates spec and returns Context.
//
// NewContext replaces a nil base context with context.Background, replaces a
// zero StartedAt with time.Now().UTC(), and copies all mutable input state.
func NewContext(spec ContextSpec) (Context, error) {
	base := spec.Context
	if base == nil {
		base = context.Background()
	}

	startedAt := spec.StartedAt
	if startedAt.IsZero() {
		startedAt = time.Now().UTC()
	}

	commandContext := Context{
		base:       base,
		node:       spec.Node,
		invocation: spec.Invocation.Clone(),
		startedAt:  startedAt,
		fields:     cloneStringMap(spec.Fields),
	}

	if err := commandContext.Validate(); err != nil {
		return Context{}, err
	}

	return commandContext, nil
}

// MustContext validates spec and returns Context.
//
// MustContext panics on invalid input. It is intended for tests and static
// command wiring where invalid context metadata is a programmer error.
func MustContext(spec ContextSpec) Context {
	commandContext, err := NewContext(spec)
	if err != nil {
		panic(err)
	}

	return commandContext
}

// BackgroundContext returns a valid command context backed by
// context.Background.
//
// The returned context is associated with node and an empty invocation.
func BackgroundContext(node Node) Context {
	return MustContext(ContextSpec{
		Context: context.Background(),
		Node:    node,
	})
}
