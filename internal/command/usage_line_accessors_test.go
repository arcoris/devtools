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

// TestUsageLineAccessors verifies token-oriented helper methods.
func TestUsageLineAccessors(t *testing.T) {
	t.Parallel()

	line := MustUsageLine("bench compare <old> <new> [flags]")

	if got, want := line.TokenCount(), 5; got != want {
		t.Fatalf("TokenCount() = %d, want %d", got, want)
	}

	assertStringSlicesEqual(t, line.Tokens(), []string{"bench", "compare", "<old>", "<new>", "[flags]"})

	first, ok := line.FirstToken()
	if !ok || first != "bench" {
		t.Fatalf("FirstToken() = %q, %t; want %q, true", first, ok, "bench")
	}

	last, ok := line.LastToken()
	if !ok || last != "[flags]" {
		t.Fatalf("LastToken() = %q, %t; want %q, true", last, ok, "[flags]")
	}

	token, ok := line.Token(2)
	if !ok || token != "<old>" {
		t.Fatalf("Token(2) = %q, %t; want %q, true", token, ok, "<old>")
	}

	if _, ok := line.Token(-1); ok {
		t.Fatalf("Token(-1) ok = true, want false")
	}

	if _, ok := line.Token(99); ok {
		t.Fatalf("Token(99) ok = true, want false")
	}

	if !line.HasToken("compare") {
		t.Fatalf("HasToken(compare) = false, want true")
	}

	if line.HasToken("comp") {
		t.Fatalf("HasToken(comp) = true, want false")
	}

	if !line.StartsWithToken("bench") {
		t.Fatalf("StartsWithToken(bench) = false, want true")
	}

	if !line.EndsWithToken("[flags]") {
		t.Fatalf("EndsWithToken([flags]) = false, want true")
	}

	tokens := line.Tokens()
	tokens[0] = "changed"
	if got := line.Tokens()[0]; got != "bench" {
		t.Fatalf("Tokens() returned shared slice; first token = %q", got)
	}
}

// TestUsageLineZeroAccessors verifies zero-value helper behavior.
func TestUsageLineZeroAccessors(t *testing.T) {
	t.Parallel()

	var line UsageLine

	if !line.IsZero() {
		t.Fatalf("zero UsageLine IsZero() = false, want true")
	}

	if line.IsValid() {
		t.Fatalf("zero UsageLine IsValid() = true, want false")
	}

	if got := line.TokenCount(); got != 0 {
		t.Fatalf("zero TokenCount() = %d, want 0", got)
	}

	if tokens := line.Tokens(); tokens != nil {
		t.Fatalf("zero Tokens() = %v, want nil", tokens)
	}

	if _, ok := line.FirstToken(); ok {
		t.Fatalf("zero FirstToken() ok = true, want false")
	}

	if _, ok := line.LastToken(); ok {
		t.Fatalf("zero LastToken() ok = true, want false")
	}
}
