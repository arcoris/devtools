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

// TestArgumentRequirementValidation verifies requirement validation.
func TestArgumentRequirementValidation(t *testing.T) {
	t.Parallel()

	required := ArgumentRequirementRequired
	if !required.IsKnown() || !required.IsRequired() {
		t.Fatalf("required predicates are invalid")
	}

	optional := ArgumentRequirementOptional
	if !optional.IsKnown() || !optional.IsOptional() {
		t.Fatalf("optional predicates are invalid")
	}

	if got, want := ArgumentRequirement("").OrDefault(), ArgumentRequirementRequired; got != want {
		t.Fatalf("zero requirement OrDefault() = %q, want %q", got, want)
	}

	invalid := ArgumentRequirement("mandatory")
	if err := invalid.Validate(); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("invalid requirement Validate() error = %v, want ErrInvalidArgument", err)
	}
}
