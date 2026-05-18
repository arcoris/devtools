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
	"testing"
)

// TestArgumentCardinalityValidation verifies cardinality validation.
func TestArgumentCardinalityValidation(t *testing.T) {
	t.Parallel()

	single := ArgumentCardinalitySingle
	if !single.IsKnown() || !single.IsSingle() {
		t.Fatalf("single predicates are invalid")
	}

	variadic := ArgumentCardinalityVariadic
	if !variadic.IsKnown() || !variadic.IsVariadic() {
		t.Fatalf("variadic predicates are invalid")
	}

	if got, want := ArgumentCardinality("").OrDefault(), ArgumentCardinalitySingle; got != want {
		t.Fatalf("zero cardinality OrDefault() = %q, want %q", got, want)
	}

	invalid := ArgumentCardinality("many")
	if err := invalid.Validate(); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("invalid cardinality Validate() error = %v, want ErrInvalidArgument", err)
	}
}
