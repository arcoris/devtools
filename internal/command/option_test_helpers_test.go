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

func assertOptionValueStrings(t *testing.T, got []string, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("slice length = %d, want %d; got %v, want %v", len(got), len(want), got, want)
	}

	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("slice[%d] = %q, want %q; got %v, want %v", index, got[index], want[index], got, want)
		}
	}
}

func assertOptionValuePanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatalf("function did not panic")
		}
	}()

	fn()
}
