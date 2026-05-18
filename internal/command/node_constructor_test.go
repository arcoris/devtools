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
	"testing"
)

func TestNodeKindStringAndValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		kind  NodeKind
		text  string
		valid bool
	}{
		{kind: NodeRoot, text: "root", valid: true},
		{kind: NodeFamily, text: "family", valid: true},
		{kind: NodeCommand, text: "command", valid: true},
		{kind: NodeKind("unknown"), text: "unknown", valid: false},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.text, func(t *testing.T) {
			t.Parallel()

			if got := tt.kind.String(); got != tt.text {
				t.Fatalf("String() = %q, want %q", got, tt.text)
			}

			if got := tt.kind.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestNewRootNode(t *testing.T) {
	t.Parallel()

	node, err := NewRootNode()
	if err != nil {
		t.Fatalf("NewRootNode() returned unexpected error: %v", err)
	}

	if !node.IsRoot() {
		t.Fatalf("IsRoot() = false, want true")
	}

	if !node.ID().IsZero() {
		t.Fatalf("root ID = %q, want zero ID", node.ID())
	}

	if !node.Path().IsRoot() {
		t.Fatalf("root Path() = %q, want root", node.Path())
	}

	if node.Use() != "" {
		t.Fatalf("root Use() = %q, want empty", node.Use())
	}
}

func TestNewCommandNode(t *testing.T) {
	t.Parallel()

	node, err := NewCommandNode(MustID("bench.run"), MustPath("bench", "run"), "run")
	if err != nil {
		t.Fatalf("NewCommandNode() returned unexpected error: %v", err)
	}

	if !node.IsCommand() {
		t.Fatalf("IsCommand() = false, want true")
	}

	if node.IsFamily() {
		t.Fatalf("IsFamily() = true, want false")
	}

	if node.IsRoot() {
		t.Fatalf("IsRoot() = true, want false")
	}

	if got, want := node.ID(), MustID("bench.run"); got != want {
		t.Fatalf("ID() = %q, want %q", got, want)
	}

	if got, want := node.Path(), MustPath("bench", "run"); !got.Equal(want) {
		t.Fatalf("Path() = %q, want %q", got, want)
	}

	if got, want := node.Use(), "run"; got != want {
		t.Fatalf("Use() = %q, want %q", got, want)
	}
}

func TestNewFamilyNode(t *testing.T) {
	t.Parallel()

	child := mustTestCommandNode(t, "bench.run")

	node, err := NewFamilyNode(MustID("bench"), MustPath("bench"), "bench", child)
	if err != nil {
		t.Fatalf("NewFamilyNode() returned unexpected error: %v", err)
	}

	if !node.IsFamily() {
		t.Fatalf("IsFamily() = false, want true")
	}

	if got, want := node.ChildCount(), 1; got != want {
		t.Fatalf("ChildCount() = %d, want %d", got, want)
	}
}

func TestNewNodeStoresMetadata(t *testing.T) {
	t.Parallel()

	node := MustNode(NodeSpec{
		Kind:       NodeCommand,
		ID:         MustID("check"),
		Path:       MustPath("check"),
		Use:        "check",
		Short:      "Run checks",
		Long:       "Run all configured checks.",
		Example:    "arcoris-tool check",
		Hidden:     true,
		Deprecated: "use verify",
	})

	if got, want := node.Short(), "Run checks"; got != want {
		t.Fatalf("Short() = %q, want %q", got, want)
	}

	if got, want := node.Long(), "Run all configured checks."; got != want {
		t.Fatalf("Long() = %q, want %q", got, want)
	}

	if got, want := node.Example(), "arcoris-tool check"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}

	if !node.Hidden() {
		t.Fatalf("Hidden() = false, want true")
	}

	if got, want := node.Deprecated(), "use verify"; got != want {
		t.Fatalf("Deprecated() = %q, want %q", got, want)
	}

	if !node.IsDeprecated() {
		t.Fatalf("IsDeprecated() = false, want true")
	}

	if got, want := node.Documentation().Summary(), "Run checks"; got != want {
		t.Fatalf("Documentation().Summary() = %q, want %q", got, want)
	}

	if got, want := node.Visibility(), VisibilityHidden; got != want {
		t.Fatalf("Visibility() = %q, want %q", got, want)
	}

	if !node.Metadata().IsDeprecated() {
		t.Fatalf("Metadata().IsDeprecated() = false, want true")
	}
}

func TestNewNodeStoresStructuredDomainFields(t *testing.T) {
	t.Parallel()

	binding := runtimeTestBinding()
	handler := runtimeTestOKHandler()

	node := MustNode(NodeSpec{
		Kind: NodeCommand,
		ID:   MustID("bench.run"),
		Path: MustPath("bench", "run"),
		Use:  "run",
		Documentation: MustDocumentation(DocumentationSpec{
			Summary: "Run benchmarks",
			Usage:   MustUsage(UsageSpec{Syntax: "bench run [flags] <suite>"}),
		}),
		Metadata: MustMetadata(MetadataSpec{
			Owner: "devtools",
			Deprecation: &DeprecationSpec{
				Message: "use bench execute",
			},
		}),
		Visibility: VisibilityInternal,
		Group:      MustGroup("benchmark"),
		Topics:     []Topic{MustTopic("benchmark.run")},
		Binding:    binding,
		Handler:    handler,
	})

	if got, want := node.Short(), "Run benchmarks"; got != want {
		t.Fatalf("Short() = %q, want %q", got, want)
	}

	if usage, ok := node.Usage(); !ok || usage.String() != "bench run [flags] <suite>" {
		t.Fatalf("Usage() = %q, %v; want syntax", usage, ok)
	}

	if got, want := node.Deprecated(), "use bench execute"; got != want {
		t.Fatalf("Deprecated() = %q, want %q", got, want)
	}

	if got, want := node.Visibility(), VisibilityInternal; got != want {
		t.Fatalf("Visibility() = %q, want %q", got, want)
	}

	if group, ok := node.Group(); !ok || group != MustGroup("benchmark") {
		t.Fatalf("Group() = %q, %v; want benchmark", group, ok)
	}

	topics := node.Topics()
	topics[0] = MustTopic("changed")
	if got, want := node.Topics()[0], MustTopic("benchmark.run"); got != want {
		t.Fatalf("Topics() returned mutable state: got %q, want %q", got, want)
	}

	if !node.HasBinding() {
		t.Fatalf("HasBinding() = false, want true")
	}

	if got := node.Binding().OptionCount(); got != binding.OptionCount() {
		t.Fatalf("Binding().OptionCount() = %d, want %d", got, binding.OptionCount())
	}

	if _, ok := node.RuntimeHandler(); !ok {
		t.Fatalf("RuntimeHandler() ok = false, want true")
	}
}

func TestMustNodePanicsForInvalidNode(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustNode(NodeSpec{Kind: NodeKind("unknown")})
	})
}

func TestNewNodeRejectsUnknownKind(t *testing.T) {
	t.Parallel()

	_, err := NewNode(NodeSpec{Kind: NodeKind("unknown")})
	if err == nil {
		t.Fatalf("NewNode() returned nil error")
	}

	if !errors.Is(err, ErrInvalidNode) {
		t.Fatalf("NewNode() error = %v, want ErrInvalidNode", err)
	}
}
