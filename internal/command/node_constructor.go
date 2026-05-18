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

// NewNode validates spec and returns a detached Node value.
func NewNode(spec NodeSpec) (Node, error) {
	documentation, err := newNodeDocumentation(spec)
	if err != nil {
		return Node{}, err
	}

	metadata, err := newNodeMetadata(spec)
	if err != nil {
		return Node{}, err
	}

	node := Node{
		kind:          spec.Kind,
		id:            spec.ID,
		path:          spec.Path,
		use:           spec.Use,
		aliases:       cloneStringSlice(spec.Aliases),
		documentation: documentation,
		metadata:      metadata,
		visibility:    newNodeVisibility(spec),
		group:         spec.Group,
		topics:        cloneTopics(spec.Topics),
		binding:       spec.Binding,
		handler:       spec.Handler,
		children:      cloneNodes(spec.Children),
	}

	if err := node.Validate(); err != nil {
		return Node{}, err
	}

	return node, nil
}

// MustNode validates spec and returns a Node.
//
// MustNode panics on invalid input. It is intended for static command
// definitions and tests where invalid nodes are programmer errors.
func MustNode(spec NodeSpec) Node {
	node, err := NewNode(spec)
	if err != nil {
		panic(err)
	}

	return node
}

// NewRootNode returns a root command-tree node with optional children.
func NewRootNode(children ...Node) (Node, error) {
	return NewNode(NodeSpec{
		Kind:     NodeRoot,
		Path:     RootPath(),
		Children: children,
	})
}

// MustRootNode returns a root command-tree node and panics on invalid children.
func MustRootNode(children ...Node) Node {
	return MustNode(NodeSpec{
		Kind:     NodeRoot,
		Path:     RootPath(),
		Children: children,
	})
}

// NewFamilyNode returns a command family node.
//
// A family node is a non-runnable branch in the command tree. It must contain
// at least one child.
func NewFamilyNode(id ID, path Path, use string, children ...Node) (Node, error) {
	return NewNode(NodeSpec{
		Kind:     NodeFamily,
		ID:       id,
		Path:     path,
		Use:      use,
		Children: children,
	})
}

// MustFamilyNode returns a command family node and panics on invalid input.
func MustFamilyNode(id ID, path Path, use string, children ...Node) Node {
	return MustNode(NodeSpec{
		Kind:     NodeFamily,
		ID:       id,
		Path:     path,
		Use:      use,
		Children: children,
	})
}

// NewCommandNode returns a leaf command node.
//
// This constructor only creates the structural node. Runnable command behavior
// belongs outside this value object.
func NewCommandNode(id ID, path Path, use string) (Node, error) {
	return NewNode(NodeSpec{
		Kind: NodeCommand,
		ID:   id,
		Path: path,
		Use:  use,
	})
}

// MustCommandNode returns a leaf command node and panics on invalid input.
func MustCommandNode(id ID, path Path, use string) Node {
	return MustNode(NodeSpec{
		Kind: NodeCommand,
		ID:   id,
		Path: path,
		Use:  use,
	})
}
