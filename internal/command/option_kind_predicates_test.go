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

// TestOptionKindPredicates verifies direct kind predicates.
func TestOptionKindPredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		kind            OptionKind
		zero            bool
		known           bool
		boolKind        bool
		stringKind      bool
		enumKind        bool
		enumLike        bool
		integer         bool
		signedInteger   bool
		unsignedInteger bool
		float           bool
		duration        bool
		numeric         bool
		scalar          bool
		list            bool
		repeatable      bool
		requiresValue   bool
		implicitBool    bool
	}{
		{
			name:          "zero",
			kind:          "",
			zero:          true,
			known:         false,
			requiresValue: false,
		},
		{
			name:          "bool",
			kind:          OptionKindBool,
			known:         true,
			boolKind:      true,
			scalar:        true,
			requiresValue: false,
			implicitBool:  true,
		},
		{
			name:          "string",
			kind:          OptionKindString,
			known:         true,
			stringKind:    true,
			scalar:        true,
			requiresValue: true,
		},
		{
			name:          "enum",
			kind:          OptionKindEnum,
			known:         true,
			enumKind:      true,
			enumLike:      true,
			scalar:        true,
			requiresValue: true,
		},
		{
			name:          "int",
			kind:          OptionKindInt,
			known:         true,
			integer:       true,
			signedInteger: true,
			numeric:       true,
			scalar:        true,
			requiresValue: true,
		},
		{
			name:          "int64",
			kind:          OptionKindInt64,
			known:         true,
			integer:       true,
			signedInteger: true,
			numeric:       true,
			scalar:        true,
			requiresValue: true,
		},
		{
			name:            "uint",
			kind:            OptionKindUint,
			known:           true,
			integer:         true,
			unsignedInteger: true,
			numeric:         true,
			scalar:          true,
			requiresValue:   true,
		},
		{
			name:            "uint64",
			kind:            OptionKindUint64,
			known:           true,
			integer:         true,
			unsignedInteger: true,
			numeric:         true,
			scalar:          true,
			requiresValue:   true,
		},
		{
			name:          "float64",
			kind:          OptionKindFloat64,
			known:         true,
			float:         true,
			numeric:       true,
			scalar:        true,
			requiresValue: true,
		},
		{
			name:          "duration",
			kind:          OptionKindDuration,
			known:         true,
			duration:      true,
			scalar:        true,
			requiresValue: true,
		},
		{
			name:          "string-list",
			kind:          OptionKindStringList,
			known:         true,
			list:          true,
			repeatable:    true,
			requiresValue: true,
		},
		{
			name:          "enum-list",
			kind:          OptionKindEnumList,
			known:         true,
			enumLike:      true,
			list:          true,
			repeatable:    true,
			requiresValue: true,
		},
		{
			name:          "int-list",
			kind:          OptionKindIntList,
			known:         true,
			list:          true,
			repeatable:    true,
			requiresValue: true,
		},
		{
			name:          "unknown",
			kind:          OptionKind("path"),
			known:         false,
			requiresValue: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.kind.IsZero(); got != tt.zero {
				t.Fatalf("IsZero() = %v, want %v", got, tt.zero)
			}

			if got := tt.kind.IsKnown(); got != tt.known {
				t.Fatalf("IsKnown() = %v, want %v", got, tt.known)
			}

			if got := tt.kind.IsBool(); got != tt.boolKind {
				t.Fatalf("IsBool() = %v, want %v", got, tt.boolKind)
			}

			if got := tt.kind.IsString(); got != tt.stringKind {
				t.Fatalf("IsString() = %v, want %v", got, tt.stringKind)
			}

			if got := tt.kind.IsEnum(); got != tt.enumKind {
				t.Fatalf("IsEnum() = %v, want %v", got, tt.enumKind)
			}

			if got := tt.kind.IsEnumLike(); got != tt.enumLike {
				t.Fatalf("IsEnumLike() = %v, want %v", got, tt.enumLike)
			}

			if got := tt.kind.IsInteger(); got != tt.integer {
				t.Fatalf("IsInteger() = %v, want %v", got, tt.integer)
			}

			if got := tt.kind.IsSignedInteger(); got != tt.signedInteger {
				t.Fatalf("IsSignedInteger() = %v, want %v", got, tt.signedInteger)
			}

			if got := tt.kind.IsUnsignedInteger(); got != tt.unsignedInteger {
				t.Fatalf("IsUnsignedInteger() = %v, want %v", got, tt.unsignedInteger)
			}

			if got := tt.kind.IsFloat(); got != tt.float {
				t.Fatalf("IsFloat() = %v, want %v", got, tt.float)
			}

			if got := tt.kind.IsDuration(); got != tt.duration {
				t.Fatalf("IsDuration() = %v, want %v", got, tt.duration)
			}

			if got := tt.kind.IsNumeric(); got != tt.numeric {
				t.Fatalf("IsNumeric() = %v, want %v", got, tt.numeric)
			}

			if got := tt.kind.IsScalar(); got != tt.scalar {
				t.Fatalf("IsScalar() = %v, want %v", got, tt.scalar)
			}

			if got := tt.kind.IsList(); got != tt.list {
				t.Fatalf("IsList() = %v, want %v", got, tt.list)
			}

			if got := tt.kind.IsRepeatable(); got != tt.repeatable {
				t.Fatalf("IsRepeatable() = %v, want %v", got, tt.repeatable)
			}

			if got := tt.kind.RequiresValue(); got != tt.requiresValue {
				t.Fatalf("RequiresValue() = %v, want %v", got, tt.requiresValue)
			}

			if got := tt.kind.AllowsImplicitBoolean(); got != tt.implicitBool {
				t.Fatalf("AllowsImplicitBoolean() = %v, want %v", got, tt.implicitBool)
			}
		})
	}
}
