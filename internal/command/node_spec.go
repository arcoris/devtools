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
	// An empty value means the node is not marked as deprecated at this layer.
	Deprecated string

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
// Node currently models structure and metadata only. Runnable behavior belongs
// to the execution layer and should be added through action/lifecycle types, not
// through adapter-specific command objects.
type Node struct {
	kind       NodeKind
	id         ID
	path       Path
	use        string
	aliases    []string
	short      string
	long       string
	example    string
	hidden     bool
	deprecated string
	children   []Node
}

// Spec returns a detached construction spec for node.
func (node Node) Spec() NodeSpec {
	return node.spec()
}

// spec returns a detached construction spec for node.
func (node Node) spec() NodeSpec {
	return NodeSpec{
		Kind:       node.kind,
		ID:         node.id,
		Path:       node.path,
		Use:        node.use,
		Aliases:    cloneStringSlice(node.aliases),
		Short:      node.short,
		Long:       node.long,
		Example:    node.example,
		Hidden:     node.hidden,
		Deprecated: node.deprecated,
		Children:   cloneNodes(node.children),
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
