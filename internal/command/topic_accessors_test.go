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

func TestTopicBasicAccessors(t *testing.T) {
	t.Parallel()

	topic := MustTopic("diagnostics.perf")

	if got, want := topic.String(), "diagnostics.perf"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	if got, want := topic.Key(), "diagnostics.perf"; got != want {
		t.Fatalf("Key() = %q, want %q", got, want)
	}

	if !topic.Equal(MustTopic("diagnostics.perf")) {
		t.Fatalf("Equal() = false, want true")
	}

	if topic.Equal(MustTopic("diagnostics.trace")) {
		t.Fatalf("Equal(other) = true, want false")
	}
}

func TestTopicIsZero(t *testing.T) {
	t.Parallel()

	var zero Topic
	if !zero.IsZero() {
		t.Fatalf("zero topic IsZero() = false, want true")
	}

	if MustTopic("profiling").IsZero() {
		t.Fatalf("non-zero topic IsZero() = true, want false")
	}
}

func TestTopicPartsReturnsDetachedCopy(t *testing.T) {
	t.Parallel()

	topic := MustTopic("diagnostics.perf")

	parts := topic.Parts()
	assertStringSlicesEqual(t, parts, []string{"diagnostics", "perf"})

	parts[0] = "changed"

	again := topic.Parts()
	assertStringSlicesEqual(t, again, []string{"diagnostics", "perf"})
}

func TestTopicPartsForZeroTopic(t *testing.T) {
	t.Parallel()

	var topic Topic

	if parts := topic.Parts(); parts != nil {
		t.Fatalf("zero topic Parts() = %#v, want nil", parts)
	}
}

func TestTopicDepthAndLen(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		topic Topic
		want  int
	}{
		{name: "zero", topic: "", want: 0},
		{name: "root", topic: MustTopic("profiling"), want: 1},
		{name: "child", topic: MustTopic("diagnostics.perf"), want: 2},
		{name: "deep", topic: MustTopic("a.b.c"), want: 3},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.topic.Depth(); got != test.want {
				t.Fatalf("Depth() = %d, want %d", got, test.want)
			}

			if got := test.topic.Len(); got != test.want {
				t.Fatalf("Len() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestTopicAt(t *testing.T) {
	t.Parallel()

	topic := MustTopic("diagnostics.perf.stat")

	tests := []struct {
		name  string
		index int
		want  string
		ok    bool
	}{
		{name: "negative", index: -1, want: "", ok: false},
		{name: "first", index: 0, want: "diagnostics", ok: true},
		{name: "second", index: 1, want: "perf", ok: true},
		{name: "third", index: 2, want: "stat", ok: true},
		{name: "out of range", index: 3, want: "", ok: false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, ok := topic.At(test.index)
			if ok != test.ok {
				t.Fatalf("At(%d) ok = %v, want %v", test.index, ok, test.ok)
			}

			if got != test.want {
				t.Fatalf("At(%d) = %q, want %q", test.index, got, test.want)
			}
		})
	}
}

func TestTopicLeaf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		topic Topic
		want  string
	}{
		{name: "zero", topic: "", want: ""},
		{name: "root", topic: MustTopic("profiling"), want: "profiling"},
		{name: "child", topic: MustTopic("diagnostics.perf"), want: "perf"},
		{name: "deep", topic: MustTopic("a.b.c"), want: "c"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.topic.Leaf(); got != test.want {
				t.Fatalf("Leaf() = %q, want %q", got, test.want)
			}
		})
	}
}
