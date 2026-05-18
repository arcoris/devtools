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

// OrDefault returns ArgumentCardinalitySingle when cardinality is zero.
func (cardinality ArgumentCardinality) OrDefault() ArgumentCardinality {
	if cardinality == "" {
		return ArgumentCardinalitySingle
	}

	return cardinality
}

// String returns the canonical cardinality string.
func (cardinality ArgumentCardinality) String() string {
	return string(cardinality)
}

// IsKnown reports whether cardinality is a supported non-zero state.
func (cardinality ArgumentCardinality) IsKnown() bool {
	switch cardinality {
	case ArgumentCardinalitySingle, ArgumentCardinalityVariadic:
		return true
	default:
		return false
	}
}

// Validate verifies cardinality.
func (cardinality ArgumentCardinality) Validate() error {
	if cardinality == "" {
		return fmt.Errorf("%w: cardinality is empty", ErrInvalidArgument)
	}

	if cardinality.IsKnown() {
		return nil
	}

	return fmt.Errorf("%w: unsupported cardinality %q", ErrInvalidArgument, cardinality)
}

// IsSingle reports whether cardinality is single.
func (cardinality ArgumentCardinality) IsSingle() bool {
	return cardinality == ArgumentCardinalitySingle
}

// IsVariadic reports whether cardinality is variadic.
func (cardinality ArgumentCardinality) IsVariadic() bool {
	return cardinality == ArgumentCardinalityVariadic
}
