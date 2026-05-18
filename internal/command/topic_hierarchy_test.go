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
	"errors"
	"testing"
)

func TestTopicParent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		topic     Topic
		want      Topic
		wantFound bool
	}{
		{name: "zero", topic: "", want: "", wantFound: false},
		{name: "root", topic: MustTopic("profiling"), want: "", wantFound: false},
		{name: "child", topic: MustTopic("diagnostics.perf"), want: MustTopic("diagnostics"), wantFound: true},
		{name: "deep", topic: MustTopic("a.b.c"), want: MustTopic("a.b"), wantFound: true},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, found := test.topic.Parent()
			if found != test.wantFound {
				t.Fatalf("Parent() found = %v, want %v", found, test.wantFound)
			}

			if got != test.want {
				t.Fatalf("Parent() topic = %q, want %q", got, test.want)
			}

			if test.topic.HasParent() != test.wantFound {
				t.Fatalf("HasParent() = %v, want %v", test.topic.HasParent(), test.wantFound)
			}
		})
	}
}

func TestTopicHasPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		topic  Topic
		prefix Topic
		want   bool
	}{
		{name: "same topic", topic: MustTopic("diagnostics.perf"), prefix: MustTopic("diagnostics.perf"), want: true},
		{name: "parent prefix", topic: MustTopic("diagnostics.perf"), prefix: MustTopic("diagnostics"), want: true},
		{name: "grandparent prefix", topic: MustTopic("a.b.c"), prefix: MustTopic("a"), want: true},
		{name: "similar text is not hierarchy", topic: MustTopic("diagnostic-tools.perf"), prefix: MustTopic("diagnostics"), want: false},
		{name: "sibling is false", topic: MustTopic("diagnostics.perf"), prefix: MustTopic("diagnostics.trace"), want: false},
		{name: "empty topic", topic: "", prefix: MustTopic("diagnostics"), want: false},
		{name: "empty prefix", topic: MustTopic("diagnostics.perf"), prefix: "", want: false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.topic.HasPrefix(test.prefix); got != test.want {
				t.Fatalf("%q.HasPrefix(%q) = %v, want %v", test.topic, test.prefix, got, test.want)
			}
		})
	}
}

func TestTopicTrimPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		topic  Topic
		prefix Topic
		want   Topic
		ok     bool
	}{
		{name: "parent", topic: MustTopic("diagnostics.perf.stat"), prefix: MustTopic("diagnostics"), want: MustTopic("perf.stat"), ok: true},
		{name: "self", topic: MustTopic("diagnostics.perf"), prefix: MustTopic("diagnostics.perf"), want: "", ok: true},
		{name: "missing", topic: MustTopic("diagnostics.perf"), prefix: MustTopic("benchmark"), want: "", ok: false},
		{name: "empty prefix", topic: MustTopic("diagnostics.perf"), prefix: "", want: "", ok: false},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, ok := test.topic.TrimPrefix(test.prefix)
			if ok != test.ok {
				t.Fatalf("TrimPrefix() ok = %v, want %v", ok, test.ok)
			}

			if got != test.want {
				t.Fatalf("TrimPrefix() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestTopicAppend(t *testing.T) {
	t.Parallel()

	parent := MustTopic("diagnostics")

	child, err := parent.Append("perf")
	if err != nil {
		t.Fatalf("Append returned unexpected error: %v", err)
	}

	if got, want := child, MustTopic("diagnostics.perf"); got != want {
		t.Fatalf("Append() = %q, want %q", got, want)
	}
}

func TestTopicAppendToZeroTopic(t *testing.T) {
	t.Parallel()

	var topic Topic

	child, err := topic.Append("benchmarks")
	if err != nil {
		t.Fatalf("Append returned unexpected error: %v", err)
	}

	if got, want := child, MustTopic("benchmarks"); got != want {
		t.Fatalf("Append() = %q, want %q", got, want)
	}
}

func TestTopicAppendRejectsInvalidSegment(t *testing.T) {
	t.Parallel()

	tests := []string{"", "Perf", "perf.stat", "perf_stat", "1perf", "-perf", "перф"}

	for _, segment := range tests {
		segment := segment

		t.Run(segment, func(t *testing.T) {
			t.Parallel()

			_, err := MustTopic("diagnostics").Append(segment)
			if err == nil {
				t.Fatalf("Append(%q) returned nil error", segment)
			}

			if !errors.Is(err, ErrInvalidTopic) {
				t.Fatalf("Append(%q) error = %v, want ErrInvalidTopic", segment, err)
			}
		})
	}
}

func TestTopicAppendRejectsInvalidParent(t *testing.T) {
	t.Parallel()

	_, err := Topic("Diagnostics").Append("perf")
	if err == nil {
		t.Fatalf("Append returned nil error")
	}

	if !errors.Is(err, ErrInvalidTopic) {
		t.Fatalf("Append error = %v, want ErrInvalidTopic", err)
	}
}

func TestTopicMustAppendPanicsForInvalidSegment(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustTopic("diagnostics").MustAppend("Perf")
	})
}

func TestTopicJoin(t *testing.T) {
	t.Parallel()

	joined, err := MustTopic("diagnostics").Join(MustTopic("perf.stat"))
	if err != nil {
		t.Fatalf("Join() returned unexpected error: %v", err)
	}

	if got, want := joined, MustTopic("diagnostics.perf.stat"); got != want {
		t.Fatalf("Join() = %q, want %q", got, want)
	}
}

func TestTopicJoinWithZeroSide(t *testing.T) {
	t.Parallel()

	joined, err := Topic("").Join(MustTopic("profiling"))
	if err != nil {
		t.Fatalf("Join() returned unexpected error: %v", err)
	}

	if got, want := joined, MustTopic("profiling"); got != want {
		t.Fatalf("Join() = %q, want %q", got, want)
	}

	joined, err = MustTopic("profiling").Join("")
	if err != nil {
		t.Fatalf("Join() returned unexpected error: %v", err)
	}

	if got, want := joined, MustTopic("profiling"); got != want {
		t.Fatalf("Join() = %q, want %q", got, want)
	}
}

func TestTopicJoinRejectsInvalidSide(t *testing.T) {
	t.Parallel()

	_, err := MustTopic("diagnostics").Join(Topic("Perf"))
	if err == nil {
		t.Fatalf("Join() returned nil error")
	}

	if !errors.Is(err, ErrInvalidTopic) {
		t.Fatalf("Join() error = %v, want ErrInvalidTopic", err)
	}
}

func TestTopicMustJoinPanicsForInvalidSide(t *testing.T) {
	t.Parallel()

	assertPanics(t, func() {
		_ = MustTopic("diagnostics").MustJoin(Topic("Perf"))
	})
}
