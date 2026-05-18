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

// String returns the canonical string representation of the option source.
func (source OptionSource) String() string {
	return string(source)
}

// IsZero reports whether the option source has not been set.
func (source OptionSource) IsZero() bool {
	return source == ""
}

// IsKnown reports whether source is one of the supported non-zero option
// sources.
func (source OptionSource) IsKnown() bool {
	for _, candidate := range knownOptionSources {
		if source == candidate {
			return true
		}
	}

	return false
}

// IsValid reports whether source satisfies the option-source grammar.
func (source OptionSource) IsValid() bool {
	return source.Validate() == nil
}

// Validate verifies that source is a supported non-zero option source.
func (source OptionSource) Validate() error {
	if source == "" {
		return ErrEmptyOptionSource
	}

	if source.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported value %q", ErrInvalidOptionSource, source)
}

// KnownOptionSources returns all supported option sources in stable precedence
// order, from lowest to highest.
//
// The returned slice is detached and can be safely modified by the caller.
func KnownOptionSources() []OptionSource {
	return cloneOptionSources(knownOptionSources)
}
