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
	"testing"
)

func TestNewRuntimeFromNodeCreatesExecutableRuntime(t *testing.T) {
	t.Parallel()

	node := runtimeNodeTestCommandNode()
	collector := &RuntimeEventCollector{}

	runtime, err := NewRuntimeFromNode(RuntimeFromNodeSpec{
		Node:      node,
		Clock:     FixedRuntimeClock{Time: runtimeTestTime()},
		EventSink: collector,
	})
	if err != nil {
		t.Fatalf("NewRuntimeFromNode() returned unexpected error: %v", err)
	}

	if got, want := runtime.Name(), "bench.run"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}

	if id, ok := runtime.CommandID(); !ok || id != node.ID() {
		t.Fatalf("CommandID() = %q, %v; want %q", id, ok, node.ID())
	}

	if got, want := runtime.Binding().OptionCount(), node.Binding().OptionCount(); got != want {
		t.Fatalf("Binding().OptionCount() = %d, want %d", got, want)
	}

	if got, want := runtime.Metadata().Owner(), "devtools"; got != want {
		t.Fatalf("Metadata().Owner() = %q, want %q", got, want)
	}

	if got, want := runtime.Visibility(), VisibilityHidden; got != want {
		t.Fatalf("Visibility() = %q, want %q", got, want)
	}

	result, err := runtime.Execute(context.Background(), RuntimeExecutionSpec{
		OptionValues: []OptionValue{
			MustScalarOptionValue("format", OptionKindEnum, OptionSourceCommandLine, "json"),
		},
		PositionalValues: []string{"stable"},
	})
	if err != nil {
		t.Fatalf("Execute() returned unexpected error: %v", err)
	}

	if !result.IsOK() {
		t.Fatalf("result IsOK() = false, want true")
	}
}

func TestNewRuntimeFromNodeAllowsNameOverride(t *testing.T) {
	t.Parallel()

	runtime := MustRuntimeFromNode(RuntimeFromNodeSpec{
		Name: "bench-runtime",
		Node: runtimeNodeTestCommandNode(),
	})

	if got, want := runtime.Name(), "bench-runtime"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestNewRuntimeFromNodeRejectsNonCommandNodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		node Node
	}{
		{name: "root", node: MustRootNode()},
		{
			name: "family",
			node: MustFamilyNode(
				MustID("bench"),
				MustPath("bench"),
				"bench",
				runtimeNodeTestCommandNode(),
			),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewRuntimeFromNode(RuntimeFromNodeSpec{Node: test.node})
			if err == nil {
				t.Fatalf("NewRuntimeFromNode() returned nil error")
			}

			if !errors.Is(err, ErrInvalidRuntime) {
				t.Fatalf("NewRuntimeFromNode() error = %v, want ErrInvalidRuntime", err)
			}

			if !errors.Is(err, ErrInvalidNode) {
				t.Fatalf("NewRuntimeFromNode() error = %v, want ErrInvalidNode", err)
			}
		})
	}
}

func TestNewRuntimeFromNodeRejectsInvalidNode(t *testing.T) {
	t.Parallel()

	_, err := NewRuntimeFromNode(RuntimeFromNodeSpec{})
	if err == nil {
		t.Fatalf("NewRuntimeFromNode() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRuntime) {
		t.Fatalf("NewRuntimeFromNode() error = %v, want ErrInvalidRuntime", err)
	}

	if !errors.Is(err, ErrInvalidNode) {
		t.Fatalf("NewRuntimeFromNode() error = %v, want ErrInvalidNode", err)
	}
}

func TestNewRuntimeFromNodeRejectsCommandWithoutHandler(t *testing.T) {
	t.Parallel()

	_, err := NewRuntimeFromNode(RuntimeFromNodeSpec{
		Node: mustTestCommandNode(t, "bench.run"),
	})
	if err == nil {
		t.Fatalf("NewRuntimeFromNode() returned nil error")
	}

	if !errors.Is(err, ErrInvalidRuntime) {
		t.Fatalf("NewRuntimeFromNode() error = %v, want ErrInvalidRuntime", err)
	}
}

func TestNewRuntimeFromNodeAllowsMetadataAndVisibilityOverride(t *testing.T) {
	t.Parallel()

	override := MustMetadata(MetadataSpec{Owner: "override"})

	runtime := MustRuntimeFromNode(RuntimeFromNodeSpec{
		Node:       runtimeNodeTestCommandNode(),
		Metadata:   override,
		Visibility: VisibilityInternal,
	})

	if got, want := runtime.Metadata().Owner(), "override"; got != want {
		t.Fatalf("Metadata().Owner() = %q, want %q", got, want)
	}

	if got, want := runtime.Visibility(), VisibilityInternal; got != want {
		t.Fatalf("Visibility() = %q, want %q", got, want)
	}
}

func runtimeNodeTestCommandNode() Node {
	return MustNode(NodeSpec{
		Kind:       NodeCommand,
		ID:         MustID("bench.run"),
		Path:       MustPath("bench", "run"),
		Use:        "run",
		Binding:    runtimeTestBinding(),
		Handler:    runtimeTestOKHandler(),
		Metadata:   MustMetadata(MetadataSpec{Owner: "devtools"}),
		Visibility: VisibilityHidden,
	})
}
