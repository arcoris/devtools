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
	"fmt"
)

// NodeKind describes the structural role of a command-tree node.
//
// NodeKind is framework-neutral. It does not represent any Cobra-specific type.
// A CLI adapter may later map a NodeKind to a concrete adapter command.
type NodeKind string

const (
	// NodeRoot is the root of a command tree.
	//
	// Root nodes have no ID, no Use segment, and a root Path.
	NodeRoot NodeKind = "root"

	// NodeFamily is a non-runnable command-tree branch.
	//
	// Family nodes group child commands or child families. Examples include
	// "bench", "profile", "trace", and "perf".
	NodeFamily NodeKind = "family"

	// NodeCommand is a leaf command-tree node.
	//
	// Command nodes are intended to become runnable after the execution layer is
	// added. Adapter or execution behavior belongs outside this structural model.
	NodeCommand NodeKind = "command"
)

var (
	// ErrInvalidNode reports that a command-tree node definition is malformed.
	ErrInvalidNode = errors.New("command node is invalid")
)

// String returns the stable text form of kind.
func (kind NodeKind) String() string {
	return string(kind)
}

// IsValid reports whether kind is one of the known node kinds.
func (kind NodeKind) IsValid() bool {
	return validateNodeKind(kind) == nil
}

// validateNodeKind validates that kind is one of the known node kinds.
func validateNodeKind(kind NodeKind) error {
	switch kind {
	case NodeRoot, NodeFamily, NodeCommand:
		return nil
	default:
		return fmt.Errorf("%w: unknown node kind %q", ErrInvalidNode, kind)
	}
}
