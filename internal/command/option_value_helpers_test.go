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

import "testing"

// TestIsStringLikeElementKind verifies string-like element-kind predicate.
func TestIsStringLikeElementKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		kind OptionKind
		want bool
	}{
		{kind: OptionKindString, want: true},
		{kind: OptionKindEnum, want: true},
		{kind: OptionKindInt, want: false},
		{kind: OptionKindStringList, want: false},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.kind.String(), func(t *testing.T) {
			t.Parallel()

			if got := isStringLikeElementKind(tt.kind); got != tt.want {
				t.Fatalf("isStringLikeElementKind(%q) = %v, want %v", tt.kind, got, tt.want)
			}
		})
	}
}

// TestStringSlicesEqual verifies ordered string slice comparison.
func TestStringSlicesEqual(t *testing.T) {
	t.Parallel()

	if !stringSlicesEqual([]string{"a", "b"}, []string{"a", "b"}) {
		t.Fatalf("stringSlicesEqual(equal) = false, want true")
	}

	if stringSlicesEqual([]string{"a", "b"}, []string{"b", "a"}) {
		t.Fatalf("stringSlicesEqual(different order) = true, want false")
	}

	if stringSlicesEqual([]string{"a"}, []string{"a", "b"}) {
		t.Fatalf("stringSlicesEqual(different length) = true, want false")
	}
}
