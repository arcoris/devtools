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

// Kind returns the node kind.
func (node Node) Kind() NodeKind {
	return node.kind
}

// ID returns the stable machine-facing node identifier.
//
// Root nodes return the zero ID.
func (node Node) ID() ID {
	return node.id
}

// Path returns the logical command-tree path.
func (node Node) Path() Path {
	return node.path
}

// Use returns the canonical adapter-facing command segment.
//
// Root nodes return an empty string.
func (node Node) Use() string {
	return node.use
}

// Aliases returns a detached copy of adapter-facing alias segments.
func (node Node) Aliases() []string {
	return cloneStringSlice(node.aliases)
}

// Short returns the short human-facing node description.
func (node Node) Short() string {
	return node.short
}

// Long returns the long human-facing node description.
func (node Node) Long() string {
	return node.long
}

// Example returns human-facing usage examples.
func (node Node) Example() string {
	return node.example
}

// Hidden reports whether adapters should hide this node from default help.
func (node Node) Hidden() bool {
	return node.hidden
}

// Deprecated returns the optional deprecation message.
func (node Node) Deprecated() string {
	return node.deprecated
}

// IsDeprecated reports whether the node has a deprecation message.
func (node Node) IsDeprecated() bool {
	return node.deprecated != ""
}

// IsRoot reports whether the node is the command-tree root.
func (node Node) IsRoot() bool {
	return node.kind == NodeRoot
}

// IsFamily reports whether the node is a command family.
func (node Node) IsFamily() bool {
	return node.kind == NodeFamily
}

// IsCommand reports whether the node is a leaf command.
func (node Node) IsCommand() bool {
	return node.kind == NodeCommand
}

// IsLeaf reports whether the node has no children.
func (node Node) IsLeaf() bool {
	return len(node.children) == 0
}

// HasChildren reports whether the node has one or more children.
func (node Node) HasChildren() bool {
	return len(node.children) > 0
}

// HasAlias reports whether alias is registered as an adapter-facing alias for
// this node.
func (node Node) HasAlias(alias string) bool {
	for _, candidate := range node.aliases {
		if candidate == alias {
			return true
		}
	}

	return false
}

// MatchesSegment reports whether segment matches the node's canonical Use
// segment or one of its aliases.
//
// Root nodes never match a segment because they do not have a Use segment.
func (node Node) MatchesSegment(segment string) bool {
	if node.use == "" {
		return false
	}

	if node.use == segment {
		return true
	}

	return node.HasAlias(segment)
}
