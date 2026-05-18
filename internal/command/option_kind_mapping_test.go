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

// TestOptionKindElementKind verifies list-to-scalar mapping.
func TestOptionKindElementKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		kind OptionKind
		want OptionKind
	}{
		{kind: OptionKindBool, want: OptionKindBool},
		{kind: OptionKindString, want: OptionKindString},
		{kind: OptionKindEnum, want: OptionKindEnum},
		{kind: OptionKindInt, want: OptionKindInt},
		{kind: OptionKindStringList, want: OptionKindString},
		{kind: OptionKindEnumList, want: OptionKindEnum},
		{kind: OptionKindIntList, want: OptionKindInt},
		{kind: OptionKindInt64List, want: OptionKindInt64},
		{kind: OptionKindUintList, want: OptionKindUint},
		{kind: OptionKindUint64List, want: OptionKindUint64},
		{kind: OptionKindFloat64List, want: OptionKindFloat64},
		{kind: OptionKindDurationList, want: OptionKindDuration},
		{kind: OptionKind("unknown"), want: OptionKind("unknown")},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.kind.String(), func(t *testing.T) {
			t.Parallel()

			if got := tt.kind.ElementKind(); got != tt.want {
				t.Fatalf("ElementKind() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestOptionKindListKind verifies scalar-to-list mapping.
func TestOptionKindListKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		kind OptionKind
		want OptionKind
	}{
		{kind: OptionKindBool, want: OptionKindBool},
		{kind: OptionKindString, want: OptionKindStringList},
		{kind: OptionKindEnum, want: OptionKindEnumList},
		{kind: OptionKindInt, want: OptionKindIntList},
		{kind: OptionKindInt64, want: OptionKindInt64List},
		{kind: OptionKindUint, want: OptionKindUintList},
		{kind: OptionKindUint64, want: OptionKindUint64List},
		{kind: OptionKindFloat64, want: OptionKindFloat64List},
		{kind: OptionKindDuration, want: OptionKindDurationList},
		{kind: OptionKindStringList, want: OptionKindStringList},
		{kind: OptionKind("unknown"), want: OptionKind("unknown")},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.kind.String(), func(t *testing.T) {
			t.Parallel()

			if got := tt.kind.ListKind(); got != tt.want {
				t.Fatalf("ListKind() = %q, want %q", got, tt.want)
			}
		})
	}
}
