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

// TestUsageAccessors verifies read-only helper behavior.
func TestUsageAccessors(t *testing.T) {
	t.Parallel()

	usage := MustUsage(UsageSpec{
		Syntax: "bench run",
		Alternatives: []string{
			"bench smoke",
			"bench run --suite <name>",
		},
	})

	if got, want := usage.Syntax().String(), "bench run"; got != want {
		t.Fatalf("Syntax() = %q, want %q", got, want)
	}

	if got, want := usage.Len(), 3; got != want {
		t.Fatalf("Len() = %d, want %d", got, want)
	}

	if got, want := usage.LineCount(), 3; got != want {
		t.Fatalf("LineCount() = %d, want %d", got, want)
	}

	if got, want := usage.AlternativeCount(), 2; got != want {
		t.Fatalf("AlternativeCount() = %d, want %d", got, want)
	}

	if !usage.HasAlternatives() {
		t.Fatalf("HasAlternatives() = false, want true")
	}

	alternative, ok := usage.Alternative(1)
	if !ok || alternative.String() != "bench run --suite <name>" {
		t.Fatalf("Alternative(1) = %q, %t; want %q, true", alternative, ok, "bench run --suite <name>")
	}

	if _, ok := usage.Alternative(-1); ok {
		t.Fatalf("Alternative(-1) ok = true, want false")
	}

	if _, ok := usage.Alternative(99); ok {
		t.Fatalf("Alternative(99) ok = true, want false")
	}

	line, ok := usage.Line(0)
	if !ok || line.String() != "bench run" {
		t.Fatalf("Line(0) = %q, %t; want %q, true", line, ok, "bench run")
	}

	line, ok = usage.Line(2)
	if !ok || line.String() != "bench run --suite <name>" {
		t.Fatalf("Line(2) = %q, %t; want %q, true", line, ok, "bench run --suite <name>")
	}

	if _, ok := usage.Line(-1); ok {
		t.Fatalf("Line(-1) ok = true, want false")
	}

	if _, ok := usage.Line(99); ok {
		t.Fatalf("Line(99) ok = true, want false")
	}

	assertStringSlicesEqual(t, usage.LineStrings(), []string{
		"bench run",
		"bench smoke",
		"bench run --suite <name>",
	})

	assertStringSlicesEqual(t, usageLineStrings(usage.SortedLines()), []string{
		"bench run",
		"bench run --suite <name>",
		"bench smoke",
	})

	if !usage.Contains("  bench   smoke ") {
		t.Fatalf("Contains(spaced alternative) = false, want true")
	}

	if !usage.ContainsLine(MustUsageLine("bench run")) {
		t.Fatalf("ContainsLine(primary) = false, want true")
	}

	if usage.Contains("bench\nrun") {
		t.Fatalf("Contains(invalid line) = true, want false")
	}
}

// TestUsageZeroAccessors verifies zero-value helper behavior.
func TestUsageZeroAccessors(t *testing.T) {
	t.Parallel()

	var usage Usage

	if !usage.IsZero() {
		t.Fatalf("zero Usage IsZero() = false, want true")
	}

	if usage.IsValid() {
		t.Fatalf("zero Usage IsValid() = true, want false")
	}

	if got := usage.String(); got != "" {
		t.Fatalf("zero String() = %q, want empty", got)
	}

	if got := usage.Len(); got != 0 {
		t.Fatalf("zero Len() = %d, want 0", got)
	}

	if lines := usage.Lines(); lines != nil {
		t.Fatalf("zero Lines() = %v, want nil", lines)
	}
}
