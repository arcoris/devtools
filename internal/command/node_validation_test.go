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

func TestRootNodeRejectsIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec NodeSpec
	}{
		{
			name: "id",
			spec: NodeSpec{Kind: NodeRoot, ID: MustID("root"), Path: RootPath()},
		},
		{
			name: "non-root path",
			spec: NodeSpec{Kind: NodeRoot, Path: MustPath("bench")},
		},
		{
			name: "use",
			spec: NodeSpec{Kind: NodeRoot, Path: RootPath(), Use: "root"},
		},
		{
			name: "alias",
			spec: NodeSpec{Kind: NodeRoot, Path: RootPath(), Aliases: []string{"r"}},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewNode(tt.spec)
			if err == nil {
				t.Fatalf("NewNode() returned nil error")
			}

			if !errors.Is(err, ErrInvalidNode) {
				t.Fatalf("NewNode() error = %v, want ErrInvalidNode", err)
			}
		})
	}
}

func TestNonRootNodeRejectsInvalidIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		spec    NodeSpec
		wantErr error
	}{
		{
			name: "zero id",
			spec: NodeSpec{Kind: NodeCommand, Path: MustPath("check"), Use: "check"},
		},
		{
			name:    "invalid id",
			spec:    NodeSpec{Kind: NodeCommand, ID: ID("Check"), Path: MustPath("check"), Use: "check"},
			wantErr: ErrInvalidID,
		},
		{
			name: "root path",
			spec: NodeSpec{Kind: NodeCommand, ID: MustID("check"), Path: RootPath(), Use: "check"},
		},
		{
			name:    "invalid path",
			spec:    NodeSpec{Kind: NodeCommand, ID: MustID("check"), Path: Path{segments: []string{"Check"}}, Use: "check"},
			wantErr: ErrInvalidPath,
		},
		{
			name: "empty use",
			spec: NodeSpec{Kind: NodeCommand, ID: MustID("check"), Path: MustPath("check")},
		},
		{
			name:    "invalid use",
			spec:    NodeSpec{Kind: NodeCommand, ID: MustID("check"), Path: MustPath("check"), Use: "Check"},
			wantErr: ErrInvalidCommandNameSegment,
		},
		{
			name: "use does not match path leaf",
			spec: NodeSpec{Kind: NodeCommand, ID: MustID("bench.run"), Path: MustPath("bench", "run"), Use: "execute"},
		},
		{
			name: "id does not match path",
			spec: NodeSpec{Kind: NodeCommand, ID: MustID("bench.execute"), Path: MustPath("bench", "run"), Use: "run"},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewNode(tt.spec)
			if err == nil {
				t.Fatalf("NewNode() returned nil error")
			}

			if !errors.Is(err, ErrInvalidNode) {
				t.Fatalf("NewNode() error = %v, want ErrInvalidNode", err)
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewNode() error = %v, want wrapped %v", err, tt.wantErr)
			}
		})
	}
}

func TestFamilyNodeRequiresChildren(t *testing.T) {
	t.Parallel()

	_, err := NewFamilyNode(MustID("bench"), MustPath("bench"), "bench")
	if err == nil {
		t.Fatalf("NewFamilyNode() returned nil error")
	}

	if !errors.Is(err, ErrInvalidNode) {
		t.Fatalf("NewFamilyNode() error = %v, want ErrInvalidNode", err)
	}
}

func TestCommandNodeRejectsChildren(t *testing.T) {
	t.Parallel()

	child := mustTestCommandNode(t, "bench.run.child")

	_, err := NewNode(NodeSpec{
		Kind:     NodeCommand,
		ID:       MustID("bench.run"),
		Path:     MustPath("bench", "run"),
		Use:      "run",
		Children: []Node{child},
	})
	if err == nil {
		t.Fatalf("NewNode() returned nil error")
	}

	if !errors.Is(err, ErrInvalidNode) {
		t.Fatalf("NewNode() error = %v, want ErrInvalidNode", err)
	}
}

func TestNodeRejectsDuplicateAliases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		use     string
		aliases []string
	}{
		{name: "alias duplicates use", use: "run", aliases: []string{"run"}},
		{name: "alias duplicates alias", use: "run", aliases: []string{"execute", "execute"}},
		{name: "invalid alias", use: "run", aliases: []string{"Execute"}},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewNode(NodeSpec{
				Kind:    NodeCommand,
				ID:      MustID("bench.run"),
				Path:    MustPath("bench", "run"),
				Use:     tt.use,
				Aliases: tt.aliases,
			})
			if err == nil {
				t.Fatalf("NewNode() returned nil error")
			}

			if !errors.Is(err, ErrInvalidNode) {
				t.Fatalf("NewNode() error = %v, want ErrInvalidNode", err)
			}
		})
	}
}

func TestNodeRejectsSiblingSegmentConflicts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  Node
		right Node
	}{
		{
			name:  "duplicate use",
			left:  mustTestCommandNode(t, "bench.run"),
			right: mustTestCommandNode(t, "bench.run"),
		},
		{
			name:  "left alias conflicts with right use",
			left:  mustTestCommandNode(t, "bench.run", "execute"),
			right: mustTestCommandNode(t, "bench.execute"),
		},
		{
			name:  "alias conflicts with alias",
			left:  mustTestCommandNode(t, "bench.run", "start"),
			right: mustTestCommandNode(t, "bench.execute", "start"),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewFamilyNode(MustID("bench"), MustPath("bench"), "bench", tt.left, tt.right)
			if err == nil {
				t.Fatalf("NewFamilyNode() returned nil error")
			}

			if !errors.Is(err, ErrInvalidNode) {
				t.Fatalf("NewFamilyNode() error = %v, want ErrInvalidNode", err)
			}
		})
	}
}

func TestNodeRejectsInvalidChildParent(t *testing.T) {
	t.Parallel()

	child := mustTestCommandNode(t, "profile.cpu")

	_, err := NewFamilyNode(MustID("bench"), MustPath("bench"), "bench", child)
	if err == nil {
		t.Fatalf("NewFamilyNode() returned nil error")
	}

	if !errors.Is(err, ErrInvalidNode) {
		t.Fatalf("NewFamilyNode() error = %v, want ErrInvalidNode", err)
	}
}

func TestNodeRejectsRootChild(t *testing.T) {
	t.Parallel()

	rootChild := MustRootNode()

	_, err := NewRootNode(rootChild)
	if err == nil {
		t.Fatalf("NewRootNode(rootChild) returned nil error")
	}

	if !errors.Is(err, ErrInvalidNode) {
		t.Fatalf("NewRootNode(rootChild) error = %v, want ErrInvalidNode", err)
	}
}

func TestNodeRejectsInvalidDomainFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		spec NodeSpec
	}{
		{
			name: "duplicate topics",
			spec: NodeSpec{
				Kind:   NodeCommand,
				ID:     MustID("bench.run"),
				Path:   MustPath("bench", "run"),
				Use:    "run",
				Topics: []Topic{MustTopic("benchmark.run"), MustTopic("benchmark.run")},
			},
		},
		{
			name: "family handler",
			spec: NodeSpec{
				Kind:     NodeFamily,
				ID:       MustID("bench"),
				Path:     MustPath("bench"),
				Use:      "bench",
				Handler:  runtimeTestOKHandler(),
				Children: []Node{mustTestCommandNode(t, "bench.run")},
			},
		},
		{
			name: "deprecation conflict",
			spec: NodeSpec{
				Kind:       NodeCommand,
				ID:         MustID("bench.run"),
				Path:       MustPath("bench", "run"),
				Use:        "run",
				Deprecated: "old message",
				Metadata: MustMetadata(MetadataSpec{
					Deprecation: &DeprecationSpec{Message: "new message"},
				}),
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewNode(test.spec)
			if err == nil {
				t.Fatalf("NewNode() returned nil error")
			}

			if !errors.Is(err, ErrInvalidNode) {
				t.Fatalf("NewNode() error = %v, want ErrInvalidNode", err)
			}
		})
	}
}

func TestNodeWrapsBindingValidationErrors(t *testing.T) {
	t.Parallel()

	option := MustOption(OptionSpec{Name: "format", Kind: OptionKindString})

	_, err := NewNode(NodeSpec{
		Kind: NodeCommand,
		ID:   MustID("check"),
		Path: MustPath("check"),
		Use:  "check",
		Binding: Binding{
			options: []Option{option, option},
		},
	})
	if err == nil {
		t.Fatalf("NewNode() returned nil error")
	}

	if !errors.Is(err, ErrInvalidNode) {
		t.Fatalf("NewNode() error = %v, want ErrInvalidNode", err)
	}

	if !errors.Is(err, ErrInvalidBinding) {
		t.Fatalf("NewNode() error = %v, want ErrInvalidBinding", err)
	}
}

func TestNodeRejectsLegacyStructuredDocumentationConflicts(t *testing.T) {
	t.Parallel()

	base := func() NodeSpec {
		return NodeSpec{
			Kind: NodeCommand,
			ID:   MustID("check"),
			Path: MustPath("check"),
			Use:  "check",
		}
	}

	tests := []struct {
		name   string
		mutate func(*NodeSpec)
	}{
		{
			name: "short summary",
			mutate: func(spec *NodeSpec) {
				spec.Short = "Run checks"
				spec.Documentation = MustSummaryDocumentation("Run verification")
			},
		},
		{
			name: "long description",
			mutate: func(spec *NodeSpec) {
				spec.Long = "Run checks."
				spec.Documentation = MustDocumentation(DocumentationSpec{
					Description: "Run verification.",
				})
			},
		},
		{
			name: "usage",
			mutate: func(spec *NodeSpec) {
				spec.Usage = MustSimpleUsage("check [flags]")
				spec.Documentation = MustDocumentation(DocumentationSpec{
					Usage: MustSimpleUsage("check <target>"),
				})
			},
		},
		{
			name: "example",
			mutate: func(spec *NodeSpec) {
				spec.Example = "arcoris-tool check"
				spec.Documentation = MustDocumentation(DocumentationSpec{
					Notes: []string{nodeExampleNotePrefix + "arcoris-tool verify"},
				})
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			spec := base()
			test.mutate(&spec)

			_, err := NewNode(spec)
			if err == nil {
				t.Fatalf("NewNode() returned nil error")
			}

			if !errors.Is(err, ErrInvalidNode) {
				t.Fatalf("NewNode() error = %v, want ErrInvalidNode", err)
			}
		})
	}
}
