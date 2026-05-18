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

// BoundArgument is a named positional argument after Binding.Bind.
//
// A BoundArgument is always associated with an Argument declaration. Optional
// arguments may have no values when omitted and when no defaults were declared.
type BoundArgument struct {
	argument Argument
	values   []string
}

// Argument returns the positional argument declaration.
func (bound BoundArgument) Argument() Argument {
	return bound.argument
}

// Name returns the positional argument name.
func (bound BoundArgument) Name() ArgumentName {
	return bound.argument.Name()
}

// Kind returns the positional argument value kind.
func (bound BoundArgument) Kind() OptionKind {
	return bound.argument.Kind()
}

// Values returns detached bound values.
func (bound BoundArgument) Values() []string {
	return cloneBindingStrings(bound.values)
}

// Value returns the first bound value and whether it exists.
func (bound BoundArgument) Value() (string, bool) {
	if len(bound.values) == 0 {
		return "", false
	}

	return bound.values[0], true
}

// MustValue returns the first bound value and panics if it is absent.
func (bound BoundArgument) MustValue() string {
	value, ok := bound.Value()
	if !ok {
		panic("command bound argument value is absent")
	}

	return value
}

// Len returns the number of bound values.
func (bound BoundArgument) Len() int {
	return len(bound.values)
}

// IsEmpty reports whether the argument has no bound values.
func (bound BoundArgument) IsEmpty() bool {
	return len(bound.values) == 0
}

// String returns a compact string representation of bound values.
func (bound BoundArgument) String() string {
	if len(bound.values) == 0 {
		return ""
	}

	if len(bound.values) == 1 {
		return bound.values[0]
	}

	return strings.Join(bound.values, ",")
}

// Validate verifies the bound argument against its declaration.
func (bound BoundArgument) Validate() error {
	if err := bound.argument.Validate(); err != nil {
		return fmt.Errorf("%w: invalid declaration: %w", ErrInvalidBindingValue, err)
	}

	if err := bound.argument.ValidateValues(bound.values...); err != nil {
		return err
	}

	return nil
}
