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

import "fmt"

// OptionScope describes where an option declaration is visible in the command
// tree.
type OptionScope string

const (
	// OptionScopeLocal means the option applies only to the command node that
	// declares it.
	OptionScopeLocal OptionScope = "local"

	// OptionScopeSubtree means the option applies to the declaring node and its
	// descendants.
	OptionScopeSubtree OptionScope = "subtree"

	// OptionScopeGlobal means the option applies to the whole command tree.
	OptionScopeGlobal OptionScope = "global"
)

// NewOptionScope validates raw and returns it as an OptionScope.
func NewOptionScope(raw string) (OptionScope, error) {
	scope := OptionScope(raw)
	if err := scope.Validate(); err != nil {
		return "", err
	}

	return scope, nil
}

// ParseOptionScope is an alias for NewOptionScope.
func ParseOptionScope(raw string) (OptionScope, error) {
	return NewOptionScope(raw)
}

// MustOptionScope validates raw and returns it as an OptionScope.
func MustOptionScope(raw string) OptionScope {
	scope, err := NewOptionScope(raw)
	if err != nil {
		panic(err)
	}

	return scope
}

// String returns the canonical string representation of the scope.
func (scope OptionScope) String() string {
	return string(scope)
}

// IsZero reports whether the scope has not been set.
func (scope OptionScope) IsZero() bool {
	return scope == ""
}

// OrDefault returns OptionScopeLocal when scope is zero.
func (scope OptionScope) OrDefault() OptionScope {
	if scope.IsZero() {
		return OptionScopeLocal
	}

	return scope
}

// IsKnown reports whether scope is one of the supported non-zero states.
func (scope OptionScope) IsKnown() bool {
	switch scope {
	case OptionScopeLocal, OptionScopeSubtree, OptionScopeGlobal:
		return true
	default:
		return false
	}
}

// IsValid reports whether scope satisfies policy rules.
func (scope OptionScope) IsValid() bool {
	return scope.Validate() == nil
}

// Validate verifies that scope is a supported non-zero state.
func (scope OptionScope) Validate() error {
	if scope == "" {
		return fmt.Errorf("%w: scope is empty", ErrInvalidOptionPolicy)
	}

	if scope.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported scope %q", ErrInvalidOptionPolicy, scope)
}

// IsLocal reports whether the option applies only to the declaring command node.
func (scope OptionScope) IsLocal() bool {
	return scope == OptionScopeLocal
}

// IsSubtree reports whether the option applies to the declaring node subtree.
func (scope OptionScope) IsSubtree() bool {
	return scope == OptionScopeSubtree
}

// IsGlobal reports whether the option applies to the whole command tree.
func (scope OptionScope) IsGlobal() bool {
	return scope == OptionScopeGlobal
}

// IsInheritedByChildren reports whether a declaration with this scope is
// inherited by descendant command nodes.
func (scope OptionScope) IsInheritedByChildren() bool {
	switch scope {
	case OptionScopeSubtree, OptionScopeGlobal:
		return true
	default:
		return false
	}
}
