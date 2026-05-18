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

// NodeSpec describes a command-tree node before validation.
//
// NodeSpec is intentionally a construction DTO. Node stores detached copies of
// mutable slices so callers cannot mutate a constructed value through shared
// slice backing arrays.
type NodeSpec struct {
	// Kind defines whether the node is root, family, or leaf command.
	Kind NodeKind

	// ID is the stable machine-facing node identifier.
	//
	// Root nodes must use the zero ID. Non-root nodes must use a valid non-zero
	// ID. For non-root nodes, ID must match Path.ID().
	ID ID

	// Path is the logical command-tree path.
	//
	// Root nodes must use the root path. Non-root nodes must use a non-root
	// path. For non-root nodes, Path.Leaf() must equal Use.
	Path Path

	// Use is the canonical adapter-facing segment for this node.
	//
	// Use is one segment, not a full command path. It uses the command-name
	// segment grammar from validate.go. Root nodes must leave Use empty.
	Use string

	// Aliases are alternative adapter-facing segments for the same node.
	//
	// Aliases are not stable identities. They are convenience spellings for
	// adapters and users. Aliases use the same segment grammar as Use.
	Aliases []string

	// Short is a short one-line human-facing description.
	//
	// Structural validation does not require Short. Registry-level invariants
	// can require documentation for public nodes later.
	Short string

	// Long is a longer human-facing description.
	Long string

	// Example contains human-facing usage examples.
	Example string

	// Hidden reports whether adapters should hide the node from default help.
	Hidden bool

	// Deprecated contains an optional deprecation message.
	//
	// Deprecated is a compatibility field. New declarations should set
	// Metadata.Deprecation instead.
	Deprecated string

	// Documentation contains structured human-facing documentation.
	//
	// Short, Long, Example, and Usage are merged into Documentation during
	// construction when their structured equivalents are absent.
	Documentation Documentation

	// Usage contains structured usage syntax.
	//
	// Usage is a convenience field for declarations. Node stores usage inside
	// Documentation so the same concept is not represented twice.
	Usage Usage

	// Metadata contains machine-facing command metadata.
	//
	// Deprecated is merged into Metadata.Deprecation during construction when
	// structured deprecation metadata is absent.
	Metadata Metadata

	// Visibility controls default help, docs, and discovery exposure.
	//
	// A zero visibility defaults to public unless Hidden is true.
	Visibility Visibility

	// Group optionally classifies the node for help, generated docs, reporting,
	// and policy lookup.
	Group Group

	// Topics optionally classify the node by functional subject.
	Topics []Topic

	// Binding declares framework-neutral command inputs.
	//
	// Binding does not parse command-line syntax or read config/environment
	// values. It binds already resolved OptionValue and positional values.
	Binding Binding

	// Handler optionally contains canonical runtime executable behavior for
	// command nodes.
	//
	// RuntimeHandler is the canonical execution contract. Framework adapters may
	// also keep handler wiring outside Node when they prefer a separate registry.
	Handler RuntimeHandler

	// Children contains child nodes.
	//
	// Root and family nodes may have children. Command nodes must not have
	// children.
	Children []Node
}

// Node is an immutable-style framework-neutral command-tree node.
//
// Node represents either the root of the command tree, a command family, or a
// leaf command. It is intentionally independent from Cobra or any other CLI
// framework.
//
// Node stores declaration-level command metadata and may optionally reference a
// canonical RuntimeHandler for leaf commands. It still does not own adapter
// objects, parser state, rendering, process termination, or artifact storage.
type Node struct {
	kind          NodeKind
	id            ID
	path          Path
	use           string
	aliases       []string
	documentation Documentation
	metadata      Metadata
	visibility    Visibility
	group         Group
	topics        []Topic
	binding       Binding
	handler       RuntimeHandler
	children      []Node
}

// Spec returns a detached construction spec for node.
func (node Node) Spec() NodeSpec {
	return node.spec()
}

// spec returns a detached construction spec for node.
func (node Node) spec() NodeSpec {
	return NodeSpec{
		Kind:          node.kind,
		ID:            node.id,
		Path:          node.path,
		Use:           node.use,
		Aliases:       cloneStringSlice(node.aliases),
		Short:         node.Short(),
		Long:          node.Long(),
		Example:       node.Example(),
		Hidden:        node.Hidden(),
		Deprecated:    node.Deprecated(),
		Documentation: node.documentation,
		Usage:         node.UsageOrZero(),
		Metadata:      node.metadata,
		Visibility:    node.visibility,
		Group:         node.group,
		Topics:        cloneTopics(node.topics),
		Binding:       node.binding,
		Handler:       node.handler,
		Children:      cloneNodes(node.children),
	}
}

// cloneNodes returns a detached copy of nodes.
func cloneNodes(nodes []Node) []Node {
	if nodes == nil {
		return nil
	}

	out := make([]Node, len(nodes))
	copy(out, nodes)

	return out
}
