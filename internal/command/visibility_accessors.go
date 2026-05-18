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

// String returns the canonical string representation of the visibility value.
func (visibility Visibility) String() string {
	return string(visibility)
}

// Key returns the stable map key for visibility.
func (visibility Visibility) Key() string {
	return string(visibility)
}

// IsZero reports whether the visibility value has not been set.
func (visibility Visibility) IsZero() bool {
	return visibility == ""
}

// OrDefault returns DefaultVisibility when visibility is zero.
//
// This method is useful in constructors that support omitted visibility values
// but still need a concrete policy before validation or adapter rendering.
func (visibility Visibility) OrDefault() Visibility {
	if visibility.IsZero() {
		return DefaultVisibility()
	}

	return visibility
}

// Equal reports whether two visibility values are exactly the same state.
func (visibility Visibility) Equal(other Visibility) bool {
	return visibility == other
}

// IsKnown reports whether visibility is one of the supported non-zero states.
func (visibility Visibility) IsKnown() bool {
	for _, candidate := range knownVisibilities {
		if visibility == candidate {
			return true
		}
	}

	return false
}

// IsValid reports whether visibility satisfies the visibility grammar.
func (visibility Visibility) IsValid() bool {
	return visibility.Validate() == nil
}

// IsPublic reports whether visibility is public.
func (visibility Visibility) IsPublic() bool {
	return visibility == VisibilityPublic
}

// IsHidden reports whether visibility is hidden.
func (visibility Visibility) IsHidden() bool {
	return visibility == VisibilityHidden
}

// IsInternal reports whether visibility is internal.
func (visibility Visibility) IsInternal() bool {
	return visibility == VisibilityInternal
}
