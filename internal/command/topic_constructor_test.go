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

func TestNewTopicAcceptsValidTopics(t *testing.T) {
	t.Parallel()

	tests := []string{
		"profiling",
		"tracing",
		"perf",
		"diagnostics",
		"diagnostics.perf",
		"benchmark.compare",
		"release-notes",
		"a",
		"a1",
		"a-1",
		"a.b.c",
	}

	for _, raw := range tests {
		raw := raw

		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			topic, err := NewTopic(raw)
			if err != nil {
				t.Fatalf("NewTopic(%q) returned unexpected error: %v", raw, err)
			}

			if got := topic.String(); got != raw {
				t.Fatalf("String() = %q, want %q", got, raw)
			}

			if !topic.IsValid() {
				t.Fatalf("IsValid() = false, want true")
			}
		})
	}
}

func TestParseTopicIsAliasForNewTopic(t *testing.T) {
	t.Parallel()

	const raw = "diagnostics.perf"

	fromNew, err := NewTopic(raw)
	if err != nil {
		t.Fatalf("NewTopic(%q) returned unexpected error: %v", raw, err)
	}

	fromParse, err := ParseTopic(raw)
	if err != nil {
		t.Fatalf("ParseTopic(%q) returned unexpected error: %v", raw, err)
	}

	if fromParse != fromNew {
		t.Fatalf("ParseTopic(%q) = %q, want %q", raw, fromParse, fromNew)
	}
}

func TestNewTopicPartsBuildsCanonicalTopic(t *testing.T) {
	t.Parallel()

	topic, err := NewTopicParts("diagnostics", "perf")
	if err != nil {
		t.Fatalf("NewTopicParts() returned unexpected error: %v", err)
	}

	if got, want := topic, MustTopic("diagnostics.perf"); got != want {
		t.Fatalf("NewTopicParts() = %q, want %q", got, want)
	}
}

func TestMustTopicReturnsValidTopic(t *testing.T) {
	t.Parallel()

	topic := MustTopic("diagnostics.perf")

	if got, want := topic.String(), "diagnostics.perf"; got != want {
		t.Fatalf("MustTopic returned %q, want %q", got, want)
	}
}

func TestMustTopicPanicsForInvalidTopic(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustTopic("Diagnostics.Perf")
	})
}

func TestMustTopicPartsPanicsForInvalidSegment(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustTopicParts("diagnostics", "Perf")
	})
}
