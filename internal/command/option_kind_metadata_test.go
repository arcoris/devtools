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

// TestKnownOptionKindsReturnsDetachedStableOrder verifies stable declaration order.
func TestKnownOptionKindsReturnsDetachedStableOrder(t *testing.T) {
	t.Parallel()

	got := KnownOptionKinds()

	want := []OptionKind{
		OptionKindBool,
		OptionKindString,
		OptionKindEnum,
		OptionKindInt,
		OptionKindInt64,
		OptionKindUint,
		OptionKindUint64,
		OptionKindFloat64,
		OptionKindDuration,
		OptionKindStringList,
		OptionKindEnumList,
		OptionKindIntList,
		OptionKindInt64List,
		OptionKindUintList,
		OptionKindUint64List,
		OptionKindFloat64List,
		OptionKindDurationList,
	}

	if len(got) != len(want) {
		t.Fatalf("KnownOptionKinds length = %d, want %d", len(got), len(want))
	}

	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("KnownOptionKinds()[%d] = %q, want %q", index, got[index], want[index])
		}
	}

	got[0] = OptionKind("changed")

	again := KnownOptionKinds()
	if again[0] != OptionKindBool {
		t.Fatalf("KnownOptionKinds returned mutable state: got %q, want %q", again[0], OptionKindBool)
	}
}

// TestOptionKindValueMetavar verifies default metavar rendering.
func TestOptionKindValueMetavar(t *testing.T) {
	t.Parallel()

	tests := []struct {
		kind OptionKind
		want string
	}{
		{kind: OptionKindBool, want: "BOOL"},
		{kind: OptionKindString, want: "STRING"},
		{kind: OptionKindEnum, want: "VALUE"},
		{kind: OptionKindInt, want: "INT"},
		{kind: OptionKindInt64, want: "INT"},
		{kind: OptionKindUint, want: "UINT"},
		{kind: OptionKindUint64, want: "UINT"},
		{kind: OptionKindFloat64, want: "FLOAT"},
		{kind: OptionKindDuration, want: "DURATION"},
		{kind: OptionKindStringList, want: "STRING"},
		{kind: OptionKindEnumList, want: "VALUE"},
		{kind: OptionKindIntList, want: "INT"},
		{kind: OptionKindInt64List, want: "INT"},
		{kind: OptionKindUintList, want: "UINT"},
		{kind: OptionKindUint64List, want: "UINT"},
		{kind: OptionKindFloat64List, want: "FLOAT"},
		{kind: OptionKindDurationList, want: "DURATION"},
		{kind: OptionKind("unknown"), want: "VALUE"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.kind.String(), func(t *testing.T) {
			t.Parallel()

			if got := tt.kind.ValueMetavar(); got != tt.want {
				t.Fatalf("ValueMetavar() = %q, want %q", got, tt.want)
			}
		})
	}
}
