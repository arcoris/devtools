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

// String returns the canonical string representation of the option kind.
func (kind OptionKind) String() string {
	return string(kind)
}

// IsZero reports whether the option kind has not been set.
func (kind OptionKind) IsZero() bool {
	return kind == ""
}

// IsKnown reports whether kind is one of the supported non-zero option kinds.
func (kind OptionKind) IsKnown() bool {
	for _, candidate := range knownOptionKinds {
		if kind == candidate {
			return true
		}
	}

	return false
}

// IsValid reports whether kind satisfies the option-kind grammar.
func (kind OptionKind) IsValid() bool {
	return kind.Validate() == nil
}

// Validate verifies that kind is a supported non-zero option kind.
func (kind OptionKind) Validate() error {
	if kind == "" {
		return ErrEmptyOptionKind
	}

	if kind.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported value %q", ErrInvalidOptionKind, kind)
}

// KnownOptionKinds returns all supported option kinds in stable declaration
// order.
//
// The returned slice is detached and can be safely modified by the caller.
func KnownOptionKinds() []OptionKind {
	return cloneOptionKinds(knownOptionKinds)
}
