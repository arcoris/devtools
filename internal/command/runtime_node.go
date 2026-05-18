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
	"fmt"
	"strings"
)

// RuntimeFromNodeSpec describes framework-neutral Runtime construction from a
// command Node.
type RuntimeFromNodeSpec struct {
	// Name optionally overrides the runtime name.
	//
	// Empty defaults to Node.ID().String().
	Name string

	// Node is the command node used as the runtime declaration source.
	Node Node

	// Clock provides runtime timestamps.
	Clock RuntimeClock

	// EventSink receives lifecycle events.
	EventSink RuntimeEventSink

	// Options controls runtime behavior.
	Options RuntimeOptions

	// Metadata optionally overrides node metadata.
	//
	// Zero metadata means Node.Metadata() is used.
	Metadata Metadata

	// Visibility optionally overrides node visibility.
	//
	// Zero visibility means Node.Visibility() is used.
	Visibility Visibility
}

// NewRuntimeFromNode creates a Runtime from a command Node declaration.
//
// The helper copies only framework-neutral kernel values: command ID, Binding,
// RuntimeHandler, Metadata, and Visibility. It does not create adapter commands,
// parse flags, render output, write artifacts, or terminate the process.
func NewRuntimeFromNode(spec RuntimeFromNodeSpec) (Runtime, error) {
	node := spec.Node
	if err := node.Validate(); err != nil {
		return Runtime{}, fmt.Errorf("%w: invalid node: %w", ErrInvalidRuntime, err)
	}

	if !node.IsCommand() {
		return Runtime{}, fmt.Errorf(
			"%w: %w: node %q must be a command node, got %q",
			ErrInvalidRuntime,
			ErrInvalidNode,
			node.Path(),
			node.Kind(),
		)
	}

	handler, ok := node.RuntimeHandler()
	if !ok {
		return Runtime{}, fmt.Errorf(
			"%w: node %q has no runtime handler",
			ErrInvalidRuntime,
			node.Path(),
		)
	}

	name := strings.TrimSpace(spec.Name)
	if name == "" {
		name = node.ID().String()
	}

	metadata := spec.Metadata
	if metadata.IsZero() {
		metadata = node.Metadata()
	}

	visibility := spec.Visibility
	if visibility.IsZero() {
		visibility = node.Visibility()
	}

	runtime, err := NewRuntime(RuntimeSpec{
		Name:       name,
		CommandID:  node.ID(),
		Binding:    node.Binding(),
		Handler:    handler,
		Clock:      spec.Clock,
		EventSink:  spec.EventSink,
		Options:    spec.Options,
		Metadata:   metadata,
		Visibility: visibility,
	})
	if err != nil {
		return Runtime{}, err
	}

	return runtime, nil
}

// MustRuntimeFromNode creates a Runtime from a command Node and panics on
// invalid input.
func MustRuntimeFromNode(spec RuntimeFromNodeSpec) Runtime {
	runtime, err := NewRuntimeFromNode(spec)
	if err != nil {
		panic(err)
	}

	return runtime
}
