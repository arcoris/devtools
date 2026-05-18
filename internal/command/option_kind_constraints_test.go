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
	"testing"
)

// TestOptionKindConstraintCapabilities verifies allowed-value and range policies.
func TestOptionKindConstraintCapabilities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		kind                  OptionKind
		canHaveAllowedValues  bool
		requiresAllowedValues bool
		canHaveRange          bool
	}{
		{
			kind:                 OptionKindString,
			canHaveAllowedValues: true,
		},
		{
			kind:                  OptionKindEnum,
			canHaveAllowedValues:  true,
			requiresAllowedValues: true,
		},
		{
			kind:                 OptionKindStringList,
			canHaveAllowedValues: true,
		},
		{
			kind:                  OptionKindEnumList,
			canHaveAllowedValues:  true,
			requiresAllowedValues: true,
		},
		{
			kind:         OptionKindInt,
			canHaveRange: true,
		},
		{
			kind:         OptionKindInt64,
			canHaveRange: true,
		},
		{
			kind:         OptionKindUint,
			canHaveRange: true,
		},
		{
			kind:         OptionKindUint64,
			canHaveRange: true,
		},
		{
			kind:         OptionKindFloat64,
			canHaveRange: true,
		},
		{
			kind:         OptionKindDuration,
			canHaveRange: true,
		},
		{
			kind:         OptionKindIntList,
			canHaveRange: true,
		},
		{
			kind: OptionKindBool,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.kind.String(), func(t *testing.T) {
			t.Parallel()

			if got := tt.kind.CanHaveAllowedValues(); got != tt.canHaveAllowedValues {
				t.Fatalf("CanHaveAllowedValues() = %v, want %v", got, tt.canHaveAllowedValues)
			}

			if got := tt.kind.RequiresAllowedValues(); got != tt.requiresAllowedValues {
				t.Fatalf("RequiresAllowedValues() = %v, want %v", got, tt.requiresAllowedValues)
			}

			if got := tt.kind.CanHaveRange(); got != tt.canHaveRange {
				t.Fatalf("CanHaveRange() = %v, want %v", got, tt.canHaveRange)
			}
		})
	}
}
